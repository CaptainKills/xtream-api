package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	xtream "github.com/CaptainKills/xtream-api/api"
)

const (
	// Mandatory Environment Variables
	ENV_URL      = "XTREAM_URL"
	ENV_USERNAME = "XTREAM_USERNAME"
	ENV_PASSWORD = "XTREAM_PASSWORD"

	// Optional Environment Variables
	ENV_IMAGES   = "XTREAM_IMAGES"
	ENV_REQUESTS = "XTREAM_REQUESTS"

	ENV_BANNED_LIVE   = "XTREAM_BANNED_LIVE"
	ENV_BANNED_MOVIES = "XTREAM_BANNED_MOVIES"
	ENV_BANNED_SERIES = "XTREAM_BANNED_SERIES"
)

func GetEnvironmentCredentials() (string, string, string) {
	// Xtream URL
	url := os.Getenv(ENV_URL)
	if url == "" {
		log.Printf("[ERROR] '%s' Environment Variable Not Specified!\n", ENV_URL)
	}

	// Xtream Username
	username := os.Getenv(ENV_USERNAME)
	if username == "" {
		log.Printf("[ERROR] '%s' Environment Variable Not Specified!\n", ENV_USERNAME)
	}

	// Xtream Password
	password := os.Getenv(ENV_PASSWORD)
	if password == "" {
		log.Printf("[ERROR] '%s' Environment Variable Not Specified!\n", ENV_PASSWORD)
	}

	return url, username, password
}

func GetEnvironmentOptions() xtream.XtreamOptions {
	var options xtream.XtreamOptions
	var err error

	// Images Enabled
	images := os.Getenv(ENV_IMAGES)
	switch images {
	case "true":
		options.ImagesEnabled = true
	case "false":
		options.ImagesEnabled = false
	default:
		log.Printf("[WARNING] '%s' Environment Variable Invalid!\n", ENV_IMAGES)
		options.ImagesEnabled = true
	}

	// Request Per Minute
	requests := os.Getenv(ENV_REQUESTS)
	if requests == "" {
		options.RequestPerMinute = 1000
	} else {
		options.RequestPerMinute, err = strconv.Atoi(requests)
		if err != nil {
			log.Printf("[WARNING] '%s' Environment Variable Invalid!\n", ENV_REQUESTS)
			options.RequestPerMinute = 1000
		}
	}

	// Banned LiveStream Categories
	banned := os.Getenv(ENV_BANNED_LIVE)
	options.BannedLiveStreams = strings.Split(banned, ",")

	// Banned Movie Categories
	banned = os.Getenv(ENV_BANNED_MOVIES)
	options.BannedMovies = strings.Split(banned, ",")

	// Banned Series Categories
	banned = os.Getenv(ENV_BANNED_SERIES)
	options.BannedSeries = strings.Split(banned, ",")

	return options
}
