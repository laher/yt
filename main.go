package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/spf13/afero"
)

const (
	defaultTemplate      = "{{.|yaml}}"
	maxBufferSizeDefault = bufio.MaxScanTokenSize
)

var (
	query = flag.String("q", defaultTemplate, "Main yaml query [unless overridden by -t templates]")
	//documentSplitQuery = flag.String("dq", "", "Document split query")
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
		fmt.Fprint(flag.CommandLine.Output(), `YAML tarantosaurus

Usage:
 yt -d data.yaml
 cat data.yaml|yt -q '{{ .kind }}'
 
`)
		flag.PrintDefaults()
	}
	flag.Var(&dataSources, "d", "Data source(s)")
	flag.Var(&templates, "t", "Template file(s)")
	flag.Parse()
	// get data ...
	data, err := getSources(dataSources)
	if err != nil {
		log.Fatal(err)
	}
	if _, ok := data[mainSource]; !ok {
		data[mainSource] = source{
			typ: stdin,
		}
	}
	tpls, err := getSources(templates)
	if err != nil {
		log.Fatal(err)
	}
	if _, ok := tpls[mainSource]; !ok {
		tpls[mainSource] = source{
			typ: str,
			str: *query,
		}
	}

	err = spout(os.Stdout, data, tpls, *maxBufferSize)
	if err != nil {
		log.Fatal(err)
	}
}
