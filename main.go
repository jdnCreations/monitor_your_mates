package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/jdnCreations/monitor_your_mates/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Event struct {
	ID          int    `json:"Id"`
	Message     string `json:"Message"`
	TimeCreated string `json:"TimeCreated"`
	Severity    string `json:"Severity"`
}

type apiConfig struct {
	db     *database.Queries
	secret string
}

func (cfg *apiConfig) handleLog(w http.ResponseWriter, r *http.Request) {
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
		var events []Event
		if err := json.Unmarshal(body, &events); err != nil {
			http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
			return
		}

		for _, event := range events {
			if event.Severity != "Critical" && event.Severity != "Error" && event.Severity != "Warning" {
				return
			}
			time, err := parseDate(event.TimeCreated)
			if err != nil {
				fmt.Println("cannot parse timecreated to time.Time")
				break
			}
			_, err = cfg.db.GetEventById(r.Context(), database.GetEventByIdParams{
				ID:        int32(event.ID),
				CreatedAt: time,
			})
			if err != nil {
				fmt.Println("event doesn't exist, adding to db.")
				cfg.db.CreateEvent(r.Context(), database.CreateEventParams{
					ID:        int32(event.ID),
					Message:   sql.NullString{String: event.Message, Valid: true},
					CreatedAt: time,
					Severity:  sql.NullString{String: event.Severity, Valid: true},
				})
			} else {
				fmt.Println("event already exists")
			}
		}

		// save event to database to be able to use on frontend maybe?

		// Check event ID to see if it's new or not - don't save to DB if it's old

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Event received successfully"))
	}

}

func parseDate(dateStr string) (time.Time, error) {
	// Use a regular expression to extract the timestamp from the string
	re := regexp.MustCompile(`\/Date\((\d+)\)\/`)
	matches := re.FindStringSubmatch(dateStr)

	if len(matches) < 2 {
		return time.Time{}, fmt.Errorf("invalid date format")
	}

	// Convert the extracted timestamp (milliseconds) to an integer
	millis, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	// Convert milliseconds to seconds and nanoseconds
	seconds := millis / 1000
	nanoseconds := (millis % 1000) * 1e6

	// Return the parsed time
	return time.Unix(seconds, nanoseconds), nil
}

func (cfg *apiConfig) handleEvents(w http.ResponseWriter, r *http.Request) {
	events, err := cfg.db.GetEvents(r.Context(), 10)
	if err != nil {
		http.Error(w, "Unable to fetch events", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		http.Error(w, "Error encoding events to JSON", http.StatusInternalServerError)
		return
	}
}

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	secret := os.Getenv("SECRET")
	port := os.Getenv("PORT")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Could not open connection to db", err)
	}
	println("Connected to db")

	dbQueries := database.New(db)
	apiConf := apiConfig{}
	apiConf.db = dbQueries
	apiConf.secret = secret

	http.HandleFunc("/logEvent", apiConf.handleLog)
	http.HandleFunc("/api/events", apiConf.handleEvents)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))

	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
}
