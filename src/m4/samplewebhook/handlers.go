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
// env is the event grid envelope that was parsed by the event grid handler
func personHandler(w http.ResponseWriter, r *http.Request, env *eventgrid.Envelope) {

	// get the values from env.Data
	var p person

	if err := json.Unmarshal(env.Data, &p); err != nil {
		w.WriteHeader(500)
		log.Println("ERROR:", err)
		return
	}

	w.WriteHeader(200)

	// event grid doesn't inspect the body on a 200

	log.Println("person Handler: ", env.ID, p.FirstName, p.LastName)

	// TODO this is where you would process the "person message"
}
