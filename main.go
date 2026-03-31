package main

import (
	"log"
	"os"

	xtream "github.com/CaptainKills/xtream-api/api"
)

const (
	ENV_URL      = "XTREAM_URL"
	ENV_USERNAME = "XTREAM_USERNAME"
	ENV_PASSWORD = "XTREAM_PASSWORD"
)

func main() {
	// Environment Variables
	url := os.Getenv(ENV_URL)
	if url == "" {
		log.Printf("[ERROR] '%s' Environment Variable Not Specified!\n", ENV_URL)
	}

	username := os.Getenv(ENV_USERNAME)
	if username == "" {
		log.Printf("[ERROR] '%s' Environment Variable Not Specified!\n", ENV_USERNAME)
	}

	password := os.Getenv(ENV_PASSWORD)
	if password == "" {
		log.Printf("[ERROR] '%s' Environment Variable Not Specified!\n", ENV_PASSWORD)
	}

	if url == "" || username == "" || password == "" {
		log.Fatalln("[ERROR] Missing Environment Variables! Exiting Program...")
	}

	// Xtream Client
	client := xtream.NewClient(url, username, password)
	_, err := client.GetAccountInfo()
	if err != nil {
		log.Fatalf("[ERROR] Authentication Failed: %q\n", err)
	}
	log.Printf("[INFO] Authentication Successful: %s\n", url)

	// Fetch Categories
	live_categories, err := client.GetLiveStreamCategories()
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve LiveStream Categories: %q\n", err)
	}
	log.Printf("[INFO] Found %6d Categories for Livestreams\n", len(live_categories))

	movie_categories, err := client.GetMovieCategories()
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve Movie Categories: %q\n", err)
	}
	log.Printf("[INFO] Found %6d Categories for Movies\n", len(movie_categories))

	series_categories, err := client.GetSeriesCategories()
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve Series Categories: %q\n", err)
	}
	log.Printf("[INFO] Found %6d Categories for Series\n", len(series_categories))

	// Fetch Streams
	livestreams, err := client.GetLiveStreams()
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve LiveStreams: %q\n", err)
	}
	log.Printf("[INFO] Found %6d Livestreams\n", len(livestreams))

	movies, err := client.GetMovies()
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve Movies: %q\n", err)
	}
	log.Printf("[INFO] Found %6d Movies\n", len(movies))

	series, err := client.GetSeries()
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve Series: %q\n", err)
	}
	log.Printf("[INFO] Found %6d Series\n", len(series))

	// Export Streams
	err = client.ExportLiveStreams()
	if err != nil {
		log.Printf("[ERROR] Unable to export LiveStreams: %q\n", err)
	}

	err = client.ExportMovies()
	if err != nil {
		log.Printf("[ERROR] Unable to export Movies: %q\n", err)
	}

	err = client.ExportSeries()
	if err != nil {
		log.Printf("[ERROR] Unable to export Series: %q\n", err)
	}
}
