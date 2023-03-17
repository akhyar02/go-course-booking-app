package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var handlerTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{
		"home", "/", "GET", http.StatusOK,
	},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, test := range handlerTests {
		response, err := ts.Client().Get(ts.URL + test.url)
		if err != nil {
			t.Fatal("Error making request:", err)
		}
		if response.StatusCode != test.expectedStatusCode {
			t.Errorf("for %s, expected %d but got %d", test.name, test.expectedStatusCode, response.StatusCode)
		}
	}
}
