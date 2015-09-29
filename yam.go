// Yam (Yet Another Mux)

package yam

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// An easy type to reference standard http handler functions
type handler func(http.ResponseWriter, *http.Request)

// Configuration type that allows configuration of YAM
type Config struct {
	Options        bool
	OptionsHandler func(*Route) http.Handler
	Trace          bool
	TraceHandler   func(*Route) http.Handler
	AddHeadOnGet   bool
	HeadHandler    func(http.Handler) http.Handler
}

// Constructs a new Config instance with default values
func NewConfig() *Config {
	return &Config{
		Options:        true,
		OptionsHandler: DefaultOptionsHandler,
		Trace:          false,
		TraceHandler:   DefaultTraceHandler,
		AddHeadOnGet:   true,
		HeadHandler:    DefaultHeadHandler,
	}
}

// Base YAM type contain the Root Routes and Configuration
type Yam struct {
	Root   *Route
	Config *Config
}

// Constructs a new YAM instance with default configuration
func New() *Yam {
	y := &Yam{}
	y.Root = &Route{yam: y}
	y.Config = NewConfig()

	return y
}

// Creates a new base Route - Effectively a constructor for Route
func (y *Yam) Route(path string) *Route {
	route := y.Root.Route(path)
	route.yam = y

	return route
}

// Registers a new Route on the Path, building a Tree structure of Routes
func route(path string, router *Route, y *Yam) *Route {
	parts := strings.Split(path, "/")[1:]
	routes := router.Routes
	fullPath := router.path + path

	for i, part := range parts {
		if i == len(parts)-1 {
			for _, route := range routes {
				if route.leaf == part {
					return route
				}
			}
			route := &Route{leaf: part, path: fullPath, yam: y}
			router.Routes = append(router.Routes, route)
			return route
		} else {
			for _, route := range routes {
				if route.leaf == part {
					router = route
				} else {
					route := &Route{leaf: part, path: router.path + path, yam: y}
					router.Routes = append(router.Routes, route)
					router = route
				}
			}
		}
	}

	return nil
}

// Implements the http.Handler Interface.  Finds the correct handler for
// a path based on the path and http verb of the request.
func (y *Yam) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")[1:]
	routes := y.Root.Routes
	for i, part := range parts {
		for _, route := range routes {
			match := false
			// Pattern Match
			if strings.HasPrefix(route.leaf, ":") {
				match = true
				values := url.Values{}
				values.Add(route.leaf, part)
				r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery
			} else { // Exact match
				if route.leaf == part {
					match = true
				}
			}
			// Did we get a match
			if match {
				// If we are not at the end of the path, then we go round again
				if i < len(parts)-1 {
					routes = route.Routes
					break
				} else { // We are at the end of the path - we have a route
					handler := route.handlers[r.Method]
					// Do we have a handler for this Verb
					if handler != nil {
						// Yes, serve and return
						handler.ServeHTTP(w, r)
						return
					}
					// We do not - Serve a 405 Method Not Allowed
					w.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
			}
		}
	}

	// If we get here then we have not found a route
	w.WriteHeader(http.StatusNotFound)
}

// This type contains all the handlers for each path, each Route can also hold
// a list of child routes
type Route struct {
	leaf   string   // a part of a URL path, /foo/bar - a leaf would be foo and bar
	path   string   // full url path
	Routes []*Route // Routes that live under this route

	yam *Yam // Reference to Yam and global configuration

	// Verb handlers
	handlers map[string]http.Handler
}

// Adds a new route to the tree, and depending on configuration implements
// default handler implementation for OPTIONS and TRACE requests
func (r *Route) Route(path string) *Route {
	r = route(path, r, r.yam)

	if r.handlers == nil {
		r.handlers = make(map[string]http.Handler)
	}

	if r.yam.Config.Options {
		if r.handlers["OPTIONS"] == nil {
			r.handlers["OPTIONS"] = r.yam.Config.OptionsHandler(r)
		}
	}

	if r.yam.Config.Trace {
		if r.handlers["TRACE"] == nil {
			r.handlers["TRACE"] = r.yam.Config.TraceHandler(r)
		}
	}

	return r
}

// Adds a new handler to the route based on http Verb
func (r *Route) Add(method string, handler http.Handler) *Route {
	r.handlers[method] = handler

	return r
}

// Set a HEAD request handler for the route
func (r *Route) Head(h handler) *Route {
	r.Add("HEAD", http.HandlerFunc(h))

	return r
}

// Set a OPTIONS request for the route, overrides default implementation
func (r *Route) Options(h handler) *Route {
	r.Add("OPTIONS", http.HandlerFunc(h))

	return r
}

// Set a TRACE request for the route, overrides default implementation
func (r *Route) Trace(h handler) *Route {
	r.Add("TRACE", http.HandlerFunc(h))

	return r
}

// Set a GET request handler for the Route. A default HEAD request implementation
// will also be implemented since HEAD requests should perform the same as a GET
// request but simply not return the response body.
func (r *Route) Get(h handler) *Route {
	r.Add("GET", http.HandlerFunc(h))

	if r.yam.Config.AddHeadOnGet {
		// Apply the head middleware to the head handler
		r.Add("HEAD", r.yam.Config.HeadHandler(http.HandlerFunc(h)))
	}

	return r
}

// Set a POST request handler for the route.
func (r *Route) Post(h handler) *Route {
	r.Add("POST", http.HandlerFunc(h))

	return r
}

// Set a PUT request handler for the route.
func (r *Route) Put(h handler) *Route {
	r.Add("PUT", http.HandlerFunc(h))

	return r
}

// Set a DELETE request handler for the route.
func (r *Route) Delete(h handler) *Route {
	r.Add("DELETE", http.HandlerFunc(h))

	return r
}

// Set a PATCH request handler for the route.
func (r *Route) Patch(h handler) *Route {
	r.Add("PATCH", http.HandlerFunc(h))

	return r
}

// Default HTTP Handler function for OPTIONS requests. Adds the Allow header
// with the http verbs the route supports
func DefaultOptionsHandler(route *Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		methods := []string{}
		for key, _ := range route.handlers {
			methods = append(methods, key)
		}
		w.Header().Add("Allow", strings.Join(methods, ", "))
	})
}

// Default HTTP handler function for TRACE requests. Dumps the request as the Response.
func DefaultTraceHandler(route *Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dump, _ := httputil.DumpRequest(r, false)
		w.Write(dump)
	})
}

// Default HEAD Request Handler. Automatically added to GET requests.
func DefaultHeadHandler(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve the handler
		n.ServeHTTP(w, r)
		// Flush the body so we don't write to the client
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	})
}
