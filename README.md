markdown [![GoDoc](http://godoc.org/github.com/opennota/markdown?status.svg)](http://godoc.org/github.com/opennota/markdown)
========

opennota/markdown package provides CommonMark-compliant markdown parser and renderer, written in Go.

## Installation

    go get github.com/opennota/markdown

You can also go get [mdtool](https://github.com/opennota/mdtool), an example command-line tool:

    go get github.com/opennota/mdtool

## Standards support

Currently supported CommonMark spec: [v0.19](http://spec.commonmark.org/0.19/).

## Extensions

Besides the features required by CommonMark, opennota/markdown supports:

  * Tables (GFM)
  * Strikethrough (GFM)
  * Autoconverting plain-text URLs to links
  * Typographic replacements (smart quotes and other)

## Usage

``` go
md := markdown.New(markdown.XHTMLOutput(true), markdown.Nofollow(true))
fmt.Println(md.RenderToString([]byte("Header\n===\nText")))
```

Check out [the source of mdtool](https://github.com/opennota/mdtool/blob/master/main.go) for a more complete example.

The following options are currently supported:

  Name            |  Type  |                        Description                          | Default
  --------------- | ------ | ----------------------------------------------------------- | ---------
  HTML            | bool   | whether to enable raw HTML                                  | false
  Tables          | bool   | whether to enable GFM tables                                | true
  Linkify         | bool   | whether to autoconvert plain-text URLs to links             | true
  Typographer     | bool   | whether to enable typographic replacements                  | true
  Quotes          | string | double + single quote replacement pairs for the typographer | “”‘’
  MaxNesting      | int    | maximum block nesting level                                 | 20
  LangPrefix      | string | CSS language prefix for fenced blocks                       | language-
  Breaks          | bool   | whether to convert newlines inside paragraphs into `<br>`   | false
  Nofollow        | bool   | whether to add `rel="nofollow"` to links                      | false
  XHTMLOutput     | bool   | whether to output XHTML instead of HTML                     | false

## Benchmarks

Rendering spec/spec-0.19.txt on a Intel(R) Core(TM) i5-2400 CPU @ 3.10GHz

    BenchmarkRenderSpec019NoHTML     200           6116606 ns/op         2921682 B/op       7929 allocs/op
    BenchmarkRenderSpec019           100          15744611 ns/op         4851484 B/op      41032 allocs/op
    BenchmarkRenderSpecBlackFriday   200           7450171 ns/op         2722858 B/op      36689 allocs/op

## License

GNU GPL v3+
