package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/akhyar02/bookings/internal/config"
	"github.com/akhyar02/bookings/internal/models"
)

type reservationApiHandler struct {
	appConfig config.AppConfig
}

var ReservationApi *reservationApiHandler

func NewReservationHandler(appConfig config.AppConfig) {
	ReservationApi = &reservationApiHandler{
		appConfig: appConfig,
	}
}

func (h *reservationApiHandler) GetReservationByDate(w http.ResponseWriter, r *http.Request) {
	var UrlQueries = r.URL.Query()
	var (
		startDate = UrlQueries.Get("startDate")
		endDate   = UrlQueries.Get("endDate")
	)
	jsonResponse, _ := json.Marshal(struct {
		RoomId    int    `json:"roomId"`
		RoomType  string `json:"roomType"`
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
	}{
		RoomId:    1,
		RoomType:  "general_quarters",
		StartDate: startDate,
		EndDate:   endDate,
	})
	fmt.Fprint(w, string(jsonResponse))
}

func (h *reservationApiHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var errors = make(map[string][]string)
	var requestBody struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		RoomType  string `json:"room_type"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		fmt.Println(err)
	}

	if requestBody.RoomType == "" {
		errors["room_type"] = append(errors["room_type"], "room type is required")
	} else {
		switch strings.ToLower(requestBody.RoomType) {
		case "general quarters":
		case "major suites":
		default:
			log.Println(requestBody.RoomType)
			errors["room_type"] = append(errors["room_type"], "room type is not valid")
		}
	}

	var startDate time.Time
	var endDate time.Time
	if requestBody.StartDate == "" {
		errors["start_date"] = append(errors["start_date"], "start date is required")
	} else {
		parsedDate, err := time.Parse("2006-01-02", requestBody.StartDate)
		if err != nil {
			errors["start_date"] = append(errors["start_date"], "start date is not valid")
		}
		startDate = parsedDate
	}
	if requestBody.EndDate == "" {
		errors["end_date"] = append(errors["end_date"], "end date is required")
	} else {
		parsedDate, err := time.Parse("2006-01-02", requestBody.EndDate)
		if err != nil {
			errors["end_date"] = append(errors["end_date"], "end date is not valid")
		}
		endDate = parsedDate
	}
	if startDate.After(endDate) {
		errors["end_date"] = append(errors["end_date"], "end date must be after start date")
	}
	if requestBody.FirstName == "" {
		errors["first_name"] = append(errors["first_name"], "first name is required")
	}
	if requestBody.LastName == "" {
		errors["last_name"] = append(errors["last_name"], "last name is required")
	}
	if requestBody.Email == "" {
		errors["email"] = append(errors["email"], "email is required")
	} else {
		_, err := mail.ParseAddress(requestBody.Email)
		if err != nil {
			errors["email"] = append(errors["email"], "email is not valid")
		}
	}
	if requestBody.Phone == "" {
		errors["phone"] = append(errors["phone"], "phone is required")
	}

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	reservation := models.Reservation{
		FirstName: requestBody.FirstName,
		LastName:  requestBody.LastName,
		Email:     requestBody.Email,
		Phone:     requestBody.Phone,
		RoomID:    1,
		RoomType:  requestBody.RoomType,
		StartDate: requestBody.StartDate,
		EndDate:   requestBody.EndDate,
	}

	var responseData = struct {
		Errors map[string][]string `json:"errors"`
		Data   models.Reservation  `json:"data"`
	}{
		Errors: errors,
		Data:   reservation,
	}

	h.appConfig.Session.Put(r.Context(), "reservation", reservation)
	jsonResponse, _ := json.Marshal(responseData)
	fmt.Fprint(w, string(jsonResponse))
}
