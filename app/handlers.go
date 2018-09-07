package main

import (
	"fmt"
	"net/http"
)

type personHandler struct{}

func (h personHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "person handler")
}
