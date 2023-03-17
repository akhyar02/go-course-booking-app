package main

import (
	"net/http"
	"os"
	"testing"

	_ "github.com/akhyar02/bookings/testing-setup"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

type mockHandler struct{}

func (mh *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {}
