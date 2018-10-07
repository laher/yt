package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/template"
)

func spout(output io.Writer, dataSources map[string]source, templateSources map[string]source, maxBufferSize int) error {
	mainDS, ok := dataSources[mainSource]
	if !ok {
		return fmt.Errorf("Main data source undefined")
	}
	rdr, err := mainDS.GetReader()
	if err != nil {
		return err
	}
	data, err := unmarshalInput(rdr, maxBufferSize)
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
			rdr, err := os.Open(ds.path)
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
	trdr, err := templateSources[mainSource].GetReader()
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(trdr)
	if err != nil {
		return err
	}
	tmpl, err := template.New("main").Funcs(funcMap).Parse(string(b))
	if err != nil {
		return err
	}
	for k, s := range templateSources {
		switch s.typ {
		case str, stdin:
			rdr, err := s.GetReader()
			if err != nil {
				return err
			}
			b, err := ioutil.ReadAll(rdr)
			if err != nil {
				return err
			}
			_, err = tmpl.New(k).Parse(string(b))
			if err != nil {
				return err
			}
		case file:
			if _, err := tmpl.New(k).ParseFiles(s.path); err != nil {
				return err
			}
		}

	}

	// Run the template to verify the output.
	err = tmpl.Execute(output, data)
	if err != nil {
		return err
	}
	return nil
}
