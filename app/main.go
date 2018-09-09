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
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// application startup
func main() {

	port := flag.Int("port", 8080, "TCP port to listen on")
	logPath := flag.String("logpath", "/home/LogFiles/", "path to write log files")

	flag.Parse()

	setupLog(logPath)

	runServer(*port)
}

func runServer(port int) {

	r := mux.NewRouter()
	// r.Methods("POST")
	// r.Headers("x-forwarded-proto", "https")

	r.Handle("/person", logb.Handler(eventgrid.Handler(personHandler)))
	http.Handle("/", r)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + strconv.Itoa(port),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	log.Println("Listening on port: ", port)

	log.Println(srv.ListenAndServe())

}

func setupLog(logPath *string) {
	// TODO - mkdir -p logFilePath
	// TODO - add instance id to file name

	if !strings.HasSuffix(*logPath, "/") {
		*logPath += "/"
	}

	logFile, err := os.OpenFile(*logPath+"app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	} else {
		wrt := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(wrt)
		log.SetFlags(log.Ldate | log.Ltime)
		log.Println("init complete")
	}
}
