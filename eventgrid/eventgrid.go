package eventgrid

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// TODO - Log / Reject bad EG messages
// TODO - handle upcoming EG change that can send multiple items in a message

// Handler - handle the event grid message
func Handler(next func(w http.ResponseWriter, r *http.Request, env *Envelope)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var env Envelope
		var err error
		var msg []Envelope

		// validate the request
		if r.Body == nil {
			err = fmt.Errorf("No request body")
		}

		// decode the event grid message from the body
		if err == nil {
			err = json.NewDecoder(r.Body).Decode(&msg)
		}

		// validate the event grid envelope
		if err == nil {
			r.Body.Close()
			env = msg[0]
			err = ValidateEnvelope(&env)
		}

		if err == nil {
			// handle event grid subscription validation events
			if env.EventType == "Microsoft.EventGrid.SubscriptionValidationEvent" {
				r.URL.RawQuery = "validate"
				err = handleValidate(w, &env)
			} else {
				// call the next handler
				if next != nil {
					next(w, r, &env)
				}
			}
		}

		// log any error and return 500
		if err != nil {
			log.Println(err)
			w.WriteHeader(500)
		}
	})
}

// ValidateEnvelope - validates a message grid envelope contains required fields
func ValidateEnvelope(env *Envelope) error {
	// verify event grid ID
	if env.ID == "" {
		return fmt.Errorf("Event Grid Envelope: missing ID")
	}

	// verify event grid has data
	// TODO - should we do this? are empty data messages possible?
	if env.Data == nil {
		return fmt.Errorf("Event Grid Envelope: missing Data")
	}

	// TODO - add more validations
	return nil
}

func handleValidate(w http.ResponseWriter, msg *Envelope) error {
	// get the validationCode from the json (that's all we care about)
	var vData struct {
		ValidationCode string `json:"validationCode"`
		ValidationURL  string `json:"validationUrl"`
	}
	err := json.Unmarshal(msg.Data, &vData)

	// handle the json error
	if err != nil {
		return err
	}

	// return the validationCode as json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	fmt.Fprintf(w, "{ \"validationResponse\": \"%v\" }", vData.ValidationCode)
	log.Println("EventGridValidation: Success")

	return nil
}
