package logb

import (
	"log"
	"net/http"
	"time"
)

func LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wr := NewResponseRecorder(w)

		if next != nil {
			next.ServeHTTP(wr, r)
		}

		log.Println(wr.status, r.URL.Path, time.Now().Sub(wr.start))
	})
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	size   int
	start  time.Time
}

func NewResponseRecorder(w http.ResponseWriter) *responseRecorder {
	return &responseRecorder{
		ResponseWriter: w,
		status:         http.StatusOK,
		start:          time.Now(),
	}
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(buf []byte) (int, error) {
	n, err := r.ResponseWriter.Write(buf)
	if err == nil {
		r.size += n
	}
	return n, err
}
