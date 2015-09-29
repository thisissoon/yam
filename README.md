# YAM

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

## More Examples

``` go
y := New()
y.Config.Trace = true

y.Route("/").Get(GetRootHandler)
y.Route("/get").Get(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("GET"))
})
y.Route("/post").Post(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("POST"))
})
y.Route("/put").Put(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("PUT"))
})
y.Route("/patch").Patch(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("PATCH"))
})
y.Route("/delete").Delete(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("DELETE"))
})

a := y.Route("/a").Get(GetAHandler)
a.Route("/b").Get(GetBHandler)
a.Route("/b").Put(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("PUT B Handler"))
})
c := a.Route("/b/c").Get(GetCHandler)
c.Route("/d").Get(GetDHandler)
e := c.Route("/d/e").Get(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("E Handler"))
})
e.Route("/f").Get(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("F Handler"))
})

// Pattern Matching
a.Route("/:foo").Get(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("A :foo Handler\n"))
    w.Write([]byte(r.URL.Query().Get(":foo")))
})

bar := a.Route("/:foo/:bar").Get(func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("/a/:foo/:bar Handler\n"))
    w.Write([]byte(r.URL.Query().Get(":foo")))
    w.Write([]byte("\n"))
    w.Write([]byte(r.URL.Query().Get(":bar")))
})

bar.Route("/baz").Get(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Add("Foo", "Bar")
    w.Write([]byte("baz\n"))
    w.Write([]byte(r.URL.Query().Get(":foo")))
    w.Write([]byte("\n"))
    w.Write([]byte(r.URL.Query().Get(":bar")))
})
```
