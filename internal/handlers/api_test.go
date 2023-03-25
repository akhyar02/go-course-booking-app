package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type reservationRequestBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	RoomType  string `json:"room_type"`
	RoomId    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

var validReservationParams = reservationRequestBody{
	"john",
	"doe",
	"johnD@mail.com",
	"123456789",
	"general quarters",
	"1",
	"2025-03-03",
	"2025-03-04",
}

var apiTests = []struct {
	name               string
	url                string
	method             string
	params             interface{}
	expectedStatusCode int
}{
	{
		"Get reservation by date", "/api/reservations?startDate=2023-03-03&endDate=2023-03-04&roomId=1", "GET", nil, http.StatusOK,
	},
	{
		"Create reservation", "/api/reservations", "POST", validReservationParams, http.StatusCreated,
	},
	{
		"Create reservation with invalid date",
		"/api/reservations",
		"POST",
		reservationRequestBody{
			FirstName: validReservationParams.FirstName,
			LastName:  validReservationParams.LastName,
			Email:     validReservationParams.Email,
			Phone:     validReservationParams.Phone,
			RoomType:  validReservationParams.RoomType,
			RoomId:    validReservationParams.RoomId,
			StartDate: "2022-3-3",
			EndDate:   "2022-3-2",
		},
		http.StatusBadRequest,
	},
	{
		"Create reservation with invalid email",
		"/api/reservations",
		"POST",
		reservationRequestBody{
			FirstName: validReservationParams.FirstName,
			LastName:  validReservationParams.LastName,
			Email:     "",
			Phone:     validReservationParams.Phone,
			RoomType:  validReservationParams.RoomType,
			StartDate: validReservationParams.StartDate,
			RoomId:    validReservationParams.RoomId,
			EndDate:   validReservationParams.EndDate,
		},
		http.StatusBadRequest,
	},
	{
		"Create reservation with invalid phone",
		"/api/reservations",
		"POST",
		reservationRequestBody{
			FirstName: validReservationParams.FirstName,
			LastName:  validReservationParams.LastName,
			Email:     validReservationParams.Email,
			Phone:     "",
			RoomId:    validReservationParams.RoomId,
			RoomType:  validReservationParams.RoomType,
			StartDate: validReservationParams.StartDate,
			EndDate:   validReservationParams.EndDate,
		},
		http.StatusBadRequest,
	},
}

func TestGetReservationByDate(t *testing.T) {
	routes := getRoutes()
	var request *http.Request
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, test := range apiTests {
		var response *http.Response
		var err error
		switch test.method {
		case "GET":
			response, err = ts.Client().Get(ts.URL + test.url)
		case "POST":
			var postBody []byte
			postBody, err = json.Marshal(test.params)
			if err != nil {
				t.Fatal("Error marshalling json:", err)
			}
			request, err = http.NewRequest("POST", ts.URL+test.url, bytes.NewBuffer(postBody))
			if err != nil {
				t.Fatal("Error creating request:", err)
			}
			response, err = ts.Client().Do(request)
		}
		if err != nil {
			t.Fatal("Error making request:", err)
		}
		if response.StatusCode != test.expectedStatusCode {
			body, _ := io.ReadAll(response.Body)
			t.Errorf("for %s, expected %d but got %d, %s", test.name, test.expectedStatusCode, response.StatusCode, string(body))
		}
	}
}
