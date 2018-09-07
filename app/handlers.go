package main

import (
	"fmt"
	"net/http"
)

// this is the structure for the data portion of dc-receive messages
type person struct {
	DC    string `json:"DC"`
	Items []struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	} `json:"Items"`
}

type personHandler struct{}

func (h personHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "person handler")
}
