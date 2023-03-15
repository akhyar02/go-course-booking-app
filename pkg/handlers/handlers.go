package handlers

import (
	"net/http"

	"github.com/akhyar02/bookings/pkg/config"
	"github.com/akhyar02/bookings/pkg/models"
	"github.com/akhyar02/bookings/pkg/render"
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
	render.RenderTemplate(w, "home.page.gtpl", &models.TemplateData{})
}

// About is the about page handler
func (rp *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again"
	render.RenderTemplate(w, "about.page.gtpl", &models.TemplateData{StringMap: stringMap})
}

// GeneralQuarters is the general-quarters page handler
func (rp *Repository) GeneralQuarters(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "generals.page.gtpl", &models.TemplateData{})
}

// MajorSuites is the major-suites page handler
func (rp *Repository) MajorSuites(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "majors.page.gtpl", &models.TemplateData{})
}

// MakeReservation is the make-reservation page handler
func (rp *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "make-reservation.page.gtpl", &models.TemplateData{})
}

// SearchAvailability is the search-availability page handler
func (rp *Repository) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	var urlQuery = r.URL.Query()
	var roomType = urlQuery.Get("room")
	render.RenderTemplate(w, "search-availability.page.gtpl", &models.TemplateData{
		Data: map[string]interface{}{
			"roomType": roomType,
		},
	})
}

// ContactUs is the contact-us page handler
func (rp *Repository) ContactUs(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact-us.page.gtpl", &models.TemplateData{})
}
