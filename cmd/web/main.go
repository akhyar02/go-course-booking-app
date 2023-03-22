package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/akhyar02/bookings/internal/config"
	"github.com/akhyar02/bookings/internal/driver"
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
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

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

	dbName := flag.String("DB_NAME", "", "database name")
	dbUser := flag.String("DB_USER", "", "database user")
	dbPassowrd := flag.String("DB_PASSWORD", "", "database password")
	dbHost := flag.String("DB_HOST", "", "database host")
	dbPort := flag.String("DB_PORT", "", "database port")
	flag.Parse()
	dsn := "host=" + *dbHost + " port=" + *dbPort + " dbname=" + *dbName + " user=" + *dbUser + " password=" + *dbPassowrd
	db, err := driver.ConnectSQL(dsn)
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	defer db.SQL.Close()

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Error creating template cache:", err)
	}
	app.TemplateCache = tc
	config.TestConf.TemplateCache = tc
	repo := handlers.NewRepository(&app, db)
	helpers.NewHelpers(&app)
	handlers.NewHandlers(repo)
	handlers.NewReservationHandler(app, db)

	render.NewRenderer(&app)

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
