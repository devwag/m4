package eventgrid

import (
	"fmt"
	"net/http"
)

// Handler - handle the event grid message
func Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Print("eventGridHandler ")

		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}
