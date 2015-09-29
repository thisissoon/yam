// Yam (Yet Another Mux)

package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type handler func(http.ResponseWriter, *http.Request)

func optionsHandler(route *Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methods := []string{}
		for key, _ := range route.handlers {
			methods = append(methods, key)
		}
		w.Header().Add("Allow", strings.Join(methods, ", "))
	})
}

type Yam struct {
	Root *Route
}

func New() *Yam {
	return &Yam{
		Root: &Route{},
	}
}

func (y *Yam) Route(path string) *Route {
	r := route(path, y.Root)

	if r.handlers == nil {
		r.handlers = make(map[string]http.Handler)
	}

	if r.handlers["OPTIONS"] == nil {
		r.handlers["OPTIONS"] = optionsHandler(r)
	}

	if r.handlers["TRACE"] == nil {
		r.handlers["TRACE"] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dump, _ := httputil.DumpRequest(r, false)
			w.Write(dump)
		})
	}

	return r
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

func (y *Yam) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")[1:]
	fmt.Println(parts)
	routes := y.Root.Routes

	for i, part := range parts {
		fmt.Println(part)
		for _, route := range routes {
			fmt.Println("Leaf:", route.leaf)
			match := false
			// Pattern Match
			if strings.HasPrefix(route.leaf, ":") {
				fmt.Println("Pattern Match")
				match = true
				values := url.Values{}
				values.Add(route.leaf, part)
				r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery
			} else { // Exact match
				fmt.Println("Exact Match")
				if route.leaf == part {
					match = true
				}
			}

			if match {
				fmt.Println("Leaf ==", part)
				if i < len(parts)-1 {
					routes = route.Routes
					break
				} else {
					fmt.Println("Found: ", route.path)

					handler := route.handlers[r.Method]
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

	handlers map[string]http.Handler
}

func (r *Route) Route(path string) *Route {
	r = route(path, r)

	if r.handlers == nil {
		r.handlers = make(map[string]http.Handler)
	}

	if r.handlers["OPTIONS"] == nil {
		r.handlers["OPTIONS"] = optionsHandler(r)
	}

	if r.handlers["TRACE"] == nil {
		r.handlers["TRACE"] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dump, _ := httputil.DumpRequest(r, false)
			w.Write(dump)
		})
	}

	return r
}

func (r *Route) Add(method string, handler http.Handler) *Route {
	r.handlers[method] = handler

	return r
}

func (r *Route) Head(h handler) *Route {
	r.Add("HEAD", http.HandlerFunc(h))

	return r
}

func (r *Route) Get(h handler) *Route {
	r.Add("GET", http.HandlerFunc(h))

	// Implement the HEAD handler by default for all GET requests - HEAD
	// should not return a body so we wrap it in a middleware
	head := func(n http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Serve the handler
			n.ServeHTTP(w, r)
			// Flush the body so we don't write to the client
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		})
	}

	// Apply the head middleware to the head handler
	r.Add("HEAD", head(http.HandlerFunc(h)))

	return r
}

func (r *Route) Post(h handler) *Route {
	r.Add("POST", http.HandlerFunc(h))

	return r
}

func (r *Route) Put(h handler) *Route {
	r.Add("PUT", http.HandlerFunc(h))

	return r
}

func (r *Route) Delete(h handler) *Route {
	r.Add("DELETE", http.HandlerFunc(h))

	return r
}

func (r *Route) Patch(h handler) *Route {
	r.Add("PATCH", http.HandlerFunc(h))

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
	y := New()

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

	fmt.Printf("%+v\n", y)

	http.ListenAndServe(":5000", y)
}
