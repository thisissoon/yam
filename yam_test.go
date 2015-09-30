// Unittests for Yam

package yam

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestRoute struct {
	Path    string
	Methods []string
	Handler func(http.ResponseWriter, *http.Request)
}

type TestRequest struct {
	Path   string
	Method string
}

type TestResponse struct {
	Status int
	Body   []byte
}

// Table Driven Tests
var tests = []struct {
	// Request to make
	route    TestRoute
	req      TestRequest
	response TestResponse
}{
	// Simplest Route
	{
		TestRoute{"/foo", []string{"GET"}, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(r.URL.Path))
		}},
		TestRequest{"/foo", "GET"},
		TestResponse{http.StatusOK, []byte("/foo")},
	},
	// 404 Handling
	{
		TestRoute{"/foo", []string{"GET"}, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(r.URL.Path))
		}},
		TestRequest{"/bar", "GET"},
		TestResponse{http.StatusNotFound, nil},
	},
	// 405 Handling
	{
		TestRoute{"/foo", []string{"GET"}, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(r.URL.Path))
		}},
		TestRequest{"/foo", "POST"},
		TestResponse{http.StatusMethodNotAllowed, nil},
	},
	// Pattern Matching & Added to Query
	{
		TestRoute{"/foo/:bar", []string{"GET"}, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(r.URL.Query().Get(":bar")))
		}},
		TestRequest{"/foo/bar", "GET"},
		TestResponse{http.StatusOK, []byte("bar")},
	},
}

func TestTables(t *testing.T) {
	for _, test := range tests {
		y := New()
		r := y.Route(test.route.Path)
		for _, method := range test.route.Methods {
			r.Add(method, http.HandlerFunc(test.route.Handler))
		}

		s := httptest.NewServer(y)
		req, _ := http.NewRequest(test.req.Method, s.URL+test.req.Path, nil)
		c := &http.Client{}
		res, _ := c.Do(req)

		if res.StatusCode != test.response.Status {
			t.Errorf("Status was %v, should be %v", res.StatusCode, test.response.Status)
		}

		body, _ := ioutil.ReadAll(res.Body)
		if !bytes.Equal(body, test.response.Body) {
			t.Errorf("Body was %v, should be %v", string(body[:]), string(test.response.Body[:]))

		}

		s.Close()
	}
}
