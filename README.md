<img src="yam.jpg" height="96" align="left" />

# YAM

[![Build Status](https://img.shields.io/travis/thisissoon/yam/master.svg)](https://travis-ci.org/thisissoon/yam)
[![Coverage](https://img.shields.io/coveralls/thisissoon/yam.svg)](https://coveralls.io/github/thisissoon/yam)
[![Docs](https://img.shields.io/badge/documentation-godoc-375eab.svg)](https://godoc.org/github.com/thisissoon/yam)
[![Licsense](https://img.shields.io/badge/LICENSE-MIT-blue.svg)](https://github.com/thisissoon/yam/blob/master/LICENSE)

YAM (Yet Another Mux) is another Golang HTTP multiplexer designed to be simple, flexible
and configurable.

``` go
package main

import (
	"net/http"

	"github.com/thisissoon/yam"
)

func main() {
	mux := yam.New()
	mux.Route("/").Get(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	}))

	http.ListenAndServe(":5000", mux)
}
```

## Install

```
go get github.com/thisissoon/yam
```

## Documentaion

Full package documentation can be found at https://godoc.org/github.com/thisissoon/yam.
