package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Event struct {
	ID int `json:"Id"`
	Message string `json:"Message"`
	TimeCreated string `json:"TimeCreated"`
	Severity string `json:"Severity"`
}

func handleLog(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Method == "POST" {
		fmt.Println("POST METHOD RECVD")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "could not read body", http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()
		var event []Event
		if err := json.Unmarshal(body, &event); err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
			return
		}

		// save event to database to be able to use on frontend maybe?

		// Check event ID to see if it's new or not - don't save to DB if it's old

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Event received successfully"))
	}
}

func main() {
	fmt.Println("UP N RUNNIN")

	http.HandleFunc("/logEvent", handleLog)

	err := http.ListenAndServe(":3333", nil)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}