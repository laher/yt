package main

import (
	"html/template"
	"io"
	"log"
	"os"
)

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
