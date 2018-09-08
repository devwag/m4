package main

import (
	"encoding/json"
	"fmt"
	"log"
	"m4/eventgrid"
	"net/http"
)

// this is the structure for the data portion of dc-receive messages
type person struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func personHandler(w http.ResponseWriter, r *http.Request, env *eventgrid.Envelope) {

	var p person
	err := json.Unmarshal(env.Data, &p)

	if err == nil {
		w.WriteHeader(200)
		fmt.Fprintln(w, "person Handler: ", p.FirstName, "", p.LastName)
		log.Println("person Handler: ", env.ID, p.FirstName, "", p.LastName)
	} else {
		fmt.Println(err)
	}
}
