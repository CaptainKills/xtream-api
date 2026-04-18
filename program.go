package main

import (
	"log"

	"github.com/CaptainKills/xtream-api/api"
)

type Program[T api.Stream] struct {
	loader      func(*api.XtreamClient) (map[int]T, error)
	categoriser func(*api.XtreamClient) (map[int]api.Category, error)
	moderator   func(*api.XtreamClient) []string

	streams    map[int]T
	categories map[int]api.Category
	banned     []string

	cache        map[int]T
	cacheFile    string
	metadata     map[int]*api.XtreamMetadata
	metadataFile string

	directory string
	label     string
}

var livestreams = Program[api.LiveStream]{
	loader: func(c *api.XtreamClient) (map[int]api.LiveStream, error) {
		return c.GetLiveStreams()
	},
	categoriser: func(c *api.XtreamClient) (map[int]api.Category, error) {
		return c.GetLiveStreamCategories()
	},
	moderator: func(c *api.XtreamClient) []string {
		return c.Options.BannedLiveStreams
	},

	streams:    map[int]api.LiveStream{},
	categories: map[int]api.Category{},
	banned:     []string{},

	cache:     map[int]api.LiveStream{},
	cacheFile: cacheFileLivestreams,

	metadata:     map[int]*api.XtreamMetadata{},
	metadataFile: metadataFileLivestreams,

	directory: directoryLivestreams,
	label:     "Livestreams",
}

var movies = Program[api.Movie]{
	loader: func(c *api.XtreamClient) (map[int]api.Movie, error) {
		return c.GetMovies()
	},
	categoriser: func(c *api.XtreamClient) (map[int]api.Category, error) {
		return c.GetMovieCategories()
	},
	moderator: func(c *api.XtreamClient) []string {
		return c.Options.BannedMovies
	},

	streams:    map[int]api.Movie{},
	categories: map[int]api.Category{},
	banned:     []string{},

	cache:     map[int]api.Movie{},
	cacheFile: cacheFileMovies,

	metadata:     map[int]*api.XtreamMetadata{},
	metadataFile: metadataFileMovies,

	directory: directoryMovies,
	label:     "Movies",
}

var series = Program[api.Series]{
	loader: func(c *api.XtreamClient) (map[int]api.Series, error) {
		return c.GetSeries()
	},
	categoriser: func(c *api.XtreamClient) (map[int]api.Category, error) {
		return c.GetSeriesCategories()
	},
	moderator: func(c *api.XtreamClient) []string {
		return c.Options.BannedSeries
	},

	streams:    map[int]api.Series{},
	categories: map[int]api.Category{},
	banned:     []string{},

	cache:     map[int]api.Series{},
	cacheFile: cacheFileSeries,

	metadata:     map[int]*api.XtreamMetadata{},
	metadataFile: metadataFileSeries,

	directory: directorySeries,
	label:     "Series",
}

func (p Program[T]) Run(client *api.XtreamClient) error {
	var err error

	// Fetch Streams
	p.streams, err = p.loader(client)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve %s: %v\n", p.label, err)
	}
	log.Printf("[INFO] Found %6d %s\n", len(p.streams), p.label)

	// Fetch Categories
	p.categories, err = p.categoriser(client)
	if err != nil {
		log.Printf("[ERROR] Failed to retrieve %s Categories: %v\n", p.label, err)
	}
	log.Printf("[INFO] Found %6d Categories for %s\n", len(p.categories), p.label)

	// Fetch Banned
	p.banned = p.moderator(client)

	// Import Cache
	err = ImportCache(&p.cache, p.cacheFile, p.label)
	if err != nil {
		log.Printf("[ERROR] Failed to import %s Cache: %v\n", p.label, err)
	}

	// Import Metadata
	p.metadata, err = ImportMetadata(p.metadataFile, p.label)
	if err != nil {
		log.Printf("[ERROR] Failed to import %s Metadata: %v\n", p.label, err)
	}

	// Filter Streams
	err = FilterCategories(&p.categories, p.banned, p.label)
	if err != nil {
		log.Printf("[ERROR] Failted to filter categories: %v\n", err)
	}

	err = FilterStreams(client, &p.streams, p.categories, p.cache, &p.metadata, p.label)
	if err != nil {
		log.Printf("[ERROR] Failed to filter %s: %v\n", p.label, err)
	}

	// Export Streams
	err = Export(client, &p.streams, &p.metadata, p.directory, p.label)
	if err != nil {
		log.Printf("[ERROR] Unable to export %s: %v\n", p.label, err)
	}

	// Validate Streams
	err = Validate(p.directory, p.label)
	if err != nil {
		log.Printf("[ERROR] Unable to validate %s: %v\n", p.label, err)
	}

	// Export Cache
	err = ExportCache(&p.streams, &p.cache, &p.metadata, p.cacheFile, p.label)
	if err != nil {
		log.Printf("[ERROR] Failed to export %s Cache: %v\n", p.label, err)
	}

	// Export Metadata
	err = ExportMetadata(&p.metadata, p.metadataFile, p.label)
	if err != nil {
		log.Printf("[ERROR] Failed to export %s Metadata: %v\n", p.label, err)
	}

	return nil
}
