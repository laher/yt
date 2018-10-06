package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

var (
	query                = flag.String("q", "{{.|yaml}}", "Main yaml query")
	documentSplitQuery   = flag.String("dq", "", "Document split query")
	index                = flag.Int("di", 0, "Select doc by order of appearance in the input")
	maxBufferSizeDefault = bufio.MaxScanTokenSize
	maxBufferSize        = flag.Int("b", maxBufferSizeDefault, "Max buffer size")
	errNoMatch           = fmt.Errorf("no match")
)

func toYAML(x interface{}) (string, error) {
	bts, err := yaml.Marshal(x)
	return string(bts), err
}

func coerceKeys(in map[interface{}]interface{}) map[string]interface{} {

	tryMap := map[string]interface{}{}
	for k, v := range in {
		ks := fmt.Sprintf("%v", k)
		switch vT := v.(type) {
		case map[interface{}]interface{}:
			tryMap[ks] = coerceKeys(vT)
		default:
			tryMap[ks] = v
		}
	}
	return tryMap
}

func toJSON(x interface{}) (string, error) {
	switch xT := x.(type) {
	case map[interface{}]interface{}:
		tryMap := coerceKeys(xT)
		bts, err := json.Marshal(tryMap)
		return string(bts), err
	}
	bts, err := json.Marshal(x)
	return string(bts), err
}

func other(x interface{}) string {
	return fmt.Sprintf("%+v", x)
}

func toGo(x interface{}) string {
	return fmt.Sprintf("%#v", x)
}

func del(x interface{}, y string) interface{} {
	switch t := x.(type) {
	case map[interface{}]interface{}:
		delete(t, y)
	default:
		panic(fmt.Sprintf("x (%t) is not a map[interface{}]interface{}", x))
	}
	return x
}

func set(x map[interface{}]interface{}, k string, y interface{}) string {
	x[k] = y
	return "" // empty-string avoids output
}

type strslice []string

func (i *strslice) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *strslice) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var myStrs strslice

func main() {

	flag.Var(&myStrs, "d", "Data source(s)")
	flag.Parse()
	var (
		err   error
		input io.ReadCloser
		data  = map[string]string{}
	)
	input = os.Stdin

	if *index > 0 || *documentSplitQuery != "" {
		input, err = split(input, *documentSplitQuery, *index, *maxBufferSize)
		if err != nil {
			log.Fatal(err)
		}
	} // else just process the first doc
	for i, dataArg := range myStrs {
		parts := strings.Split(dataArg, "=")
		key := ""
		val := ""
		if len(parts) == 2 {
			key = parts[0]
			val = parts[1]
		} else if len(parts) == 1 {
			key = strconv.Itoa(i)
			val = parts[0]
		} else {
			log.Fatal("data should take the format key=filename.yaml")
		}

		data[key] = val
	}

	err = spout(input, os.Stdout, data, *query, *maxBufferSize)
	if err != nil {
		log.Fatal(err)
	}
}

func unmarshalInput(input io.ReadCloser, maxBufferSize int) (map[interface{}]interface{}, error) {
	y, err := ioutil.ReadAll(io.LimitReader(input, int64(maxBufferSize)))
	if err != nil {
		return nil, err
	}
	defer input.Close()

	data := make(map[interface{}]interface{})
	err = yaml.Unmarshal(y, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func spout(input io.ReadCloser, output io.Writer, dataSources map[string]string, query string, maxBufferSize int) error {

	data, err := unmarshalInput(input, maxBufferSize)
	if err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"yaml":      toYAML,
		"json":      toJSON,
		"o":         other,
		"go":        toGo,
		"del":       del,
		"set":       set,
		"tableflip": func() string { return "(╯°□°）╯︵ ┻━┻" },
		"ds": func(k string) interface{} {
			ds := dataSources[k]
			rdr, err := os.Open(ds)
			if err != nil {
				log.Fatal(err)
			}
			ret, err := unmarshalInput(rdr, maxBufferSize)
			if err != nil {
				log.Fatal(err)
			}
			return ret
		},
	}
	tmpl, err := template.New("test").Funcs(funcMap).Parse(query)
	if err != nil {
		return err
	}

	// Run the template to verify the output.
	err = tmpl.Execute(output, data)
	if err != nil {
		return err
	}
	return nil
}

func split(input io.ReadCloser, documentSplitQuery string, docIndex int, maxBufferSize int) (io.ReadCloser, error) {
	defer input.Close()
	s := bufio.NewScanner(input)
	s.Buffer(make([]byte, maxBufferSize), maxBufferSize)
	onNewDoc := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		buf := ""
		for i := 0; i < len(data); i++ {
			if data[i] == '\n' {
				if buf == "---" {
					return i + 1, data[:i], nil
				}
				buf = ""
			} else if data[i] == '-' {
				buf += "-"
			} else {
				buf = ""
			}
		}
		// There is one final token to be delivered, which may be the empty string.
		// Returning bufio.ErrFinalToken here tells Scan there are no more tokens after this
		// but does not trigger an error to be returned from Scan itself.
		return 0, data, bufio.ErrFinalToken
	}
	s.Split(onNewDoc)
	i := 0
	for s.Scan() {
		y := s.Bytes()
		if string(y) == "---" {
			//log.Printf("skip")
			continue
		}
		if documentSplitQuery != "" {
			//log.Printf("scanned: %s", string(y))
			data := make(map[interface{}]interface{})
			err := yaml.Unmarshal(y, &data)
			if err != nil {
				return nil, fmt.Errorf("unmarshaling: %s", err)
			}
			funcMap := template.FuncMap{
				"yaml":      toYAML,
				"o":         other,
				"tableflip": func() string { return "(╯°□°）╯︵ ┻━┻" },
			}
			tmpl, err := template.New("test").Funcs(funcMap).Parse(documentSplitQuery)
			if err != nil {
				return nil, fmt.Errorf("parsing: %s", err)
			}
			b := bytes.NewBuffer([]byte{})
			// Run the template to verify the output.
			err = tmpl.Execute(b, data)
			if err != nil {
				return nil, fmt.Errorf("execution: %s", err)
			}
			st := b.String()
			//log.Printf("output: [%s]", st)
			//println("[", st, "]", len(st), len(strings.TrimSpace(st)))

			if strings.TrimSpace(st) == "true" {
				//fmt.Printf("---%s", string(y))

				return ioutil.NopCloser(bytes.NewBuffer(y)), nil
			}
		} else {
			if i == docIndex {
				//fmt.Printf("---%s", string(y))
				return ioutil.NopCloser(bytes.NewBuffer(y)), nil
			}
		}
		i++
	}
	return nil, errNoMatch
}
