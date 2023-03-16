package handlers

import (
	"net/http"

	"github.com/akhyar02/bookings/internal/config"
	"github.com/akhyar02/bookings/internal/models"
	"github.com/akhyar02/bookings/internal/render"
)

type Repository struct {
	app *config.AppConfig
}

var Repo *Repository

func NewRepository(a *config.AppConfig) *Repository {
	return &Repository{
		app: a,
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (rp *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.gtpl", &models.TemplateData{})
}

// About is the about page handler
func (rp *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again"
	render.RenderTemplate(w, r, "about.page.gtpl", &models.TemplateData{StringMap: stringMap})
}

// GeneralQuarters is the general-quarters page handler
func (rp *Repository) GeneralQuarters(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.gtpl", &models.TemplateData{})
}

// MajorSuites is the major-suites page handler
func (rp *Repository) MajorSuites(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.gtpl", &models.TemplateData{})
}

// MakeReservation is the make-reservation page handler
func (rp *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "make-reservation.page.gtpl", &models.TemplateData{})
}

// SearchAvailability is the search-availability page handler
func (rp *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	var urlQuery = r.URL.Query()
	var roomType = urlQuery.Get("room")
	render.RenderTemplate(w, r, "search-availability.page.gtpl", &models.TemplateData{
		Data: map[string]interface{}{
			"roomType": roomType,
		},
	})
}

// ContactUs is the contact-us page handler
func (rp *Repository) ContactUs(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact-us.page.gtpl", &models.TemplateData{})
}
