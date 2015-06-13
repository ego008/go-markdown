markdown [![GoDoc](http://godoc.org/github.com/opennota/markdown?status.svg)](http://godoc.org/github.com/opennota/markdown)
========

opennota/markdown package provides CommonMark-compliant markdown parser and renderer, written in Go.

## Installation

    go get github.com/opennota/markdown

You can also go get [mdtool](https://github.com/opennota/mdtool), an example command-line tool:

    go get github.com/opennota/mdtool

## Standards support

Currently supported CommonMark spec: [v0.20](http://spec.commonmark.org/0.20/).

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
  Nofollow        | bool   | whether to add `rel="nofollow"` to links                    | false
  XHTMLOutput     | bool   | whether to output XHTML instead of HTML                     | false

## Benchmarks

Rendering spec/spec-0.20.txt on a Intel(R) Core(TM) i5-2400 CPU @ 3.10GHz

    BenchmarkRenderSpec019NoHTML      200     6285825 ns/op    2946110 B/op     8207 allocs/op
    BenchmarkRenderSpec019            100    15764980 ns/op    4879226 B/op    41539 allocs/op
    BenchmarkRenderSpecBlackFriday    200     7667491 ns/op    2750367 B/op    37247 allocs/op

## TODO

  * Improve performance with the raw HTML option enabled
  * Write an URL parser/encoder that would support decoding punycode and counting matching `[(` and `)]` in URLs

## License

GNU GPL v3+
