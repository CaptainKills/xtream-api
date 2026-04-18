package main

import (
	"log"
	"time"

	"github.com/CaptainKills/xtream-api/api"
)

const (
	MODE_DIR  = 0o755
	MODE_FILE = 0o644

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

	// Launch Time Delay
	if time.Now().Before(options.LaunchTime) {
		launch := options.LaunchTime.Format(time.DateTime)
		now := time.Now().Format(time.DateTime)

		log.Printf("[INFO] Next run scheduled at: %s (Current time: %s)\n", launch, now)
		time.Sleep(time.Until(options.LaunchTime))
	}

	for {
		start := time.Now()

		// Fetch Categories
		liveCategories, err := client.GetLiveStreamCategories()
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve LiveStream Categories: %v\n", err)
		}
		log.Printf("[INFO] Found %6d Categories for Livestreams\n", len(liveCategories))

		movieCategories, err := client.GetMovieCategories()
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve Movie Categories: %v\n", err)
		}
		log.Printf("[INFO] Found %6d Categories for Movies\n", len(movieCategories))

		seriesCategories, err := client.GetSeriesCategories()
		if err != nil {
			log.Printf("[ERROR] Failed to retrieve Series Categories: %v\n", err)
		}
		log.Printf("[INFO] Found %6d Categories for Series\n", len(seriesCategories))

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

		// Filter Streams
		err = FilterCategories(&client.Options, &liveCategories, &movieCategories, &seriesCategories)
		if err != nil {
			log.Printf("[ERROR] Failted to filter categories: %v\n", err)
		}

		err = FilterLiveStreams(client, &livestreams, liveCategories)
		if err != nil {
			log.Printf("[ERROR] Failed to filter Livestreams: %v\n", err)
		}

		err = FilterMovies(client, &movies, movieCategories)
		if err != nil {
			log.Printf("[ERROR] Failed to filter Movies: %v\n", err)
		}

		err = FilterSeries(client, &series, seriesCategories)
		if err != nil {
			log.Printf("[ERROR] Failed to filter Series: %v\n", err)
		}

		// Export LiveStreams
		err = Export(client, &livestreams, directoryLivestreams, "Livestreams")
		if err != nil {
			log.Printf("[ERROR] Unable to export LiveStreams: %v\n", err)
		}

		err = Validate(directoryLivestreams, "Livestreams")
		if err != nil {
			log.Printf("[ERROR] Unable to validate LiveStreams: %v\n", err)
		}

		// Export Movies
		err = Export(client, &movies, directoryMovies, "  Movies   ")
		if err != nil {
			log.Printf("[ERROR] Unable to export Movies: %v\n", err)
		}

		err = Validate(directoryMovies, "  Movies   ")
		if err != nil {
			log.Printf("[ERROR] Unable to validate Movies: %v\n", err)
		}

		// Export Series
		err = Export(client, &series, directorySeries, "  Series   ")
		if err != nil {
			log.Printf("[ERROR] Unable to export Series: %v\n", err)
		}

		err = Validate(directorySeries, "  Series   ")
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
