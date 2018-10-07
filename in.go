package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"text/template"

	yaml "gopkg.in/yaml.v2"
)

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
