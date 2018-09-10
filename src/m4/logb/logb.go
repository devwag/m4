package logb

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// TODO - should we replace with Gorilla logger?
// TODO - add Apache log file support

var reqLog = log.New(os.Stdout, "", log.Ldate|log.Ltime)

// Init - initialize the log writer
func Init(logPath string) {
	log.Println("request log", logPath)

	logFile, err := openLogFile(logPath)

	if err != nil {
		log.Println("init", err)
		log.Println(err)
		return
	}

	wrt := io.MultiWriter(os.Stdout, logFile)
	reqLog.SetOutput(wrt)

	reqLog.Println("testing")
}

//Handler - http handler that writes to log file(s)
func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		wr := &ResponseLogger{
			ResponseWriter: w,
			status:         0,
			start:          time.Now().UTC()}

		if next != nil {
			next.ServeHTTP(wr, r)
		}

		reqLog.Println(wr.status, time.Now().UTC().Sub(wr.start).Nanoseconds()/100000, r.Method, r.URL.Path, r.URL.RawQuery)
	})
}

// ResponseLogger - wrap http.ResponseWriter to include status and size
type ResponseLogger struct {
	http.ResponseWriter
	status int
	size   int
	start  time.Time
}

// WriteHeader - wraps http.WriteHeader
func (r *ResponseLogger) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

// Write - wraps http.Write
func (r *ResponseLogger) Write(buf []byte) (int, error) {
	n, err := r.ResponseWriter.Write(buf)

	if err == nil {
		r.size += n
	}

	return n, err
}

// Open the log file
func openLogFile(logPath string) (*os.File, error) {
	fileName := logPath

	if !strings.HasSuffix(fileName, "/") {
		fileName += "/"
	}

	// make the log directory if it doesn't exist
	if err := os.MkdirAll(fileName, 0666); err != nil {
		return nil, err
	}

	fileName += "request"

	// use instance ID to differentiate log files between instances in App Services
	if iid := os.Getenv("WEBSITE_ROLE_INSTANCE_ID"); iid != "" {
		fileName += "_" + strings.TrimSpace(iid)
	}

	fileName += ".log"

	// open the log file
	return os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
}
