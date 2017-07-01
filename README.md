# yt

`yt` is a [YAML](http://www.yaml.org/) processor in the spirit of [`jq`](https://stedolan.github.io/jq/) and [sed](https://en.wikipedia.org/wiki/Sed).

`yt` can work with multi-document YAML files.

## Background

I wrote this tool to work with multi-document YAML files (see kubernetes/kubectl). 
It's a really useful feature, but unfortunately various YAML tools do not inherently support multi-document YAML files ([go-yaml](https://github.com/go-yaml/yaml) just picks the first one)

`yt` offers 2 ways to select a document, and then it offers a simple query mechanism so that you can print out parts of your document.

`yt` uses Go's internal templating engine to achieve this. It's never going to be as terse as `jq`, but it's fine for basic stuff.

## Installation

    go get github.com/laher/yt

I'll cut a release once I'm happy

## Run

The default behaviour parses the input as YAML and spits it out in the same format

```
    yt < sample.yaml 
```

This is effectively the same as setting the main query to '{{.|yaml}}'

```
    yt -q '{{.|yaml}}' < sample.yaml 
```

## Selecting a root doc from a stream

You can select a 'document index' or a 'document query'.

### Document index

Instead of selecting the first document (index 0), `yt` finds the second document (index 1).

```
    yt -di=1 < sample.yaml 
```

### Matching documents

Instead of selecting the first document, `yt` queries documents until it finds a document matching the 'document query'.

```
    yt -dq='{{eq .kind "ConfigMap"}}' < sample.yaml 
```

## Query syntax

yt's query syntax comes from the [text/template package](https://golang.org/pkg/text/template) from [Go](https://golang.org)'s standard library.


## Acknowledgements

All of the hard stuff was done in go-yaml and in go itself. Thanks all. 

yt's name is deliberately similar to yq and jq. The y is for yaml and the t is for template. The name is also a vague nod to `kt` by Felix Geller, whose seat I stole.
