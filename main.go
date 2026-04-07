package main

import (
	"log"
	"os"
	"time"

	xtream "github.com/CaptainKills/xtream-api/api"
)

const (
	ENV_URL      = "XTREAM_URL"
	ENV_USERNAME = "XTREAM_USERNAME"
	ENV_PASSWORD = "XTREAM_PASSWORD"
	ENV_IMAGES   = "XTREAM_IMAGES"
)

var ENABLE_IMAGES bool = true

func main() {
	log.SetOutput(os.Stdout)

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

	images := os.Getenv(ENV_IMAGES)
	switch images {
	case "true":
		ENABLE_IMAGES = true
	case "false":
		ENABLE_IMAGES = false
	default:
		ENABLE_IMAGES = true
	}

	if url == "" || username == "" || password == "" {
		log.Fatalln("[ERROR] Missing Environment Variables! Exiting Program...")
	}

	for {
		start := time.Now()

		// Xtream Client
		client := xtream.NewClient(url, username, password)
		_, err := client.GetAccountInfo()
		if err != nil {
			log.Fatalf("[ERROR] Authentication Failed: %v\n", err)
		}
		log.Printf("[INFO] Authentication Successful: %s\n", url)

		// Fetch Categories
		live_categories, err := client.GetLiveStreamCategories()
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve LiveStream Categories: %v\n", err)
		}
		log.Printf("[INFO] Found %6d Categories for Livestreams\n", len(live_categories))

		movie_categories, err := client.GetMovieCategories()
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve Movie Categories: %v\n", err)
		}
		log.Printf("[INFO] Found %6d Categories for Movies\n", len(movie_categories))

		series_categories, err := client.GetSeriesCategories()
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve Series Categories: %v\n", err)
		}
		log.Printf("[INFO] Found %6d Categories for Series\n", len(series_categories))

		// Fetch Streams
		livestreams, err := client.GetLiveStreams()
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve LiveStreams: %v\n", err)
		}
		log.Printf("[INFO] Found %6d Livestreams\n", len(livestreams))

		movies, err := client.GetMovies()
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve Movies: %v\n", err)
		}
		log.Printf("[INFO] Found %6d Movies\n", len(movies))

		series, err := client.GetSeries()
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve Series: %v\n", err)
		}
		log.Printf("[INFO] Found %6d Series\n", len(series))

		// Export LiveStreams
		err = client.ExportLiveStreams(ENABLE_IMAGES)
		if err != nil {
			log.Printf("[ERROR] Unable to export LiveStreams: %v\n", err)
		}

		err = client.ValidateLiveStreams()
		if err != nil {
			log.Printf("[ERROR] Unable to validate LiveStreams: %v\n", err)
		}

		// Export Movies
		err = client.ExportMovies(ENABLE_IMAGES)
		if err != nil {
			log.Printf("[ERROR] Unable to export Movies: %v\n", err)
		}

		err = client.ValidateMovies()
		if err != nil {
			log.Printf("[ERROR] Unable to validate Movies: %v\n", err)
		}

		// Export Series
		err = client.ExportSeries(ENABLE_IMAGES)
		if err != nil {
			log.Printf("[ERROR] Unable to export Series: %v\n", err)
		}

		err = client.ValidateSeries()
		if err != nil {
			log.Printf("[ERROR] Unable to validate Series: %v\n", err)
		}

		// Wait until next run
		end := time.Now()
		diff := end.Sub(start).Round(time.Millisecond)
		next := start.Add(24 * time.Hour)

		log.Printf("[INFO] Run Duration: %s. Next run scheduled at: %s\n", diff.String(), next.Format(time.DateTime))
		time.Sleep(24 * time.Hour)
	}
}
