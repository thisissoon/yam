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
