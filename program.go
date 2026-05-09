package main

import (
	"log"
	"time"

	"github.com/CaptainKills/xtream-api/api"
)

type Program[T api.Stream] struct {
	disabled func(Config) bool

	loader      func(*api.XtreamClient) (map[int]T, error)
	categoriser func(*api.XtreamClient) (map[int]api.Category, error)
	moderator   func(Config) []string

	cacheFile string
	stateFile string
	directory string

	label string
}

type Config struct {
	LaunchTime time.Time

	DisabledLive   bool
	DisabledMovies bool
	DisabledSeries bool

	BannedLiveStreams []string
	BannedMovies      []string
	BannedSeries      []string
}

type State struct {
	Strm  bool `json:"strm"`
	Nfo   bool `json:"nfo"`
	Image bool `json:"image"`
}

var livestreams = Program[api.LiveStream]{
	disabled: func(c Config) bool {
		return c.DisabledLive
	},

	loader: func(c *api.XtreamClient) (map[int]api.LiveStream, error) {
		return c.GetLiveStreams()
	},
	categoriser: func(c *api.XtreamClient) (map[int]api.Category, error) {
		return c.GetLiveStreamCategories()
	},
	moderator: func(c Config) []string {
		return c.BannedLiveStreams
	},

	cacheFile: cacheFileLivestreams,
	stateFile: metadataFileLivestreams,
	directory: directoryLivestreams,

	label: "Livestreams",
}

var movies = Program[api.Movie]{
	disabled: func(c Config) bool {
		return c.DisabledMovies
	},

	loader: func(c *api.XtreamClient) (map[int]api.Movie, error) {
		return c.GetMovies()
	},
	categoriser: func(c *api.XtreamClient) (map[int]api.Category, error) {
		return c.GetMovieCategories()
	},
	moderator: func(c Config) []string {
		return c.BannedMovies
	},

	cacheFile: cacheFileMovies,
	stateFile: metadataFileMovies,
	directory: directoryMovies,

	label: "Movies",
}

var series = Program[api.Series]{
	disabled: func(c Config) bool {
		return c.DisabledSeries
	},

	loader: func(c *api.XtreamClient) (map[int]api.Series, error) {
		return c.GetSeries()
	},
	categoriser: func(c *api.XtreamClient) (map[int]api.Category, error) {
		return c.GetSeriesCategories()
	},
	moderator: func(c Config) []string {
		return c.BannedSeries
	},

	cacheFile: cacheFileSeries,
	stateFile: metadataFileSeries,
	directory: directorySeries,

	label: "Series",
}

func (p *Program[T]) Run(client *api.XtreamClient, config Config) error {
	var err error

	// Check if program is disabled
	if p.disabled(config) {
		log.Printf("[INFO] Program Disabled: %s\n", p.label)
		return nil
	}
	log.Printf("[INFO] Running Program: %s\n", p.label)

	// Initialise Program
	streams := map[int]T{}
	categories := map[int]api.Category{}
	banned := []string{}
	cache := map[int]T{}
	state := map[int]*State{}

	// Fetch Streams
	streams, err = p.loader(client)
	if err != nil {
		log.Printf("[ERROR] (%s) Failed to retrieve Streams: %v\n", p.label, err)
	}
	log.Printf("[INFO] (%s) Found %6d Streams\n", p.label, len(streams))

	// Fetch Categories
	categories, err = p.categoriser(client)
	if err != nil {
		log.Printf("[ERROR] (%s) Failed to retrieve Categories: %v\n", p.label, err)
	}
	log.Printf("[INFO] (%s) Found %6d Categories\n", p.label, len(categories))

	// Fetch Banned
	banned = p.moderator(config)
	log.Printf("[INFO] (%s) Found %6d Banned Prefixes\n", p.label, len(banned))

	// Import Cache
	err = ImportCache(&cache, p.cacheFile, p.label)
	if err != nil {
		log.Printf("[ERROR] (%s) Failed to import Cache: %v\n", p.label, err)
	}

	// Import Metadata
	state, err = LoadState(p.stateFile, p.label)
	if err != nil {
		log.Printf("[ERROR] (%s) Failed to load State: %v\n", p.label, err)
	}

	// Filter Streams
	err = FilterCategories(&categories, banned, p.label)
	if err != nil {
		log.Printf("[ERROR] (%s) Failed to filter Categories: %v\n", p.label, err)
	}

	err = FilterStreams(client, &streams, categories, cache, &state, p.label)
	if err != nil {
		log.Printf("[ERROR] (%s) Failed to filter Streams: %v\n", p.label, err)
	}

	// Export Streams
	err = Export(client, &streams, &state, p.directory, p.label)
	if err != nil {
		log.Printf("[ERROR] (%s) Unable to export Streams: %v\n", p.label, err)
	}

	// Validate Streams
	err = Validate(p.directory, p.label)
	if err != nil {
		log.Printf("[ERROR] (%s) Unable to validate Streams: %v\n", p.label, err)
	}

	// Export Cache
	err = ExportCache(&streams, &cache, &state, p.cacheFile, p.label)
	if err != nil {
		log.Printf("[ERROR] (%s) Failed to export Cache: %v\n", p.label, err)
	}

	// Export Metadata
	err = SaveState(&state, p.stateFile, p.label)
	if err != nil {
		log.Printf("[ERROR] (%s) Failed to save State: %v\n", p.label, err)
	}

	return nil
}
