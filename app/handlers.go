package main

import (
	"fmt"
	"net/http"
)

type fooHandler struct{}

func (h fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "foo")
}

type barHandler struct{}

func (h barHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "bar")
}
