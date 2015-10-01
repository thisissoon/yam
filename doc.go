// Copyright 2015 SOON_ London Limited. All rights reserved.
// Use of this source code is governed by The MIT License (MIT).
// This can be found in the LICENSE file at the repository root.

/*
YAM (Yet Another Mux) is a simple HTTP Multiplexer (Router).

Overview

YAM's goal is to be simple, flexible and configurable. Above all
it does break the standard interfaces and function handlers of
the net/http package.

This is YAM's most simplest implementation:

	mux := yam.New()
	mux.Route("/").Get(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})
	http.ListenAndServe(":5000", mux)

Features

Yam has the following features:

	- Method based routing, returning 405 when a route does not implement a specific verb
	- Simple pattern matching using the "/foo/:bar/baz" syntax, values are placed onto the request URL parameters
	- Support for all the standard HTTP verbs out of the box (OPTIONS, GET, HEAD, POST, PUT, PATCH, DELETE, TRACE)
	- Sub Routing
	- Configuration, allowing default handler functions overrides and flags for OPTIONS and TRACE

Method Based Routing

To implement a method on a route simple call the routes method function for that method:

	mux := yam.New()
	mux.Route("/").Get(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Get Request"))
	})
	mux.Route("/").Post(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Post Request"))
	})
	http.ListenAndServe(":5000", mux)

In the above example only "GET" and "POST" are supported, other methods such as "PUT"
would return a "405 Method Not Allowed".

Pattern Matching

YAM implements a very simple "/foo/:bar" pattern matching system, values from those patterns
are placed on the requests URL as query parameters and therefore no extra dependencies are
required. The values persist down the path.

	mux := yam.New()
	mux.Route("/foo/:bar").Get(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.Query().Get(":bar")))
	})
	mux.Route("/foo/:bar/baz/:fiz").Get(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.Query().Get(":bar")+"\n"))
		w.Write([]byte(r.URL.Query().Get(":fix")+"\n"))
	})

	http.ListenAndServe(":5000", mux)

Methods

YAM supports all the standard HTTP verbs, with the exception of CONNECT (https://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol#Request_methods).
Each verb has it's own function on a specific route, if the route does not implement a verb and a request
is made with that verb a 405 will be returned.

	mux := yam.New()
	mux.Route("/foo").
		Get(get).
		Post(post).
		Put(put).
		Patch(patch).
		Delete(delete)
	http.ListenAndServe(":5000", mux)

In the case of "OPTIONS" and "HEAD" these are handled for you provided YAM is configured to do (which it is by default).

The default "OPTIONS" handler function will add an "Allow" header which will be populated with the verbs
supported by the route. This behaviour can be turned off or overridden on a route by route basis:

	mux := yam.New()
	mux.Route("/foo").Get(myGetHandler).Options(myOptionsHandler)


The default "HEAD" handler will only be implemented for the route if a Get handler has been added. This will simply
return the Headers for a response minus a response body, this behaviour can be turned off or overridden for specific
routes:

	mux := yam.New()
	mux.Route("/foo").Get(myGetHandler).Head(myHeadHandler)

The "TRACE" verb is disabled by default but is useful for debugging purposes, to enable it alter YAM's configuration:

	mux := yam.New()
	mux.Config.Trace = true
	mux.Route("/foo").Get(myGetHandler)

You can also enable "TRACE" on a route by route basis:

	mux := yam.New()
	mux.Route("/foo").Trace(yam.DefaultTraceHandler)
	mux.Route("/bar").Trace(myTraceHandler)

You can also use the "Add" function:

	mux := yam.New()
	mux.Route("/foo").Add("VERB", http.Handler)

Sub Routing

Routes can also be broken up to avoid repetition and broken up across packages.

	mux := yam.New()
	foo := mux.Route("/foo").Get(myGetHandler)
	bar := foo.Route("/bar").Post(myPostHandler).Delete(myDeleteHandler)
	bar.Route("/baz").Put(myPutHandler)

The above would register the following routes.

	GET /foo
	POST & DELETE /foo/bar
	PUT /foo/bar/baz

Configuration

Finally if you do not like any of the default settings of "YAM", you can change them! The Config type
allows this:

	mux := yam.New()
	config := NewConfig() // Defaults
	config.Options = false // OPTIONS will no longer be supported on all routes
	config.OptionsHandler = func(r *yam.Route) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("My-Custom-Header", "Foo")
		})
	} // Set a custom OPTIONS handler used when config.Options is true
	config.Trace = true // TRACE will now be supported on all routes
	config.TraceHandler = func(r *yam.Route) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dump, _ := httputil.DumpRequest(r, false)
			w.Write(dump)
		})
	} // Set a custom handler function for TRACE when config.Trace is true
	config.AddHeadOnGet = false // HEAD support will not longer be added by default on Get
	mux.Config = config

*/
package yam
