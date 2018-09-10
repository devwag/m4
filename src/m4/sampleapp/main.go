package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
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

// port to listen on (can be changed with -port option)
var port = 8080

func setupLogs(logPath string) error {
	if err := setupLog(logPath); err != nil {
		return err
	}

	// setup logb multiwriter
	if err := logb.SetLogPath(logPath); err != nil {
		return err
	}

	return nil
}

func validateFlags() (string, error) {
	// read the flags
	logPath := flag.String("logpath", "./logs/", "path to write log files")
	p := flag.Int("port", port, "TCP port to listen on")
	flag.Parse()

	fmt.Println("logpath", *logPath)
	fmt.Println("port", *p)

	// validate the logpath flag
	// TODO - add more checks
	lp := strings.TrimSpace(*logPath)
	if lp == "" {
		return "", errors.New("invalid -logpath flag")
	}

	// validate the port flag
	if *p >= 0 && *p <= 65535 {
		// set the listen port
		port = *p
	} else {
		return lp, errors.New("invalid -port flag")
	}

	return lp, nil
}

// application startup
func main() {
	logPath, err := validateFlags()

	if err != nil {
		flag.PrintDefaults()
		log.Fatal(err)
	}

	// setup the log multi writer
	if err = setupLogs(logPath); err != nil {
		log.Fatal(err)
	}

	log.Println("Listening on port: ", port)

	// run the server
	if err := runServer(port); err != nil {
		log.Println("ERROR:", err)
	}

	log.Println("Server Exit")
}

func runServer(port int) error {

	// use gorilla mux
	r := mux.NewRouter()

	// this is our only handler
	// chain the handlers together as middleware
	// app services does https offloading, so check for the x-forwarded-proto header
	// only accept POST requests
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
func setupLog(logPath string) error {

	// prepend date and time to log entries
	log.SetFlags(log.Ldate | log.Ltime)

	// open the log file
	logFile, err := openLogFile(logPath, "app", ".log")

	if err != nil {
		return err
	}

	// setup a multiwriter to log to file and stdout
	wrt := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(wrt)

	return nil
}

// Open the log file
func openLogFile(logPath string, logPrefix string, logExtension string) (*os.File, error) {

	fileName, err := getFullLogName(logPath, logPrefix, logExtension)

	if err != nil {
		return nil, err
	}

	// open the log file
	return os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}

// get the full log file name and make the directory if necessary
func getFullLogName(logPath string, logPrefix string, logExtension string) (string, error) {
	if !strings.HasSuffix(logPath, "/") {
		logPath += "/"
	}

	// make the log directory if it doesn't exist
	if err := os.MkdirAll(logPath, 0666); err != nil {
		return "", err
	}

	fileName := logPath + logPrefix

	// use instance ID to differentiate log files between instances in App Services
	if iid := os.Getenv("WEBSITE_ROLE_INSTANCE_ID"); iid != "" {
		fileName += "_" + strings.TrimSpace(iid)
	}

	// make sure logExtension has a period
	if !strings.HasPrefix(logExtension, ".") {
		logExtension = "." + logExtension
	}

	return fileName + logExtension, nil
}
