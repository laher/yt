package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
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
)

func toYAML(x interface{}) (string, error) {
	bts, err := yaml.Marshal(x)
	return string(bts), err
}

func other(x interface{}) string {
	return fmt.Sprintf("%+v", x)
}

func toGo(x interface{}) string {
	return fmt.Sprintf("%#v", x)
}

func main() {
	flag.Parse()
	var (
		err   error
		input io.ReadCloser
	)
	input = os.Stdin

	if *index > 0 || *documentSplitQuery != "" {

		input, err = split(input)
		if err != nil {
			log.Fatal(err)
		}
	} // else just process the first doc

	err = spout(input, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}

func spout(input io.ReadCloser, output io.Writer) error {
	y, err := ioutil.ReadAll(input)
	if err != nil {
		return err
	}
	defer input.Close()

	data := make(map[interface{}]interface{})
	err = yaml.Unmarshal(y, &data)
	if err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"yaml":      toYAML,
		"o":         other,
		"go":        toGo,
		"tableflip": func() string { return "(╯°□°）╯︵ ┻━┻" },
	}
	tmpl, err := template.New("test").Funcs(funcMap).Parse(*query)
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

func split(input io.ReadCloser) (io.ReadCloser, error) {
	defer input.Close()
	s := bufio.NewScanner(input)
	s.Buffer(make([]byte, *maxBufferSize), *maxBufferSize)
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
		if *documentSplitQuery != "" {
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
			tmpl, err := template.New("test").Funcs(funcMap).Parse(*documentSplitQuery)
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
			if i == *index {
				//fmt.Printf("---%s", string(y))
				return ioutil.NopCloser(bytes.NewBuffer(y)), nil
			}
		}
		i++
	}
	return nil, fmt.Errorf("no match")
}
