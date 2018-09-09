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

// TODO - mkdir -p logFilePath
// TODO - use flag for logfilepath
// TODO - add Apache log file support
// TODO - add instance id to file name

var reqLog = log.New(os.Stdout, "", log.Ldate|log.Ltime)

func init() {
	var logFilePath = "/home/LogFiles/"

	if !strings.HasSuffix(logFilePath, "/") {
		logFilePath += "/"
	}

	logFile, err := os.OpenFile(logFilePath+"request.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		log.Println(err)
	} else {
		wrt := io.MultiWriter(os.Stdout, logFile)
		reqLog.SetOutput(wrt)
	}
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

		reqLog.Println(wr.status, time.Now().UTC().Sub(wr.start).Nanoseconds()/100000, r.URL.Path, r.Method, r.URL.RawQuery)
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
