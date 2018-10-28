---
title: "Scripting Intro"
date: 2018-10-20T03:39:23+13:00
toc: true
---

`yt` uses the [text/template package](https://golang.org/pkg/text/template) from [Go](https://golang.org)'s standard library.

___NOTE__: although text/template is a templating language, and yt allows you to use it as such, the intention is more to use it as a scripting language. This allows your yaml to remain as valid yaml_

_For querying yaml, you won't really need to understand Go templates in depth, but once you're generating yaml, then it pays to understand the dialect..._

The following is only a primer on Go Templates, adapted from [Hugo's](https://gohugo.io) [documentation](https://gohugo.io/templates/introduction/). For an in-depth look into Go Templates, check the official Go docs.

Go Templates provide an extremely simple scripting language which we use in 2 ways:

 * Scripting for the view layer
 * Interpolation for individual variable within a yaml doc

## Basic Syntax

Go Templates are files with the addition of [variables]({{< relref "#variables" >}}) and [functions]({{< relref "#functions" >}}). Go Template variables and functions are accessible within `{{ }}`.

### Access a Predefined Variable

A _predefined variable_ could be a variable already existing in the
current scope (like the `.Title` example in the [Variables]({{< relref
"#variables" >}}) section below) or a custom variable (like the
`$address` example in that same section).

```go-text-template
{{ .Title }}
{{ $address }}
```

Parameters for functions are separated using spaces. The general syntax is:

```
{{ FUNCTION ARG1 ARG2 .. }}
```

The following example calls the `add` function with inputs of `1` and `2`:

```go-text-template
{{ add 1 2 }}
```

### Methods and Fields 

Methods and fields are Accessed via dot Notation

Accessing the variable `bar`:

```go-text-template
{{ .bar }}
```

### Parentheses 

Parentheses can be used to group items together for example:

```go-text-template
{{ if or (isset . "alt") (isset . "caption") }} Caption {{ end }}
```

## Variables {#variables}

Each Go Template gets a data object. In `yt`, this is the `main` data source - either STDIN or a named file.

```go-text-template
{{ .Title }}
```

Values can also be stored in custom variables and referenced later. The custom variables need to be prefixed with `$`.  For example:

```go-text-template
{{ $address := "123 Main St." }}
{{ $address }}
```

## Functions

Go Templates only ship with a few basic functions but also provide a mechanism for applications to extend the original set.

[yt scripting functions][functions] provide additional functionality specific to building websites. Functions are called by using their name followed by the required parameters separated by spaces.

### Example 1: Adding Numbers

```go-text-template
{{ add 1 2 }}
<!-- prints 3 -->
```

### Example 2: Comparing Numbers

```go-text-template
{{ lt 1 2 }}
<!-- prints true (i.e., since 1 is less than 2) -->
```
Note that both examples make use of Go Template's [math functions][].

> __NOTE__: there are more boolean operators than those listed in these docs in the [Go Template documentation](http://golang.org/pkg/text/template/#hdr-Functions).

## Includes

When including another script, you will need to pass it the data that it would need to access.

> __NOTE__: To pass along the current context, please remember to include a trailing **dot**.

### Scripts

The [`template`][template] function is used to include additional scripts using
the syntax `{{ template "_internal/<TEMPLATE>.<EXTENSION>" . }}`.

_TODO: provide a wrapper called 'script'_

Example:

```go-text-template
{{ template "myreport.tpl" . }}
```

## Logic

Go Templates provide the most basic iteration and conditional logic.

### Iteration

The Go Templates make heavy use of `range` to iterate over a _map_,
_array_, or _slice_. The following are different examples of how to
use `range`.

#### Example 1: Using Context (`.`)

```go-text-template
{{ range $array }}
    {{ . }} <!-- The . represents an element in $array -->
{{ end }}
```

#### Example 2: Declaring a variable name for an array element's value

```go-text-template
{{ range $elem_val := $array }}
    {{ $elem_val }}
{{ end }}
```

#### Example 3: Declaring variable names for an array element's index _and_ value

For an array or slice, the first declared variable will map to each
element's index.

```go-text-template
{{ range $elem_index, $elem_val := $array }}
   {{ $elem_index }} -- {{ $elem_val }}
{{ end }}
```

#### Example 4: Declaring variable names for a map element's key _and_ value

For a map, the first declared variable will map to each map element's
key.

```go-text-template
{{ range $elem_key, $elem_val := $map }}
   {{ $elem_key }} -- {{ $elem_val }}
{{ end }}
```

### Conditionals

`if`, `else`, `with`, `or`, and `and` provide the framework for handling conditional logic in Go Templates. Like `range`, each statement is closed with an `{{ end }}`.

Go Templates treat the following values as **false**:

- `false` (boolean)
- 0 (integer)
- any zero-length array, slice, map, or string

#### Example 1: `with`

It is common to write "if something exists, do this" kind of
statements using `with`.

> __NOTE:__ `with` rebinds the context `.` within its scope (just like in `range`).

It skips the block if the variable is absent, or if it evaluates to
"false" as explained above.

```go-text-template
{{ with .title }}

    <h4>{{ . }}</h4>
{{ end }}
```

#### Example 2: `with` .. `else`

Below snippet uses the "description" front-matter parameter's value if
set, else uses the default `.Summary` [variable]:

```go-text-template
{{ with . "description" }}
    {{ . }}
{{ else }}
    {{ .Summary }}
{{ end }}
```

#### Example 3: `if`

An alternative (and a more verbose) way of writing `with` is using
`if`. Here, the `.` does not get rebinded.

Below example is "Example 1" rewritten using `if`:

```go-text-template
{{ if isset . "title" }}
    <h4>{{ index . "title" }}</h4>
{{ end }}
```

#### Example 4: `if` .. `else`

Below example is "Example 2" rewritten using `if` .. `else`, and using
[`isset` function][isset] + `.` variable instead:

```go-text-template
{{ if (isset . "description") }}
    {{ index . "description" }}
{{ else }}
    {{ .Summary }}
{{ end }}
```

#### Example 5: `if` .. `else if` .. `else`

Unlike `with`, `if` can contain `else if` clauses too.

```go-text-template
{{ if (isset . "description") }}
    {{ index . "description" }}
{{ else if (isset . "summary") }}
    {{ index . "summary" }}
{{ else }}
    {{ .Summary }}
{{ end }}
```

#### Example 6: `and` & `or`

```go-text-template
{{ if (and (or (isset . "title") (isset . "caption")) (isset . "attr")) }}
```

## Pipes

One of the most powerful components of Go Templates is the ability to stack actions one after another. This is done by using pipes. Borrowed from Unix pipes, the concept is simple: each pipeline's output becomes the input of the following pipe.

Because of the very simple syntax of Go Templates, the pipe is essential to being able to chain together function calls. One limitation of the pipes is that they can only work with a single value and that value becomes the last parameter of the next pipeline.

A few simple examples should help convey how to use the pipe.

### Example 1: `shuffle`

The following two examples are functionally the same:

```go-text-template
{{ shuffle (seq 1 5) }}
```

```go-text-template
{{ (seq 1 5) | shuffle }}
```

### Example 2: `index`

The following accesses a yaml key called "title". This example also uses the [`index` function][index], which is built into Go Templates:

```go-text-template
{{ index . "title" }}
```

### Example 3: `or` with `isset`

```go-text-template
{{ if or (or (isset . "title") (isset . "caption")) (isset . "attr") }}
Stuff Here
{{ end }}
```

Could be rewritten as

```go-text-template
{{ if isset . "caption" | or isset . "title" | or isset . "attr" }}
Stuff Here
{{ end }}
```
