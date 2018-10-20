---
title: "About"
date: 2018-10-20T14:50:39+13:00
toc: true
menu:
  main:
    other: {}
---

`yt` is a [YAML](http://www.yaml.org/) processor in the spirit of [`jq`](https://stedolan.github.io/jq/) and [sed](https://en.wikipedia.org/wiki/Sed).

 * `yt` can work with multi-document YAML files.
 * `yt` is VERY experimental - please use with care.

## Goals

 * Query yaml documents
  * Support multi-document files (with `---` separator lines)
 * Manipulate yaml documents
  * Generate yaml
  * Merge and combine yaml docs
  * Keep yaml separate from looping / conditionals

## Why

Yaml seems to be everywhere at the moment, particularly in Ops tooling - see Kubernetes, Terraform, Concourse, Circle CI ...

I began writing this tool in 2017, initially to support multi-document files for Kubernetes. I came back to it again to try to generate and manipulate yaml docs for Concourse CI templates.

I wanted a tool to parse and generate yaml. yt helps in some interesting ways. 

