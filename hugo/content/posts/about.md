---
title: "About"
date: 2018-10-20T14:50:39+13:00
toc: true
menu:
  main:
    other: {}
---

`yt` is a [YAML](http://www.yaml.org/) processor in the spirit of [`jq`](https://stedolan.github.io/jq/) and [sed](https://en.wikipedia.org/wiki/Sed).  It can take in many yml files, process them using gomplate language, and then create a new .yml (or .json) file.

```
    +------+                                      
    |      |                                      
    | .yml |\                                     
    |      | -\                                   
    +------+   -\                                 
                 -\                               
    +------+       -\+-----------+        +------+
    |      |         -           |        |      |
    | .yml | --------|    yt     |--------| .yml |
    |      |        -|           |        |      |
    +------+      -/ +-----------+        +------+
                 /                                
    +------+   -/                                 
    |      | -/                                   
    | .yml |/                                     
    |      |                                      
    +------+                                      
```

 * `yt` can work with multi-document YAML files.
 * `yt` is experimental - please use with care.

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

I wanted a tool to parse and generate yaml. `yt` does this and also helps in some interesting ways.

## Acknowledgements

All of the hard stuff was done in [go-yaml](https://github.com/go-yaml/yaml) and in go itself. Thanks all.

yt's name is deliberately similar to yq and jq. The y is for yaml and the t might be for 'template', maybe for 'tool', or most likely for 'tarantosaurus'.

Thanks to [gomplate](https://github.com/hairyhenderson/gomplate) for a wonderful example of a cli tool harnessing go's text/template. I have borrowed some ideas from gomplate, especially the argument style for data sources and scripts/templates. I did contribute towards gomplate's templates feature, but the focus of yt is different, so there it is.

Thanks also to [Hugo](https://gohugo.io) for similar inspiration
