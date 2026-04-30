package main

import (
	"net/http"
	"os"
	"testing"
)

type myHandler struct {
}

func TestMain(m *testing.M) {
	// run the main function
	os.Exit(m.Run())
	// run the tests

}
func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// handle the request
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, world!"))
}
