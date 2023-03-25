package handlers

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/akhyar02/bookings/internal/config"
	"github.com/akhyar02/bookings/internal/driver"
	"github.com/akhyar02/bookings/internal/helpers"
	"github.com/akhyar02/bookings/internal/models"
	"github.com/akhyar02/bookings/internal/render"
	_ "github.com/akhyar02/bookings/testing-setup"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"gopkg.in/gomail.v2"
)

var app config.AppConfig
var session *scs.SessionManager
var db *driver.DB

func TestMain(m *testing.M) {
	defer func() {
		if db.SQL != nil {
			db.SQL.Close()
		}
		close(app.MailChan)
	}()
	os.Exit(m.Run())
}

func getRoutes() http.Handler {
	app.InProduction = false
	app.UseCache = false
	gob.Register(models.Reservation{})
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog
	app.InfoLog = infoLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.HttpOnly = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction
	session.Cookie.Path = "/"
	app.Session = session

	app.MailChan = make(chan models.MailData)
	listenForMail()

	var err error
	db, err = driver.ConnectSQL("host=127.0.0.1 port=5432 dbname=booking user=postgres password=postgres")
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	// defer db.SQL.Close()

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Error creating template cache:", err)
	}
	app.TemplateCache = tc
	config.TestConf.TemplateCache = tc
	repo := NewRepository(&app, db)
	helpers.NewHelpers(&app)
	NewHandlers(repo)
	NewReservationHandler(app, db)
	render.NewRenderer(&app)

	// ? pat
	// mux := pat.New()

	// mux.Get("/", http.HandlerFunc(Repo.Home))
	// mux.Get("/about", http.HandlerFunc(Repo.About))

	// ? chi
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)
	// mux.Use(writeToConsole)
	// mux.Use(noSurf)
	mux.Use(sessionLoad)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/contact-us", Repo.ContactUs)
	mux.Get("/general-quarters", Repo.GeneralQuarters)
	mux.Get("/major-suites", Repo.MajorSuites)
	mux.Get("/search-availability", Repo.SearchAvailability)
	mux.Get("/make-reservation", Repo.MakeReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Route("/api", func(r chi.Router) {
		r.Get("/reservations", ReservationApi.GetReservationByDate)
		r.Post("/reservations", ReservationApi.CreateReservation)
	})

	return mux
}

func noSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

func sessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func listenForMail() {
	go func() {
		server := gomail.NewDialer("localhost", 587, "", "")
		var s gomail.SendCloser
		open := false
		for {
			select {
			case m, ok := <-app.MailChan:
				if !ok {
					return
				}
				sendMail(m, server, &s, &open)
			case <-time.After(30 * time.Second):
				if !open {
					continue
				}
				s.Close()
				open = false
			}
		}
	}()
}

func sendMail(m models.MailData, server *gomail.Dialer, s *gomail.SendCloser, open *bool) {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.From)
	msg.SetHeader("To", m.To...)
	msg.SetHeader("Subject", m.Subject)
	msg.SetBody("text/html", string(m.Content))

	if !*open {
		var err error
		*s, err = server.Dial()
		if err != nil {
			log.Println(err)
			return
		} else {
			*open = true
		}
	}

	if err := gomail.Send(*s, msg); err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent")
	}

}
