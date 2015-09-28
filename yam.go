// Yam (Yet Another Mux)

package main

import (
	"fmt"
	"net/http"
	"strings"
)

type handler func(http.ResponseWriter, *http.Request)

type Yam struct {
	Root *Route
}

func NewYam() *Yam {
	return &Yam{
		Root: &Route{},
	}
}

func (y *Yam) Route(path string) *Route {
	parts := strings.Split(path, "/")[1:]
	routes := y.Root.Routes
	router := y.Root

	for i, part := range parts {
		fmt.Println(part)
		if i == len(parts)-1 {
			for _, route := range routes {
				if route.leaf == part {
					fmt.Println("Route Exists")
					return route
				}
			}

			route := &Route{leaf: part}
			router.Routes = append(router.Routes, route)

			return route
		} else {
			fmt.Println("Not Last Part")
			for _, route := range routes {
				if route.leaf == part {
					router = route
				} else {
					route := &Route{leaf: part}
					router.Routes = append(router.Routes, route)
					router = route
				}
			}
		}
	}

	return nil
}

func (y *Yam) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")[1:]

	routes := y.Root.Routes
	for i, part := range parts {
		fmt.Println(part)
		for _, route := range routes {
			fmt.Println(route.leaf)
			if route.leaf == part {
				fmt.Println("Found Leaf")
				if i < len(parts)-1 {
					routes = route.Routes
					break
				} else {
					fmt.Println("Found")

					var handler http.Handler
					switch r.Method {
					case "GET":
						handler = route.GetHandler
					}

					if handler != nil {
						handler.ServeHTTP(w, r)
						return
					}

					fmt.Println("No handler for method")
					w.WriteHeader(http.StatusMethodNotAllowed)
					return

				}
			}
		}
	}

	// If we get here then we have not found a route
	fmt.Println("Not Found")
	w.WriteHeader(http.StatusNotFound)
}

type Route struct {
	leaf   string   // a part of a URL path, /foo/bar - a leaf would be foo and bar
	Routes []*Route // Sub Routes

	GetHandler http.Handler
}

func (r *Route) Get(h handler) *Route {
	r.GetHandler = http.HandlerFunc(h)

	return r
}

func GetRootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Root Handler"))
}

func GetFooHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get Foo Handler"))
}

func GetAHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get A Handler"))
}

func GetBHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get B Handler"))
}

func main() {
	y := NewYam()

	y.Root.Routes = append(y.Root.Routes, &Route{
		leaf:   "foo",
		Routes: make([]*Route, 0),

		GetHandler: http.HandlerFunc(GetFooHandler),
	})

	y.Root.Routes[0].Routes = append(y.Root.Routes[0].Routes, &Route{
		leaf: "bar",
	})

	y.Root.Routes = append(y.Root.Routes, &Route{
		leaf:   "baz",
		Routes: make([]*Route, 0),
	})

	y.Route("/").Get(GetRootHandler)
	y.Route("/foo")
	y.Route("/a").Get(GetAHandler)
	y.Route("/a/b").Get(GetBHandler)

	fmt.Printf("%+v\n", y)

	http.ListenAndServe(":5000", y)
}
