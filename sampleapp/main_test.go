package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"m4/eventgrid"
	"m4/logb"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLogb(t *testing.T) {

	//	s := `{"id":"1001","topic":"","subject":"person","eventType":"person","eventTime":"2018-08-20T18:04:26Z", "dataVersion": "1.0","data":{"firstName": "John", "LastName": "Doe"}}`

	// TODO use Envelope to build the json
	env := eventgrid.Envelope{Subject: "person", EventType: "person", DataVersion: "1.0"}
	env.ID = "1001"
	env.EventTime = time.Now().UTC().Format("2006-01-02T15:04:05Z")

	p := person{FirstName: "John", LastName: "Doe"}
	json.Unmarshal(env.Data, &p)

	var b []byte
	err := json.Unmarshal(b, &env)

	fmt.Println(err)

	fmt.Println(string(b))

	return

	r, err := http.NewRequest("POST", "https://www.logb.com/person", bytes.NewBufferString(string(b))) // bytes.NewBufferString(s))
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
