package main

import (
	"log"
	"time"

	"github.com/CaptainKills/xtream-api/api"
)

const (
	directoryRoot        = "media/"
	directoryLivestreams = directoryRoot + "live/"
	directoryMovies      = directoryRoot + "movies/"
	directorySeries      = directoryRoot + "series/"

	directoryCache = "cache/"

	debugPercent = 100
)

func main() {
	// Environment Variables
	url, username, password := GetEnvironmentCredentials()
	options := GetEnvironmentOptions()

	if url == "" || username == "" || password == "" {
		log.Fatalln("[ERROR] Missing Environment Variables! Exiting Program...")
	}

	// Xtream Client
	client := api.NewClient(url, username, password, options)
	_, err := client.GetAccountInfo()
	if err != nil {
		log.Fatalf("[ERROR] Authentication Failed: %v\n", err)
	}
	log.Printf("[INFO] Authentication Successful: %s\n", url)

	for {
		start := time.Now()

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

		// Filter Categories
		err = FilterCategories(&client.Options, &live_categories, &movie_categories, &series_categories)
		if err != nil {
			log.Printf("[ERROR] Failted to filter categories: %v\n", err)
		}

		// Filter Streams
		err = FilterLiveStreams(client, &livestreams, live_categories)
		if err != nil {
			log.Printf("[ERROR] Failed to filter Livestreams: %v\n", err)
		}

		err = FilterMovies(client, &movies, movie_categories)
		if err != nil {
			log.Printf("[ERROR] Failed to filter Movies: %v\n", err)
		}

		err = FilterSeries(client, &series, series_categories)
		if err != nil {
			log.Printf("[ERROR] Failed to filter Series: %v\n", err)
		}

		// Export LiveStreams
		err = ExportLiveStreams(client, &livestreams)
		if err != nil {
			log.Printf("[ERROR] Unable to export LiveStreams: %v\n", err)
		}

		err = ValidateLiveStreams()
		if err != nil {
			log.Printf("[ERROR] Unable to validate LiveStreams: %v\n", err)
		}

		// Export Movies
		err = ExportMovies(client, &movies)
		if err != nil {
			log.Printf("[ERROR] Unable to export Movies: %v\n", err)
		}

		err = ValidateMovies()
		if err != nil {
			log.Printf("[ERROR] Unable to validate Movies: %v\n", err)
		}

		// Export Series
		err = ExportSeries(client, &series)
		if err != nil {
			log.Printf("[ERROR] Unable to export Series: %v\n", err)
		}

		err = ValidateSeries()
		if err != nil {
			log.Printf("[ERROR] Unable to validate Series: %v\n", err)
		}

		// Wait until next run
		end := time.Now()
		diff := end.Sub(start).Round(time.Millisecond)
		next := start.Add(24 * time.Hour)

		log.Printf("[INFO] Run Duration: %s. Next run scheduled at: %s\n", diff.String(), next.Format(time.DateTime))
		time.Sleep(time.Until(next))
	}
}
