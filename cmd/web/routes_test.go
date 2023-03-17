package main

import (
	"testing"

	"github.com/akhyar02/bookings/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	app := &config.AppConfig{}

	mux := routes(app)

	switch v := mux.(type) {
	case *chi.Mux:
	default:
		t.Errorf("expected *Chi.Mux, received %T", v)
	}
}
