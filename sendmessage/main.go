package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var id int
var endpoint *string
var key *string
var fname *string
var lname *string

func main() {
	endpoint = flag.String("endpoint", "", "Event Grid Endpoint")
	key = flag.String("key", "", "Event Grid Key")
	fname = flag.String("fname", "John", "First Name")
	lname = flag.String("lname", "Doe", "Last Name")
	flag.Parse()

	if *endpoint == "" || !strings.HasPrefix(*endpoint, "https://") || *key == "" || len(*key) < 32 {
		flag.Usage()
		return
	}

	sendMessage()
}

func sendMessage() {
	// event grid expects the message to be an array, so wrap in [ ]
	var wrapper []envelope
	wrapper = append(wrapper, genMessage())

	// Marshal the json
	if json, err := json.Marshal(wrapper); err != nil {
		log.Println(err)
		return
	}

	// create the request
	req, err := http.NewRequest("POST", *endpoint, bytes.NewBuffer(json))

	// add the SAS key
	req.Header.Add("aeg-sas-key", *key)

	// execute and get a response
	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
		return
	}
	// close the request
	defer resp.Body.Close()

	// read the response
	if _, err = ioutil.ReadAll(resp.Body); err != nil {
		log.Println(err)
		return
	}

	fmt.Println("Sent:", id)
}

// generate a message to post
func genMessage() envelope {
	var msg envelope

	if id < 10000 {
		rand.Seed(time.Now().UnixNano())
	}

	id = rand.Intn(89990) + 10005

	// event grid stuff
	msg.ID = strconv.Itoa(int(id))
	// msg.Topic = "" // don't set this
	msg.Subject = "person"
	msg.EventType = msg.Subject
	msg.EventTime = time.Now().UTC().Format("2006-01-02T15:04:05Z")

	// private data payload
	msg.Data.FirstName = *fname
	msg.Data.LastName = *lname

	return msg
}

// event grid message structure
// Data varies by message so is not processed
type envelope struct {
	ID        string `json:"id"`
	Topic     string `json:"topic"`
	Subject   string `json:"subject"`
	EventType string `json:"eventType"`
	EventTime string `json:"eventTime"`
	Data      person `json:"data"`
}

// private message structure
type person struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
