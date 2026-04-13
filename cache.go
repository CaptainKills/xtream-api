package main

import (
	"encoding/json"
	"os"

	"github.com/CaptainKills/xtream-api/api"
	"github.com/CaptainKills/xtream-api/utils"
)

const (
	fileLiveCategories   = "live_categories.json"
	fileMovieCategories  = "movie_categories.json"
	fileSeriesCategories = "series_categories.json"

	fileLiveStreams = "livestreams.json"
	fileMovies      = "movies.json"
	fileSeries      = "series.json"
)

func ExportCache(c *api.XtreamClient) error {
	// Create Root Directory
	err := os.MkdirAll(api.DirectoryCache, 0o750)
	if err != nil {
		return err
	}

	// Export Categories
	err = utils.WriteJson(api.DirectoryCache+fileLiveCategories, c.Data.LiveCategories)
	if err != nil {
		return err
	}

	err = utils.WriteJson(api.DirectoryCache+fileMovieCategories, c.Data.MovieCategories)
	if err != nil {
		return err
	}

	err = utils.WriteJson(api.DirectoryCache+fileSeriesCategories, c.Data.SeriesCategories)
	if err != nil {
		return err
	}

	// Export Streams
	err = utils.WriteJson(api.DirectoryCache+fileLiveStreams, c.Data.Livestreams)
	if err != nil {
		return err
	}

	err = utils.WriteJson(api.DirectoryCache+fileMovies, c.Data.Movies)
	if err != nil {
		return err
	}

	err = utils.WriteJson(api.DirectoryCache+fileSeries, c.Data.Series)
	if err != nil {
		return err
	}

	return nil
}

func ImportCache(c *api.XtreamClient) error {
	c.Old.LiveCategories = map[int]api.Category{}
	c.Old.MovieCategories = map[int]api.Category{}
	c.Old.SeriesCategories = map[int]api.Category{}
	c.Old.Livestreams = map[int]api.LiveStream{}
	c.Old.Movies = map[int]api.Movie{}
	c.Old.Series = map[int]api.Series{}

	// Import Categories: LiveStreams
	file, err := os.ReadFile(api.DirectoryCache + fileLiveCategories)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if string(file) == "" {
		file = []byte("{}") // Create Empty JSON Object if file doesn't exist.
	}

	err = json.Unmarshal(file, &c.Old.LiveCategories)
	if err != nil {
		return err
	}

	// Import Categories: Movies
	file, err = os.ReadFile(api.DirectoryCache + fileMovieCategories)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if string(file) == "" {
		file = []byte("{}") // Create Empty JSON Object if file doesn't exist.
	}

	err = json.Unmarshal(file, &c.Old.MovieCategories)
	if err != nil {
		return err
	}

	// Import Categories: Series
	file, err = os.ReadFile(api.DirectoryCache + fileSeriesCategories)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if string(file) == "" {
		file = []byte("{}") // Create Empty JSON Object if file doesn't exist.
	}

	err = json.Unmarshal(file, &c.Old.SeriesCategories)
	if err != nil {
		return err
	}

	// Import Streams: LiveStreams
	file, err = os.ReadFile(api.DirectoryCache + fileLiveStreams)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if string(file) == "" {
		file = []byte("{}") // Create Empty JSON Object if file doesn't exist.
	}

	err = json.Unmarshal(file, &c.Old.Livestreams)
	if err != nil {
		return err
	}

	// Import Streams: Movies
	file, err = os.ReadFile(api.DirectoryCache + fileMovies)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if string(file) == "" {
		file = []byte("{}") // Create Empty JSON Object if file doesn't exist.
	}

	err = json.Unmarshal(file, &c.Old.Movies)
	if err != nil {
		return err
	}

	// Import Streams: Series
	file, err = os.ReadFile(api.DirectoryCache + fileSeries)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if string(file) == "" {
		file = []byte("{}") // Create Empty JSON Object if file doesn't exist.
	}

	err = json.Unmarshal(file, &c.Old.Series)
	if err != nil {
		return err
	}

	return nil
}
