package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler

	h := NoSurf(&myH)
	if h == nil {
		t.Errorf("NoSurf returned nil handler")
	}

	switch v := h.(type) {
	case http.Handler:
		// ok
	default:
		t.Errorf("NoSurf returned a non http.Handler type: %T", v)
	}
}
