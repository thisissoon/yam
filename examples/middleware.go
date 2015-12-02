package main

import (
	"net/http"

	"github.com/thisissoon/yam"
)

func GlobalMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("GobalMiddlware", "True")

		// Serve the next handler
		next.ServeHTTP(w, r)
	})
}

func RouteOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("RouteOnly", "True")

		// Serve the next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := yam.New()
	mux.Route("/").Get(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	}))

	mux.Route("/foo").Get(RouteOnly(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})))

	http.ListenAndServe(":5000", GlobalMiddleware(mux))
}
