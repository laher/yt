# yt

`yt` is a [YAML](http://www.yaml.org/) processor in the spirit of [`jq`](https://stedolan.github.io/jq/) and [sed](https://en.wikipedia.org/wiki/Sed).

`yt` can work with multi-document YAML files.

`yt` is VERY experimental - please use with care.

## Background

I wrote this tool in order to work with multi-document YAML files (kubernetes/kubectl configuration files). 

YAML's multi-document syntax is a really useful feature, but unfortunately some YAML tools do not inherently support multi-document files. ([go-yaml](https://github.com/go-yaml/yaml) itself just picks the first document and ignores the rest. I have also had other problems with tools I've tried, particularly with nested json inside yaml (I know, I know).

## Features

`yt` offers:

 * 2 ways to select a document within a multiple-document file/stream.
 * a simple query mechanism so that you can print out parts of your document.
 * `yt` uses Go's text/template as its querying engine. It's never going to be as terse or as flexible as `jq`, but it's fine forthe basics, and I (or you) can keep on adding template functions.

## Installation

I'll create a release once I'm happy. Until then, you'll need to [install go](https://golang.org/doc/install) in order to install yt.

    go get github.com/laher/yt

## Run

The default behaviour parses the input as YAML and spits it out in the same format.

Try running `yt` against the `sample.yaml` file provided in this repository (it's an anonymised kubernetes file, containing 3 yaml documents)

```
    yt < sample.yaml 
```

This is effectively the same as setting the main query to `'{{.|yaml}}'`

```
    yt -q '{{.|yaml}}' < sample.yaml 
```

### Queries by example

Please see golang.org/pkg/text/template for more details. These are just a few examples

#### Nested items

(outputted in Go's 'sprintf' format)

```
   yt -q '{{.metadata.labels.app}}' < sample.yaml
```

#### Selected functions

Please see golang.org/pkg/text/template for a comprehensive list of built-in functions. These are just a few examples

##### index is useful when one of your keys itself contains a dot

```
   yt -q '{{index .data "config.json"}}' < sample.yaml
```

##### `js` escapes javascript

(pipes are similar to jq's piping syntax)

```
  yt -q '{{index .data "config.json"|js}}' < sample.yaml
```


##### `go` generates a Go-syntax representation of the result

```
  yt -q '{{.metadata|go}}' < sample.yaml
```

##### or

```
  yt -q '{{or .metadata.labels.app .spec.replicas}}' < sample.yaml
```

There's lots of other built-in stuff, check go's docs.

## Selecting a root doc from a multi-document input

You can select a 'document index' or a 'document query'.

### Document index

Instead of selecting the first document (index 0), `yt` can use the second document (index 1).

```
    yt -di=1 < sample.yaml 
```

### Matching documents

Instead of selecting the first document, `yt` can query documents to find a document matching the 'document query'.

```
    yt -dq='{{eq .kind "ConfigMap"}}' < sample.yaml 
```

## Query syntax

yt's query syntax comes from the [text/template package](https://golang.org/pkg/text/template) from [Go](https://golang.org)'s standard library.


## Acknowledgements

All of the hard stuff was done in go-yaml and in go itself. Thanks all. 

yt's name is deliberately similar to yq and jq. The y is for yaml and the t is for template. 

The name is also a vague nod to [kt](https://github.com/fgeller/kt), which was written by a man who gave me his chair.
