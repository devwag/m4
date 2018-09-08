package eventgrid

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGet(t *testing.T) {

	r, err := http.NewRequest("GET", "https://www.eventgrid.com/", nil)

	if err != nil {
		t.Error("NewRequest: ", err)
	}

	w := httptest.NewRecorder()

	h := Handler(nil)
	h.ServeHTTP(w, r)

	if err != nil {
		t.Error("Request Error: ", err.Error())
	}

	if w.Code != 500 {
		t.Error("Error Code: ", w.Code)
	}

}

func TestEmptyBody(t *testing.T) {

	r, err := http.NewRequest("POST", "https://www.eventgrid.com/", nil)

	if err != nil {
		t.Error("NewRequest: ", err)
	}

	w := httptest.NewRecorder()

	h := Handler(nil)
	h.ServeHTTP(w, r)

	if err != nil {
		t.Error("Request Error: ", err.Error())
	}

	if w.Code != 500 {
		t.Error("Error Code: ", w.Code)
	}

}

func TestGoodBody(t *testing.T) {

	s := `{"id":"1001","topic":"","subject":"person","eventType":"person","eventTime":"2018-08-20T18:04:26Z","dataVersion":"1.0","data":{"firstName": "John", "LastName": "Doe"}}`

	r, err := http.NewRequest("POST", "https://www.eventgrid.com/", bytes.NewBufferString(s))

	if err != nil {
		t.Error("NewRequest: ", err)
	}

	w := httptest.NewRecorder()

	h := Handler(nil)
	h.ServeHTTP(w, r)

	if err != nil {
		t.Error("Request Error: ", err.Error())
	}

	if w.Code != 200 {
		t.Error("Error Code: ", w.Code)
	}

}

func TestBadJson(t *testing.T) {

	s := `{"id":"1001","topic":"","subject":"person","eventType":"person","eventTime":"2018-08-20T18:04:26Z","data":{"firstName": "John", "LastName": "Doe"},}`

	r, err := http.NewRequest("POST", "https://www.eventgrid.com/", bytes.NewBufferString(s))

	if err != nil {
		t.Error("NewRequest: ", err)
	}

	w := httptest.NewRecorder()

	h := Handler(nil)
	h.ServeHTTP(w, r)

	if err != nil {
		t.Error("Request Error: ", err.Error())
	}

	if w.Code != 500 {
		t.Error("Error Code: ", w.Code)
	}

}

func TestMissingId(t *testing.T) {

	s := `{"topic":"","subject":"person","eventType":"person","eventTime":"2018-08-20T18:04:26Z","dataVersion":"1.0","data":{"firstName": "John", "LastName": "Doe"}}`

	r, err := http.NewRequest("POST", "https://www.eventgrid.com/", bytes.NewBufferString(s))

	if err != nil {
		t.Error("NewRequest: ", err)
	}

	w := httptest.NewRecorder()

	h := Handler(nil)
	h.ServeHTTP(w, r)

	if err != nil {
		t.Error("Request Error: ", err.Error())
	}

	if w.Code != 500 {
		t.Error("Error Code: ", w.Code)
	}

}

func TestMissingData(t *testing.T) {

	s := `{"id":"1001","topic":"","subject":"person","eventType":"person","eventTime":"2018-08-20T18:04:26Z","dataVersion":"1.0"}`

	r, err := http.NewRequest("POST", "https://www.eventgrid.com/", bytes.NewBufferString(s))

	if err != nil {
		t.Error("NewRequest: ", err)
	}

	w := httptest.NewRecorder()

	h := Handler(nil)
	h.ServeHTTP(w, r)

	if err != nil {
		t.Error("Request Error: ", err.Error())
	}

	if w.Code != 500 {
		t.Error("Error Code: ", w.Code)
	}

}
