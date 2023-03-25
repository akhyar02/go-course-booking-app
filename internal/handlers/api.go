package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
	"time"

	"github.com/akhyar02/bookings/internal/config"
	"github.com/akhyar02/bookings/internal/driver"
	"github.com/akhyar02/bookings/internal/helpers"
	"github.com/akhyar02/bookings/internal/models"
	"github.com/akhyar02/bookings/internal/repository"
	"github.com/akhyar02/bookings/internal/repository/dbrepo"
)

type reservationApiHandler struct {
	appConfig config.AppConfig
	db        repository.DatabaseRepo
}

var ReservationApi *reservationApiHandler

func NewReservationHandler(appConfig config.AppConfig, db *driver.DB) {
	ReservationApi = &reservationApiHandler{
		appConfig: appConfig,
		db:        dbrepo.NewPostgresRepo(db.SQL, &appConfig),
	}
}

func (h *reservationApiHandler) GetReservationByDate(w http.ResponseWriter, r *http.Request) {
	var UrlQueries = r.URL.Query()
	var (
		startDateParam = UrlQueries.Get("startDate")
		endDateParam   = UrlQueries.Get("endDate")
		roomIdParam    = UrlQueries.Get("roomId")
	)

	roomId, err := strconv.Atoi(roomIdParam)
	if err != nil {
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}
	startDate, err := time.Parse("2006-01-02", startDateParam)
	if err != nil {
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}
	endDate, err := time.Parse("2006-01-02", endDateParam)
	if err != nil {
		helpers.ClientError(w, http.StatusBadRequest)
		return
	}
	isAvailable, err := h.db.SearchAvailibilityByDatesAndRoomId(startDate, endDate, roomId)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if !isAvailable {
		helpers.ClientError(w, http.StatusConflict)
		return
	}

	jsonResponse, _ := json.Marshal(struct {
		RoomId    int    `json:"roomId"`
		RoomType  string `json:"roomType"`
		StartDate string `json:"startDate"`
		EndDate   string `json:"endDate"`
	}{
		RoomId:    1,
		RoomType:  "general_quarters",
		StartDate: startDateParam,
		EndDate:   endDateParam,
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
		RoomId    string `json:"room_id"`
		RoomType  string `json:"room_type"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		helpers.ServerError(w, err)
		return
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
	var roomId int
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
	if requestBody.RoomId == "" {
		errors["room_type"] = append(errors["room_type"], "room id is required")
	} else {
		roomId, err = strconv.Atoi(requestBody.RoomId)
		if err != nil {
			errors["room_type"] = append(errors["room_type"], "room id must be a number")
		}
	}

	reservation := models.Reservation{
		FirstName: requestBody.FirstName,
		LastName:  requestBody.LastName,
		Email:     requestBody.Email,
		Phone:     requestBody.Phone,
		RoomId:    roomId,
		Room: models.Room{
			RoomName: requestBody.RoomType,
		},
		StartDate: startDate,
		EndDate:   endDate,
	}

	if len(errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		isRoomAvailable, err := h.db.SearchAvailibilityByDatesAndRoomId(startDate, endDate, roomId)
		if err != nil {
			log.Println(err)
			helpers.ServerError(w, err)
		}
		if !isRoomAvailable {
			w.WriteHeader(http.StatusConflict)
			return
		}

		w.WriteHeader(http.StatusCreated)
		reservationId, err := h.db.InsertReservation(reservation)
		if err != nil {
			log.Println(err)
			helpers.ServerError(w, err)
			return
		}

		roomRestriction := models.RoomRestriction{
			StartDate:     startDate,
			EndDate:       endDate,
			RoomId:        roomId,
			ReservationId: reservationId,
			RestrictionId: 1,
		}

		err = h.db.InsertRoomRestriction(roomRestriction)
		if err != nil {
			log.Println(err)
			helpers.ServerError(w, err)
			return
		}

		h.appConfig.MailChan <- models.MailData{
			To:      []string{reservation.Email},
			From:    "xyzbooking@xyz.com",
			Subject: "Reservation Confirmation",
			Content: `<strong>Reservation Confirmation</strong><br/>
			Thank you for your reservation. We look forward to seeing you.`,
		}
		h.appConfig.MailChan <- models.MailData{
			To:      []string{"admin@xyz.com"},
			From:    "xyzbooking@xyz.com",
			Subject: "Room Reservation",
			Content: `<strong>Room is reserved</strong><br/>`,
		}
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
