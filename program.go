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

	// Reset Program
	p.streams = map[int]T{}
	p.categories = map[int]api.Category{}
	p.banned = []string{}
	p.cache = map[int]T{}
	p.metadata = map[int]*api.XtreamMetadata{}

	// Fetch Streams
	p.streams, err = p.loader(client)
	if err != nil {
		log.Printf("[ERROR] (%s) Failed to retrieve Streams: %v\n", p.label, err)
	}
	log.Printf("[INFO] (%s) Found %6d Streams\n", p.label, len(p.streams))

	// Fetch Categories
	p.categories, err = p.categoriser(client)
	if err != nil {
		log.Printf("[ERROR] (%s) Failed to retrieve Categories: %v\n", p.label, err)
	}
	log.Printf("[INFO] (%s) Found %6d Categories\n", p.label, len(p.categories))

	// Fetch Banned
	p.banned = p.moderator(client)
	log.Printf("[INFO] (%s) Found %6d Banned Prefixes\n", p.label, len(p.banned))

	// Import Cache
	err = ImportCache(&p.cache, p.cacheFile, p.label)
	if err != nil {
		log.Printf("[ERROR] (%s) Failed to import Cache: %v\n", p.label, err)
	}

	// Import Metadata
	p.metadata, err = ImportMetadata(p.metadataFile, p.label)
	if err != nil {
		log.Printf("[ERROR] (%s) Failed to import Metadata: %v\n", p.label, err)
	}

	// Filter Streams
	err = FilterCategories(&p.categories, p.banned, p.label)
	if err != nil {
		log.Printf("[ERROR] (%s) Failed to filter Categories: %v\n", p.label, err)
	}

	err = FilterStreams(client, &p.streams, p.categories, p.cache, &p.metadata, p.label)
	if err != nil {
		log.Printf("[ERROR] (%s) Failed to filter Streams: %v\n", p.label, err)
	}

	// Export Streams
	err = Export(client, &p.streams, &p.metadata, p.directory, p.label)
	if err != nil {
		log.Printf("[ERROR] (%s) Unable to export Streams: %v\n", p.label, err)
	}

	// Validate Streams
	err = Validate(p.directory, p.label)
	if err != nil {
		log.Printf("[ERROR] (%s) Unable to validate Streams: %v\n", p.label, err)
	}

	// Export Cache
	err = ExportCache(&p.streams, &p.cache, &p.metadata, p.cacheFile, p.label)
	if err != nil {
		log.Printf("[ERROR] (%s) Failed to export Cache: %v\n", p.label, err)
	}

	// Export Metadata
	err = ExportMetadata(&p.metadata, p.metadataFile, p.label)
	if err != nil {
		log.Printf("[ERROR] (%s) Failed to export Metadata: %v\n", p.label, err)
	}

	// Clear Program
	p.streams = map[int]T{}
	p.categories = map[int]api.Category{}
	p.banned = []string{}
	p.cache = map[int]T{}
	p.metadata = map[int]*api.XtreamMetadata{}

	return nil
}
