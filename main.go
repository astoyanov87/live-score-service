package main

import (
	"encoding/json"
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

	// Run a loop to send scraped data to the client
	for {
		select {
		case <-ticker.C:
			// Scrape data
			score, err := handlers.FetchLiveScore("e8f0ffdd-56ee-4128-8e56-552d6d0ed069")
			if err != nil {
				fmt.Fprintf(w, "data: Error scraping data: %v\n\n", err)
			} else {

				data := struct {
					HomePlayerFrames              int `json:"homePlayerFrames"`
					HomePlayerScoreInCurrentFrame int `json:"homePlayerPointsInCurrentFrame"`
					HomePlayerCurrentBreak        int `json:"homePlayerCurrentBreak"`
					AwayPlayerFrames              int `json:"awayPlayerFrames"`
					AwayPlayerScoreInCurrentFrame int `json:"awayPlayerPointsInCurrentFrame"`
					AwayPlayerCurrentBreak        int `json:"awayPlayerCurrentBreak"`
				}{
					HomePlayerFrames:              score.HomePlayerFrames,
					HomePlayerScoreInCurrentFrame: score.HomeplayerPointsInCurrentFrame,
					HomePlayerCurrentBreak:        score.HomePlayerCurrentBreak,
					AwayPlayerFrames:              score.AwayPlayerFrames,
					AwayPlayerScoreInCurrentFrame: score.AwayPlayerPointsInCurrentFrame,
					AwayPlayerCurrentBreak:        score.AwayPlayerCurrentBreak,
				}
				jsonData, err := json.Marshal(data)
				if err != nil {
					log.Println("Error marshaling JSON:", err)
					return
				}
				fmt.Printf("Sending data to client: %d", score.HomeplayerPointsInCurrentFrame)
				fmt.Fprintf(w, "data: %s\n\n", jsonData)
			}
			flusher.Flush()

			// Check if the client has closed the connection
			if r.Context().Err() != nil {
				log.Println("Client closed connection")
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/events", handleSse)

	fmt.Println("Server is listening on port 8888...")
	log.Fatal(http.ListenAndServe(":8888", nil))
}
