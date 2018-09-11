package main

import (
	"bytes"
	"encoding/json"
	"m4/eventgrid"
	"m4/logb"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestMain(t *testing.T) {

	if err := os.MkdirAll("./logs/", 0666); err != nil {
		t.Error(err)
	}
}

// TODO - this isn't working :(
func testFlags(t *testing.T) {

	lp, err := validateFlags()

	if err != nil {
		t.Error(err)
	}

	if lp != "./logs/" {
		t.Error("invalid logPath:", lp)
	}
}

// test the main() app

// TODO - this is failing in CI/CD - need to debug

// func TestMainFunc(t *testing.T) {
// 	go main()
// 	time.Sleep(500 * time.Millisecond)

// 	osChan <- os.Interrupt

// 	time.Sleep(500 * time.Millisecond)
// }

// test person handler
func TestPersonHandler(t *testing.T) {

	e, err := genMessage(t)

	if err != nil {
		t.Error(err)
	}

	r, err := http.NewRequest("POST", "https://www.logb.com/person", bytes.NewBuffer(e)) // bytes.NewBufferString(s))
	if err != nil {
		t.Error("NewRequest: ", err)
	}

	w := httptest.NewRecorder()

	// chain the person handler
	h := http.Handler(logb.Handler(eventgrid.Handler(personHandler)))

	h.ServeHTTP(w, r)

	if w.Code != 200 {
		t.Error("Return Code: ", w.Code)
	}

}

// helper function to generate valid event grid json
func genMessage(t *testing.T) ([]byte, error) {
	env := eventgrid.Envelope{Subject: "person", EventType: "person", DataVersion: "1.0"}
	env.ID = "1001"
	env.EventTime = time.Now().UTC().Format("2006-01-02T15:04:05Z")

	p := person{FirstName: "John", LastName: "Doe"}

	var err error
	env.Data, err = json.Marshal(&p)

	if err != nil {
		return nil, err
	}

	var wrapper []eventgrid.Envelope
	wrapper = append(wrapper, env)

	return json.Marshal(wrapper)
}
