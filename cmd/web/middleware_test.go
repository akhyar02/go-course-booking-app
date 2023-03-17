package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	handler := &mockHandler{}
	h := noSurf(handler)

	switch v := h.(type) {
	case http.Handler:
		// do nothing - test passed
	default:
		t.Errorf("type is not http.Handler, got %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	handler := &mockHandler{}
	h := sessionLoad(handler)

	switch v := h.(type) {
	case http.Handler:
		// do nothing - test passed
	default:
		t.Errorf("type is not http.Handler, got %T", v)
	}
}
