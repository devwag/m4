package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"m4/eventgrid"
	"m4/logb"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// application startup
func main() {

	port := flag.Int("port", 8080, "TCP port to listen on")
	flag.Parse()

	// TODO - this is temporary
	logFile, err := os.OpenFile("/home/LogFiles/app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	} else {
		wrt := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(wrt)
		log.SetFlags(log.Ldate | log.Ltime)
		log.Println("init complete")
	}

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
