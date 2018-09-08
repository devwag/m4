package main

import (
	"bytes"
	"encoding/json"
	"m4/eventgrid"
	"m4/logb"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLogb(t *testing.T) {

	e, err := genMessage(t)

	if err != nil {
		t.Error(err)
	}

	r, err := http.NewRequest("POST", "https://www.logb.com/person", bytes.NewBuffer(e)) // bytes.NewBufferString(s))
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

// helper function to generate valid event grid message
func genMessage(t *testing.T) ([]byte, error) {
	env := eventgrid.Envelope{Subject: "person", EventType: "person", DataVersion: "1.0"}
	env.ID = "1001"
	env.EventTime = time.Now().UTC().Format("2006-01-02T15:04:05Z")

	p := person{FirstName: "John", LastName: "Doe"}

	var err error
	env.Data, err = json.Marshal(&p)

	if err != nil {
		t.Error(err)
	}

	var wrapper []eventgrid.Envelope
	wrapper = append(wrapper, env)

	return json.Marshal(wrapper)
}
