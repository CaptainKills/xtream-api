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

	ENV_DISABLED_LIVE   = "XTREAM_DISABLED_LIVE"
	ENV_DISABLED_MOVIES = "XTREAM_DISABLED_MOVIES"
	ENV_DISABLED_SERIES = "XTREAM_DISABLED_SERIES"

	ENV_BANNED_LIVE   = "XTREAM_BANNED_LIVE"
	ENV_BANNED_MOVIES = "XTREAM_BANNED_MOVIES"
	ENV_BANNED_SERIES = "XTREAM_BANNED_SERIES"
)

func GetCredentials() (string, string, string) {
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

func GetApplicationConfig() Config {
	var config Config

	// Launch Time
	launch := os.Getenv(ENV_LAUNCH)
	if launch != "" {
		t, err := time.Parse(time.TimeOnly, launch)
		if err != nil {
			log.Printf("[WARNING] '%s' Environment Variable Invalid! %v\n", ENV_LAUNCH, err)
			config.LaunchTime = time.Date(0, 0, 1, 0, 0, 0, 0, time.Now().Location())
		} else {
			now := time.Now()
			config.LaunchTime = time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), t.Second(), 0, now.Location())

			if now.After(config.LaunchTime) {
				config.LaunchTime = config.LaunchTime.AddDate(0, 0, 1)
			}
		}
	} else {
		config.LaunchTime = time.Date(0, 0, 1, 0, 0, 0, 0, time.Now().Location())
	}

	// Enabled Programs
	disabled := os.Getenv(ENV_DISABLED_LIVE)
	if disabled != "" {
		switch disabled {
		case "true":
			config.DisabledLive = true
		case "false":
			config.DisabledLive = false
		default:
			log.Printf("[WARNING] '%s' Environment Variable Invalid!\n", ENV_DISABLED_LIVE)
			config.DisabledLive = false
		}
	} else {
		config.DisabledLive = false
	}

	disabled = os.Getenv(ENV_DISABLED_MOVIES)
	if disabled != "" {
		switch disabled {
		case "true":
			config.DisabledMovies = true
		case "false":
			config.DisabledMovies = false
		default:
			log.Printf("[WARNING] '%s' Environment Variable Invalid!\n", ENV_DISABLED_MOVIES)
			config.DisabledMovies = false
		}
	} else {
		config.DisabledMovies = false
	}

	disabled = os.Getenv(ENV_DISABLED_SERIES)
	if disabled != "" {
		switch disabled {
		case "true":
			config.DisabledSeries = true
		case "false":
			config.DisabledSeries = false
		default:
			log.Printf("[WARNING] '%s' Environment Variable Invalid!\n", ENV_DISABLED_SERIES)
			config.DisabledSeries = false
		}
	} else {
		config.DisabledSeries = false
	}

	// Banned Categories
	banned := os.Getenv(ENV_BANNED_LIVE)
	config.BannedLiveStreams = strings.Split(banned, ",")

	banned = os.Getenv(ENV_BANNED_MOVIES)
	config.BannedMovies = strings.Split(banned, ",")

	banned = os.Getenv(ENV_BANNED_SERIES)
	config.BannedSeries = strings.Split(banned, ",")

	return config
}

func GetXtreamOptions() api.XtreamOptions {
	var options api.XtreamOptions

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

	return options
}
