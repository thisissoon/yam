# YAM

![Build Status](https://img.shields.io/travis/thisissoon/yam.svg)
![Coverage](https://img.shields.io/coveralls/thisissoon/yam.svg)
[![Docs](https://img.shields.io/badge/documentation-godoc-375eab.svg)](https://godoc.org/github.com/thisissoon/yam)
![Licsense](https://img.shields.io/badge/LICENSE-MIT-blue.svg)

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
    mux.Route("/").Get(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello World!"))
    })

    http.ListenAndServe(":5000", mux)
}
```

## Install

```
go get github.com/thisissoon/yam
```

## Documentaion

Full package documentation can be found at https://godoc.org/github.com/thisissoon/yam.
