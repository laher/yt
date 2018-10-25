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
	defaultQuery         = "{{.|yaml}}"
	maxBufferSizeDefault = bufio.MaxScanTokenSize
)

var (
	query = flag.String("q", defaultQuery, "Main yaml query [unless overridden by -s script.yt]")
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
var scripts strslice

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
	flag.Var(&scripts, "s", "Script file(s)")
	flag.Parse()
	args := flag.Args()
	// get data ...
	data, err := getSources(dataSources)
	if err != nil {
		log.Fatal(err)
	}
	if _, ok := data[mainSource]; !ok {
		switch len(args) {
		case 1:
			data[mainSource] = source{
				typ:  file,
				path: args[0],
			}
		case 0:
			data[mainSource] = source{
				typ: stdin,
			}
		default:
			log.Fatal("Unsupported (multiple cli args). Please pipe to yt or use a single cli arg")

		}
	}
	tpls, err := getSources(scripts)
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
