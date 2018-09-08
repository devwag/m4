package main

import (
	"bytes"
	"m4/eventgrid"
	"m4/logb"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogb(t *testing.T) {

	s := `{"id":"1001","topic":"","subject":"person","eventType":"person","eventTime":"2018-08-20T18:04:26Z", "dataVersion": "1.0","data":{"firstName": "John", "LastName": "Doe"}}`

	// var env eventgrid.Envelope
	// err := json.Unmarshal([]byte(s), &env)

	// var p person
	// json.Unmarshal(env.Data, &p)

	r, err := http.NewRequest("POST", "https://www.logb.com/person", bytes.NewBufferString(s))
	if err != nil {
		t.Error("NewRequest: ", err)
	}

	w := httptest.NewRecorder()

	h := http.Handler(logb.Handler(eventgrid.Handler(personHandler)))

	h.ServeHTTP(w, r)

	if err != nil {
		t.Error("Request Error: ", err.Error())
	}

	if w.Code != 200 {
		t.Error("Error Code: ", w.Code)
	}

}
