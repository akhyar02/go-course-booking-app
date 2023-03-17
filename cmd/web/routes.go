package main

import (
	"net/http"

	"github.com/akhyar02/bookings/internal/config"
	"github.com/akhyar02/bookings/internal/handlers"
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

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact-us", handlers.Repo.ContactUs)
	mux.Get("/general-quarters", handlers.Repo.GeneralQuarters)
	mux.Get("/major-suites", handlers.Repo.MajorSuites)
	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Get("/make-reservation", handlers.Repo.MakeReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	mux.Route("/api", func(r chi.Router) {
		r.Get("/reservations", handlers.ReservationApi.GetReservationByDate)
		r.Post("/reservations", handlers.ReservationApi.CreateReservation)
	})

	return mux
}
