// Copyright 2015 SOON_ London Limited. All rights reserved.
// Use of this source code is governed by The MIT License (MIT).
// This can be found in the LICENSE file at the repository root.

package yam

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// Configuration type that allows configuration of YAM
type Config struct {
	Options        bool
	OptionsHandler func(*Route) http.Handler
	Trace          bool
	TraceHandler   func(*Route) http.Handler
	AddHeadOnGet   bool
}

// Constructs a new Config instance with default values
func NewConfig() *Config {
	return &Config{
		Options:        true,
		OptionsHandler: DefaultOptionsHandler,
		Trace:          false,
		TraceHandler:   DefaultTraceHandler,
		AddHeadOnGet:   true,
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
	y.Config = NewConfig()
	y.Root = &Route{yam: y}

	return y
}

// Creates a new base Route - Effectively a constructor for Route
func (y *Yam) Route(path string) *Route {
	route := y.Root.Route(path)

	return route
}

// Gets or Creates the route for the path. As it traverses the tree routes
// are either created if they do not exist. At the end the function returns
// the last leaf of the tree
func route(path string, base *Route, y *Yam) *Route {
	var route *Route
	// /foo/bar/baz [foo bar baz]
	parts := strings.Split(path, "/")[1:]
	// Our starting routes we loop over should be the base route
	route = base

	// Iterate over the parts
	var found bool
	for _, part := range parts {
		// Iterate over the routes on
		found = false
		for _, r := range route.Routes {
			// This part of the path already exists in the routes
			if r.leaf == part {
				// Set our base route to now be this route
				route = r
				found = true
				break
			}
		}
		if !found {
			// The part of the path does not exist in the routes, create it
			r := &Route{leaf: part, yam: y}
			// Add the route to the list of routes
			route.Routes = append(route.Routes, r)
			// Set the next route to be the one we just created
			route = r
		}
	}

	return route
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
func (r *Route) Add(method string, f http.HandlerFunc) *Route {
	r.handlers[method] = http.HandlerFunc(f)

	return r
}

// Set a HEAD request handler for the route
func (r *Route) Head(f http.HandlerFunc) *Route {
	r.Add("HEAD", f)

	return r
}

// Set a OPTIONS request for the route, overrides default implementation
func (r *Route) Options(f http.HandlerFunc) *Route {
	r.Add("OPTIONS", f)

	return r
}

// Set a TRACE request for the route, overrides default implementation
func (r *Route) Trace(f http.HandlerFunc) *Route {
	r.Add("TRACE", f)

	return r
}

// Set a GET request handler for the Route. A default HEAD request implementation
// will also be implemented since HEAD requests should perform the same as a GET
// request but simply not return the response body.
func (r *Route) Get(f http.HandlerFunc) *Route {
	r.Add("GET", f)

	if r.yam.Config.AddHeadOnGet {
		// Apply the head middleware to the head handler
		r.Add("HEAD", f)
	}

	return r
}

// Set a POST request handler for the route.
func (r *Route) Post(f http.HandlerFunc) *Route {
	r.Add("POST", f)

	return r
}

// Set a PUT request handler for the route.
func (r *Route) Put(f http.HandlerFunc) *Route {
	r.Add("PUT", f)

	return r
}

// Set a DELETE request handler for the route.
func (r *Route) Delete(f http.HandlerFunc) *Route {
	r.Add("DELETE", f)

	return r
}

// Set a PATCH request handler for the route.
func (r *Route) Patch(f http.HandlerFunc) *Route {
	r.Add("PATCH", f)

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
