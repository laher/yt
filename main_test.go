package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

func TestSpout(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		query    string
		expected string
	}{
		{name: "single value", input: `number: 1`, query: `{{.number}}`, expected: `1`},
		{name: "to yaml", input: `number: 1`, query: `{{.|yaml}}`, expected: `number: 1`},
		{name: "to go", input: `number: 1`, query: `{{.}}`, expected: `map[number:1]`},
		{name: "nested query", input: `spec:
  replicas: 1
  template:
    metadata:
      labels:
        name: myservice
        app: myservice`, query: `{{.spec.template.metadata.labels.name}}`, expected: `myservice`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//i := ioutil.NopCloser(strings.NewReader(test.input))
			o := bytes.NewBuffer([]byte{})
			templates := map[string]source{
				mainSource: source{
					typ: stdin,
					Reader: func() io.ReadCloser {
						return ioutil.NopCloser(strings.NewReader(test.query))
					},
				},
			}
			data := map[string]source{
				mainSource: source{
					typ: stdin,
					Reader: func() io.ReadCloser {
						return ioutil.NopCloser(strings.NewReader(test.input))
					},
				},
			}
			err := spout(o, data, templates, maxBufferSizeDefault)
			if err != nil {
				t.Errorf("Unexpected error")
				return
			}
			actual := strings.TrimSpace(o.String())
			if test.expected != actual {
				t.Errorf("Expected [%s] not equal to actual [%s]", test.expected, actual)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		docIndex    int
		expected    string
		expectedErr error
	}{
		{name: "simple", input: `number: 1`, docIndex: 0, expected: `number: 1`},
		{name: "non-matching docIndex", input: `number: 1`, docIndex: 1, expectedErr: errNoMatch},
		{name: "index second doc", input: `---
number: 1
---
number: 2`, docIndex: 1, expected: `number: 2`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			i := ioutil.NopCloser(strings.NewReader(test.input))
			o, err := split(i, test.docIndex, maxBufferSizeDefault)
			if test.expectedErr == nil {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
					return
				}
				b, err := ioutil.ReadAll(o)
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				actual := strings.TrimSpace(string(b))
				if test.expected != actual {
					t.Errorf("Expected [%s] not equal to actual [%s]", test.expected, actual)
				}
			} else {
				if err != test.expectedErr {
					t.Errorf("Expected error [%v] but got [%v]", test.expectedErr, err)
				}
			}
		})
	}
}
