package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CaptainKills/xtream-api/api"
)

const (
	// Mandatory Environment Variables
	ENV_URL      = "XTREAM_URL"
	ENV_USERNAME = "XTREAM_USERNAME"
	ENV_PASSWORD = "XTREAM_PASSWORD"

	// Optional Environment Variables
	ENV_LAUNCH   = "XTREAM_LAUNCH"
	ENV_IMAGES   = "XTREAM_IMAGES"
	ENV_METADATA = "XTREAM_METADATA"
	ENV_REQUESTS = "XTREAM_REQUESTS"
	ENV_TIMEOUT  = "XTREAM_TIMEOUT"

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

func GetEnvironmentOptions() api.XtreamOptions {
	var options api.XtreamOptions

	// Launch Time
	launch := os.Getenv(ENV_LAUNCH)
	if launch != "" {
		t, err := time.Parse(time.TimeOnly, launch)
		if err != nil {
			log.Printf("[WARNING] '%s' Environment Variable Invalid! %v\n", ENV_LAUNCH, err)
			options.LaunchTime = time.Date(0, 0, 1, 0, 0, 0, 0, time.Now().Location())
		} else {
			now := time.Now()
			options.LaunchTime = time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), 0, now.Location())

			if now.After(options.LaunchTime) {
				options.LaunchTime = options.LaunchTime.AddDate(0, 0, 1)
			}
		}
	} else {
		options.LaunchTime = time.Date(0, 0, 1, 0, 0, 0, 0, time.Now().Location())
	}

	// Images Enabled
	images := os.Getenv(ENV_IMAGES)
	if images != "" {
		switch images {
		case "true":
			options.ImagesEnabled = true
		case "false":
			options.ImagesEnabled = false
		default:
			log.Printf("[WARNING] '%s' Environment Variable Invalid!\n", ENV_IMAGES)
			options.ImagesEnabled = false
		}
	} else {
		options.ImagesEnabled = false
	}

	// NFO Enabled
	metadata := os.Getenv(ENV_METADATA)
	if metadata != "" {
		switch metadata {
		case "true":
			options.MetadataEnabled = true
		case "false":
			options.MetadataEnabled = false
		default:
			log.Printf("[WARNING] '%s' Environment Variable Invalid!\n", ENV_METADATA)
			options.MetadataEnabled = false
		}
	} else {
		options.MetadataEnabled = false
	}

	// Request Per Minute
	requests := os.Getenv(ENV_REQUESTS)
	if requests == "" {
		options.RequestPerMinute = time.Duration(1000)
	} else {
		r, err := strconv.Atoi(requests)
		if err != nil {
			log.Printf("[WARNING] '%s' Environment Variable Invalid! %v\n", ENV_REQUESTS, err)
			options.RequestPerMinute = time.Duration(1000)
		} else {
			options.RequestPerMinute = time.Duration(r)
		}
	}

	// Request Timeout
	timeout := os.Getenv(ENV_TIMEOUT)
	if timeout == "" {
		options.RequestTimeout = time.Duration(30)
	} else {
		t, err := strconv.Atoi(timeout)
		if err != nil {
			log.Printf("[WARNING] '%s' Environment Variable Invalid! %v\n", ENV_TIMEOUT, err)
			options.RequestTimeout = time.Duration(30)
		} else {
			options.RequestTimeout = time.Duration(t)
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
