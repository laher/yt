package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/afero"
	yaml "gopkg.in/yaml.v2"
)

const (
	mainSource = "."
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

func split(input io.ReadCloser, docIndex int, maxBufferSize int) (io.ReadCloser, error) {
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
		if i == docIndex {
			//fmt.Printf("---%s", string(y))
			return ioutil.NopCloser(bytes.NewBuffer(y)), nil
		}
		i++
	}
	return nil, errNoMatch
}

const (
	stdin = "stdin"
	str   = "string"
	file  = "file"
)

type source struct {
	typ  string
	path string
	str  string
}

func (s source) GetReader() (io.ReadCloser, error) {
	switch s.typ {
	case str:
		return ioutil.NopCloser(strings.NewReader(s.str)), nil
	case stdin:
		return os.Stdin, nil
	case file:
		return fs.Open(s.path)
	}
	return nil, fmt.Errorf("Unhandled source type")
}

func getSources(args []string) (map[string]source, error) {
	data := map[string]source{}
	for _, dataArg := range args {
		parts := strings.Split(dataArg, "=")
		key := ""
		val := ""
		if len(parts) == 2 {
			key = parts[0]
			val = parts[1]
		} else if len(parts) == 1 {
			key = mainSource
			val = parts[0]
		} else {
			return nil, errors.New("data should take the format key=filename.yaml")
		}
		filenames, err := afero.Glob(fs, val)
		if err != nil {
			return nil, err
		}
		i := 0
		for _, f := range filenames {
			fi, err := fs.Stat(f)
			if err != nil {
				return nil, err
			}
			if !fi.IsDir() {
				k := f
				if key != mainSource {
					k = fmt.Sprintf("%s:%s", key, f)
				}
				data[k] = source{typ: file, path: f}
				i++
			}
		}

		if i == 0 {
			return nil, fmt.Errorf("pattern matches no files: %#q", val)
		}
		// not a glob, nor a dir
		if len(filenames) == 1 && i == 1 && (!strings.ContainsAny(val, "*?")) {
			data[key] = source{typ: file, path: filenames[0]}
		}
	}
	return data, nil
}
