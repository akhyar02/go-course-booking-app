package main

import (
	"net/http"

	"github.com/akhyar02/bookings/pkg/config"
	"github.com/akhyar02/bookings/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	// ? pat
	// mux := pat.New()

	// mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	// mux.Get("/about", http.HandlerFunc(handlers.Repo.About))

	// ? chi
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	mux.Use(writeToConsole)
	mux.Use(noSurf)
	mux.Use(sessionLoad)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	mux.Get("/contact-us", http.HandlerFunc(handlers.Repo.ContactUs))
	mux.Get("/general-quarters", http.HandlerFunc(handlers.Repo.GeneralQuarters))
	mux.Get("/major-suites", http.HandlerFunc(handlers.Repo.MajorSuites))
	mux.Get("/search-availability", http.HandlerFunc(handlers.Repo.SearchAvailability))
	mux.Get("/make-reservation", http.HandlerFunc(handlers.Repo.MakeReservation))

	return mux
}
