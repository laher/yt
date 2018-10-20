---
title: "Getting Started"
date: 2018-10-20T15:58:49+13:00
toc: true
menu:
  main:
      {}
---

## Installation

I'll create a release once I'm happy. Until then, you'll need to [install go](https://golang.org/doc/install) in order to install yt.

    go get github.com/laher/yt

## Run

The default behaviour parses the input as YAML and spits it out in the same format.

Try running `yt` against the `testdata/k8s.yaml` file provided in this repository (it's an anonymised kubernetes file, containing 3 yaml documents)

```
    yt < testdata/k8s.yaml 
```

This is effectively the same as setting the main query to `'{{.|yaml}}'`

```
    yt -q '{{.|yaml}}' < testdata/k8s.yaml 
```

## Query syntax

See our [Introduction to Templating](/yt/posts/templating/)

### Queries by example

These are just a few examples of querying with `yt`

#### Nested items

(outputted in Go's 'sprintf' format)

```
   yt -q '{{.metadata.labels.app}}' < testdata/k8s.yaml
```

#### Selected functions

Please see golang.org/pkg/text/template for a comprehensive list of built-in functions. These are just a few examples

##### index is useful when one of your keys itself contains a dot

```
   yt -q '{{index .data "config.json"}}' < testdata/k8s.yaml
```

##### `js` escapes javascript

(pipes are similar to jq's piping syntax)

```
  yt -q '{{index .data "config.json"|js}}' < testdata/k8s.yaml
```

##### `go` generates a Go-syntax representation of the result

```
  yt -q '{{.metadata|go}}' < testdata/k8s.yaml
```

##### or

```
  yt -q '{{or .metadata.labels.app .spec.replicas}}' < testdata/k8s.yaml
```

There's lots of other built-in stuff, check go's docs.

## Data sources:

You can specify multiple sources of data.yaml

```
  yt -d testdata/k8s.yaml -d additional=testdata/additional-data.yaml -q '{{ . |yaml}}{{ with (ds "additional) }}{{ .|yaml}}{{ end }}'
```

## Selecting a root doc from a multi-document input

You can select a different root document during datasource selection.

e.g. for the second document (index 1): `yt -q '{{ (ds "." 1).kind }}' < testdata/k8s.yaml`

### Merging data sources:

Write some data from another doc into this doc

```
yt -q '{{set . "data" (ds "x").data}}{{.|yaml}}' -d x=additional-data.yaml < testdata/k8s.yaml
```

## Large files and streams

`yt` is not very efficient with large files. Don't use it for streams, it's not ready for that yet. Perhaps I'll convert it to use more stream-oriented parsing in the future.

In the meantime, use `-maxBufferSize=1000000` to manage large files. The default should work fine with smmallish files.

