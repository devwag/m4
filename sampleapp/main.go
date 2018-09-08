package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os/signal"

	// TODO should change these to github.com/bartr/m4 once stable
	"m4/eventgrid"
	"m4/logb"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// channel used to send os.Interrupts
var osChan = make(chan os.Signal, 1)

// application startup
func main() {

	// parse command line flags
	port := flag.Int("port", 8080, "TCP port to listen on")
	logPath := flag.String("logpath", "/home/LogFiles/", "path to write log files")
	flag.Parse()

	setupLog(logPath)

	log.Println("Listening on port: ", *port)

	// run the server
	err := runServer(*port)

	if err != nil {
		log.Println("ERROR:", err)
	}

	log.Println("Server Exit")
}

func runServer(port int) error {

	// use gorilla mux
	r := mux.NewRouter()
	// TODO this causes app services to reject the request
	//     need to debug
	// putting the methods call on the handler works correctly
	// r.Methods("POST")

	// TODO make sure the request used https
	// r.Headers("x-forwarded-proto", "https")

	// this is our only handler
	// chain the handlers together as middleware
	r.Handle("/person", logb.Handler(eventgrid.Handler(personHandler))).Methods("POST").Headers("x-forwarded-proto", "https")
	http.Handle("/", r)

	// setup the server
	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + strconv.Itoa(port),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	// run webserver in a Go routine so we can cancel
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println("ERROR:", err)
		}
	}()

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	signal.Notify(osChan, os.Interrupt)

	// Block until we receive our signal
	<-osChan

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait until the timeout deadline
	return srv.Shutdown(ctx)
}

// setup log multi writer
func setupLog(logPath *string) {

	// prepend date and time to log entries
	log.SetFlags(log.Ldate | log.Ltime)

	// open the log file
	logFile, err := openLogFile(logPath)

	// we ignore the open error which means everything gets logged to stdout
	if err != nil {
		log.Println(err)
	} else {
		// setup a multiwriter to log to file and stdout
		wrt := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(wrt)
		log.Println("init complete")
	}
}

// Open the log file
func openLogFile(logPath *string) (*os.File, error) {
	fileName := *logPath

	if !strings.HasSuffix(fileName, "/") {
		fileName += "/"
	}

	// make the log directory if it doesn't exist
	err := os.MkdirAll(fileName, 0666)
	if err != nil {
		return nil, err
	}

	fileName += "app"

	// use instance ID to differentiate log files between instances in App Services
	iid := os.Getenv("WEBSITE_ROLE_INSTANCE_ID")
	if iid != "" {
		fileName += "_" + strings.TrimSpace(iid)
	}

	fileName += ".log"

	// open the log file
	return os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}
