package handlers

import (
	"log"
	"net/http"

	"github.com/akhyar02/bookings/internal/config"
	"github.com/akhyar02/bookings/internal/driver"
	"github.com/akhyar02/bookings/internal/models"
	"github.com/akhyar02/bookings/internal/render"
	"github.com/akhyar02/bookings/internal/repository"
	"github.com/akhyar02/bookings/internal/repository/dbrepo"
)

type Repository struct {
	app *config.AppConfig
	db  repository.DatabaseRepo
}

var Repo *Repository

func NewRepository(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		app: a,
		db:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (rp *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "home.page.gtpl", &models.TemplateData{})
}

// About is the about page handler
func (rp *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again"
	render.Template(w, r, "about.page.gtpl", &models.TemplateData{StringMap: stringMap})
}

// GeneralQuarters is the general-quarters page handler
func (rp *Repository) GeneralQuarters(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.gtpl", &models.TemplateData{})
}

// MajorSuites is the major-suites page handler
func (rp *Repository) MajorSuites(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.gtpl", &models.TemplateData{})
}

// MakeReservation is the make-reservation page handler
func (rp *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	var urlQuery = r.URL.Query()
	var roomType = urlQuery.Get("roomType")
	var startDate = urlQuery.Get("startDate")
	var endDate = urlQuery.Get("endDate")
	render.Template(w, r, "make-reservation.page.gtpl", &models.TemplateData{
		Data: map[string]interface{}{
			"roomType":  roomType,
			"startDate": startDate,
			"endDate":   endDate,
		},
	})
}

// SearchAvailability is the search-availability page handler
func (rp *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	var urlQuery = r.URL.Query()
	var roomType = urlQuery.Get("room")
	render.Template(w, r, "search-availability.page.gtpl", &models.TemplateData{
		Data: map[string]interface{}{
			"roomType": roomType,
		},
	})
}

// ContactUs is the contact-us page handler
func (rp *Repository) ContactUs(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact-us.page.gtpl", &models.TemplateData{})
}

// ReservationSummary is the reservation-summary page handler
func (rp *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := rp.app.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("Type error")
		http.Redirect(w, r, "/make-reservation", http.StatusTemporaryRedirect)
		return
	}
	render.Template(w, r, "reservation-summary.page.gtpl", &models.TemplateData{
		Data: map[string]interface{}{
			"roomType":  reservation.Room.RoomName,
			"startDate": reservation.StartDate,
			"endDate":   reservation.EndDate,
			"firstName": reservation.FirstName,
			"lastName":  reservation.LastName,
			"email":     reservation.Email,
			"phone":     reservation.Phone,
		},
	})
}
