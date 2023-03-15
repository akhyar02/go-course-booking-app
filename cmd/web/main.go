package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/akhyar02/bookings/pkg/config"
	"github.com/akhyar02/bookings/pkg/handlers"
	"github.com/akhyar02/bookings/pkg/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = 8080

var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.InProduction = false

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
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	// ? using pat
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", portNumber),
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error: server not started")
	}
}
