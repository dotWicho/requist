# requist

[![Go](https://github.com/dotWicho/requist/workflows/Go/badge.svg?branch=master)](https://github.com/dotWicho/requist)
[![Quality Report](https://goreportcard.com/badge/github.com/dotWicho/requist)](https://goreportcard.com/badge/github.com/dotWicho/requist)
[![GoDoc](https://godoc.org/github.com/dotWicho/requist?status.svg)](https://pkg.go.dev/github.com/dotWicho/requist?tab=doc)

## Requist is a Go Library to manage HTTP/S Requests easily

## Getting started

- API documentation is available via [godoc](https://godoc.org/github.com/dotWicho/requist).
- The [examples](./examples) directory contains more elaborate example applications.

## Installation

To install Requist package, you need to install Go and set your Go workspace first.

1 - The first need [Go](https://golang.org/) installed (**version 1.13+ is required**).
Then you can use the below Go command to install Requist

```bash
$ go get -u github.com/dotWicho/requist
```

And then Import it in your code:

``` go
package main

import "github.com/dotWicho/requist"
```
Or

2 - Use as module in you project (go.mod file):

``` go
module myclient

go 1.13

require (
	github.com/dotWicho/requist v1.2.2
)
```

## Contributing

- Get started by checking our [contribution guidelines](https://github.com/dotWicho/requist/blob/master/CONTRIBUTING.md).
- Read the [dotWicho requist wiki](https://github.com/dotWicho/requist/wiki) for more technical and design details.
- If you have any questions, just ask!

