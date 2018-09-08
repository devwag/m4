package eventgrid

import (
	"encoding/json"
	"log"
	"net/http"
)

// Handler - handle the event grid message
func Handler(next func(w http.ResponseWriter, r *http.Request, env *Envelope)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// only support post
		// TODO - should we always check this?
		if r.Method != "POST" {
			w.WriteHeader(500)
			log.Println(r.Method, " Not supported")
		}

		// TODO - should we add an https check?

		var env Envelope

		// decode the event grid message from the body
		if r.Body != nil {
			err := json.NewDecoder(r.Body).Decode(&env)
			if err != nil {
				w.WriteHeader(500)
				log.Println(err)
				return
			}
			r.Body.Close()
		}

		// verify event grid ID
		if env.ID == "" {
			w.WriteHeader(500)
			log.Println("Event Grid Envelope: missing ID")
			return
		}

		// verify event grid has data
		// TODO - should we do this? are empty data messages possible?
		if env.Data == nil {
			w.WriteHeader(500)
			log.Println("Event Grid Envelope: missing Data")
			return
		}

		// call the next handler
		if next != nil {
			next(w, r, &env)
		}
	})
}
