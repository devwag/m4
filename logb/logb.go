package logb

import (
	"log"
	"net/http"
	"time"
)

// TODO - should we replace with Gorilla logger?

// TODO - mkdir -p logFilePath
// TODO - add MultiWriter support
// TODO - use flag for logfilepath
// TODO - add error log and method
// TODO - add verbose support
// TODO - add Apache log file support

var logFilePath = "/home/LogFiles/"

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

		log.Println(wr.status, r.URL.Path, time.Now().Sub(wr.start))
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
