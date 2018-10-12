package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

	// index is represented as '...int' just so that it can be an optional parameter
	funcs := funcMap()
	funcs["ds"] = func(k string, index ...int) map[interface{}]interface{} {
		ds, ok := dataSources[k]
		if !ok {
			// not found. Die with helpful list of data sources
			kbuf := ""
			for k := range dataSources {
				if len(kbuf) > 0 {
					kbuf += ","
				}
				kbuf += "'" + k + "'"
			}
			log.Fatalf("Data source not found. Data sources: [%s]", kbuf)

		}
		var (
			rdr io.ReadCloser
			err error
		)
		rdr, err = fs.Open(ds.path)
		if err != nil {
			log.Fatal(err)
		}
		if len(index) > 0 {
			rdr, err = split(rdr, index[0], maxBufferSize)
			if err != nil {
				log.Fatal(err)
			}
		}
		ret, err := unmarshalInput(rdr, maxBufferSize)
		if err != nil {
			log.Fatal(err)
		}
		return ret
	}
	funcs["dss"] = func() []string {
		dss := []string{}
		for k := range dataSources {
			dss = append(dss, k)
		}
		return dss
	}
	trdr, err := templateSources[mainSource].GetReader()
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(trdr)
	if err != nil {
		return err
	}
	tmpl, err := template.New("main").Funcs(funcs).Parse(string(b))
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
