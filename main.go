package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/spf13/afero"
)

const (
	defaultTemplate      = "{{.|yaml}}"
	maxBufferSizeDefault = bufio.MaxScanTokenSize
)

var (
	query = flag.String("q", defaultTemplate, "Main yaml query [unless overridden by -t templates]")
	//documentSplitQuery = flag.String("dq", "", "Document split query")
	index         = flag.Int("di", 0, "Select doc by order of appearance in the input")
	maxBufferSize = flag.Int("b", maxBufferSizeDefault, "Max buffer size per input file")
	errNoMatch    = fmt.Errorf("no match")

	fs = afero.NewOsFs()
)

type strslice []string

func (i *strslice) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *strslice) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var dataSources strslice
var templates strslice

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `YAML tool

Usage:
 %s -d data.yaml
 %s < data.yaml
 
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}
	flag.Var(&dataSources, "d", "Data source(s)")
	flag.Var(&templates, "t", "Template file(s)")
	flag.Parse()
	/*
		var (
			err   error
			input io.ReadCloser
		)
		input = os.Stdin

		if *index > 0 { // || *documentSplitQuery != "" {
			input, err = split(input, *index, *maxBufferSize)
			if err != nil {
				log.Fatal(err)
			}

		} // else just process the first doc
	*/
	// get data ...
	data, err := getSources(dataSources)
	if err != nil {
		log.Fatal(err)
	}
	if _, ok := data[mainSource]; !ok {
		data[mainSource] = source{
			typ: str,
			Reader: func() io.ReadCloser {
				return os.Stdin
			},
		}
	}
	tpls, err := getSources(templates)
	if err != nil {
		log.Fatal(err)
	}
	if _, ok := tpls[mainSource]; !ok {
		tpls[mainSource] = source{
			typ: str,
			Reader: func() io.ReadCloser {
				return ioutil.NopCloser(strings.NewReader(*query))
			},
		}
	}

	err = spout(os.Stdout, data, tpls, *maxBufferSize)
	if err != nil {
		log.Fatal(err)
	}
}
