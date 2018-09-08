package main

import (
	"encoding/json"
	"log"
	"m4/eventgrid"
	"net/http"
)

// this is the structure for the data portion of event grid messages
type person struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

// handle event grid "person" messages
// w and r are standard http.Handler params
// env is the event grid envelope that was parsed by the eveng grid handler
func personHandler(w http.ResponseWriter, r *http.Request, env *eventgrid.Envelope) {

	log.Println("x-forwarded-proto", r.Header.Get("x-forwarded-proto"))
	log.Println("x-arr-ssl", r.Header.Get("x-arr-ssl"))

	// get the values from env.Data
	var p person
	err := json.Unmarshal(env.Data, &p)

	if err == nil {
		w.WriteHeader(200)

		// event grid doesn't inspect the body on a 200

		log.Println("person Handler: ", env.ID, p.FirstName, p.LastName)

		// TODO this is where you would process the "person message"
	} else {
		w.WriteHeader(500)
		log.Println("ERROR:", err)
	}
}
