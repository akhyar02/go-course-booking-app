package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/akhyar02/bookings/internal/config"
	"github.com/akhyar02/bookings/internal/handlers"
	"github.com/akhyar02/bookings/internal/helpers"
	"github.com/akhyar02/bookings/internal/models"
	"github.com/akhyar02/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = 8080

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	app.InProduction = false
	gob.Register(models.Reservation{})

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
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

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Error creating template cache:", err)
	}
	app.TemplateCache = tc
	config.TestConf.TemplateCache = tc
	repo := handlers.NewRepository(&app)
	helpers.NewHelpers(&app)
	handlers.NewHandlers(repo)
	handlers.NewReservationHandler(app)

	render.NewTemplates(&app)

	// ? using pat
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", portNumber),
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error: server not started", err)
	}
}
