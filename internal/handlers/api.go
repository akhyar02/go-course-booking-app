package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/akhyar02/bookings/internal/config"
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
		RoomId int `json:"roomId"`
	}{
		RoomId: 1,
	})
	fmt.Println(startDate, endDate)
	fmt.Fprint(w, string(jsonResponse))
}
