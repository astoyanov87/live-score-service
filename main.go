package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/astoyanov87/live-score-service/handlers"
)

func enableCORS(w http.ResponseWriter) {
	// Allow all origins
	w.Header().Set("Access-Control-Allow-Origin", "*")
	// Allow specific methods
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	// Allow specific headers
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func handleSse(w http.ResponseWriter, r *http.Request) {

	enableCORS(w)

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// Scrape live score

			// var test int
			// test = 18
			score, err := handlers.FetchLiveScore("990f7d4d-d6a9-4054-a00f-5ebd93cd23d6")
			if err != nil {
				fmt.Fprintf(w, "data: Error scraping data: %v\n\n", err)
			} else {
				fmt.Printf("Sending data to client: %d", score.HomePlayerCurrentBreak)
				fmt.Fprintf(w, "data: %d\n\n", score.HomePlayerCurrentBreak) // Send live score to client
			}
			flusher.Flush()

			// Check if client has disconnected
			if f, ok := w.(http.CloseNotifier); ok {
				select {
				case <-f.CloseNotify():
					log.Println("Client disconnected.")
					return
				default:
				}
			}
		}
	}
}

func main() {
	http.HandleFunc("/events", handleSse)

	fmt.Println("Server is listening on port 8888...")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
