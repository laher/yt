# yt

`yt` is a [YAML](http://www.yaml.org/) processor in the spirit of [`jq`](https://stedolan.github.io/jq/) and [sed](https://en.wikipedia.org/wiki/Sed).

`yt` can work with multi-document YAML files.

`yt` is VERY experimental - please use with care.

## Primary Goals

 * Query yaml documents
 * Manipulate yaml documents
 * Merge yamls
 * Support multi-document files
## Features

`yt` offers:

 * a query mechanism so that you can print out parts of your document.
 * an approach to merging and combining multiple yaml documents

`yt` uses Go's text/template as its querying engine. It's never going to be as terse or as flexible as `jq`, but it's fine for the basics, and I (or you) can keep on adding template functions.

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

You can select a different root document during datasource selection.

e.g. for the second document (index 1): `yt -q '{{ (ds "." 1).kind }}' < sample.yaml`

### Additional data sources:

Write some data from another doc into this doc

```
yt -q '{{set . "data" (ds "x").data}}{{.|yaml}}' -d x=additional-data.yaml < sample.yaml
```

## Query syntax

yt's query syntax comes from the [text/template package](https://golang.org/pkg/text/template) from [Go](https://golang.org)'s standard library.

## Large files and streams

`yt` is not very efficient with large files. Don't use it for streams, it's not ready for that yet. Perhaps I'll convert it to use more stream-oriented parsing in the future.

In the meantime, use `-maxBufferSize=1000000` to manage large files. The default should work fine with smmallish files.

## Background

I originally wrote this tool in order to _query_ multi-document YAML files (kubernetes/kubectl configuration files). 

Some time later I found myself working with bigger, more repetitive yaml files (Concourse CI pipelines), and wanting to 'template' repetitive parts of the document.

Multi-document note: YAML's multi-document syntax is a really useful feature, but unfortunately some YAML tools do not inherently support multi-document files. ([go-yaml](https://github.com/go-yaml/yaml) itself just picks the first document and ignores the rest. I have also had other problems with tools I've tried, particularly with nested json inside yaml (I know, I know).

## Acknowledgements

All of the hard stuff was done in [go-yaml](https://github.com/go-yaml/yaml) and in go itself. Thanks all. 

yt's name is deliberately similar to yq and jq. The y is for yaml and the t might be for 'template', maybe for 'tool', or most likely for 'tarantosaurus'.

Thanks to [gomplate](https://github.com/hairyhenderson/gomplate) for a wonderful example of a cli tool harnessing go's text/template. I have borrowed some ideas from gomplate, especially the argument style for data sources / templates. I did contribute towards those templates, but there it is.
