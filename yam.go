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

func route(path string, router *Route) *Route {
	parts := strings.Split(path, "/")[1:]
	routes := router.Routes

	fmt.Println("Start Router:", router.path)
	fmt.Println("Stat Path:", path)
	fullPath := router.path + path

	for i, part := range parts {
		fmt.Println("Part:", part)
		if i == len(parts)-1 {

			for _, route := range routes {
				if route.leaf == part {
					fmt.Println("Route Exists")
					fmt.Println("--------------")
					return route
				}
			}

			route := &Route{leaf: part, path: fullPath}
			fmt.Println("Add:", route.path)
			fmt.Println("Router:", router.path)
			router.Routes = append(router.Routes, route)

			fmt.Println("--------------")

			return route

		} else {
			for _, route := range routes {
				if route.leaf == part {
					fmt.Println("Leaf:", route.leaf)
					router = route
				} else {
					fmt.Println("Router:", router.path)
					route := &Route{leaf: part, path: router.path + path}
					router.Routes = append(router.Routes, route)
					router = route
				}
			}

		}
	}

	return nil
}

func (y *Yam) Route(path string) *Route {
	return route(path, y.Root)
}

func (y *Yam) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")[1:]
	fmt.Println(parts)

	routes := y.Root.Routes
	for i, part := range parts {
		fmt.Println(part)
		for _, route := range routes {
			fmt.Println("Leaf:", route.leaf)
			if route.leaf == part {
				fmt.Println("Leaf ==", part)
				if i < len(parts)-1 {
					routes = route.Routes
					break
				} else {
					fmt.Println("Found: ", route.path)

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
	path   string   // full url path
	Routes []*Route // Routes that live under this route

	GetHandler http.Handler
}

func (r *Route) Route(path string) *Route {
	return route(path, r)
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

func GetCHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get C Handler"))
}

func GetDHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get D Handler"))
}

func main() {
	y := NewYam()

	y.Route("/").Get(GetRootHandler)
	y.Route("/foo")

	a := y.Route("/a").Get(GetAHandler)
	a.Route("/b").Get(GetBHandler)
	c := a.Route("/c").Get(GetCHandler)
	c.Route("/d").Get(GetDHandler)
	c.Route("/d/e").Get(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("E Handler"))
	})

	fmt.Printf("%+v\n", y)

	http.ListenAndServe(":5000", y)
}
