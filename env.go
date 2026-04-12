package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	xtream "github.com/CaptainKills/xtream-api/api"
)

const (
	// Mandatory Environment Variables
	ENV_URL      = "XTREAM_URL"
	ENV_USERNAME = "XTREAM_USERNAME"
	ENV_PASSWORD = "XTREAM_PASSWORD"

	// Optional Environment Variables
	ENV_IMAGES   = "XTREAM_IMAGES"
	ENV_NFO      = "XTREAM_NFO"
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

func GetEnvironmentOptions() xtream.XtreamOptions {
	var options xtream.XtreamOptions

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
	nfo := os.Getenv(ENV_NFO)
	if nfo != "" {
		switch nfo {
		case "true":
			options.NfoEnabled = true
		case "false":
			options.NfoEnabled = false
		default:
			log.Printf("[WARNING] '%s' Environment Variable Invalid!\n", ENV_NFO)
			options.NfoEnabled = false
		}
	} else {
		options.NfoEnabled = false
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
