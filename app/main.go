package main

import (
	"flag"
	"log"
	"m4/eventgrid"
	"m4/logb"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// application startup
func main() {

	port := flag.Int("port", 8080, "TCP port to listen on")
	flag.Parse()

	r := mux.NewRouter()
	// r.Methods("POST")
	// r.Headers("x-forwarded-proto", "https")

	r.Handle("/person", logb.Handler(eventgrid.Handler(personHandler)))
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + strconv.Itoa(*port),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Println("Listening on port: ", *port)

	log.Println(srv.ListenAndServe())

}
