package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	defaultTemplate      = "{{.|yaml}}"
	maxBufferSizeDefault = bufio.MaxScanTokenSize
)

var (
	query              = flag.String("q", defaultTemplate, "Main yaml query [unless overridden by -t templates]")
	documentSplitQuery = flag.String("dq", "", "Document split query")
	index              = flag.Int("di", 0, "Select doc by order of appearance in the input")
	maxBufferSize      = flag.Int("b", maxBufferSizeDefault, "Max buffer size per input file")
	errNoMatch         = fmt.Errorf("no match")
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
	for i, dataArg := range dataSources {
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
