package api

import (
	"encoding/json"
	"os"
	"strconv"
)

const (
	fileLiveCategories   = "live_categories.json"
	fileMovieCategories  = "movie_categories.json"
	fileSeriesCategories = "series_categories.json"

	fileLiveStreams = "livestreams.json"
	fileMovies      = "movies.json"
	fileSeries      = "series.json"
)

type XtreamRaw struct {
	liveCategories   []byte
	movieCategories  []byte
	seriesCategories []byte
	livestreams      []byte
	movies           []byte
	series           []byte
}

func (c *XtreamClient) ExportCache() error {
	// Create Root Directory
	err := os.MkdirAll(directoryCache, 0o750)
	if err != nil {
		return err
	}

	// Export Categories
	err = WriteJson(directoryCache+fileLiveCategories, c.raw.liveCategories)
	if err != nil {
		return err
	}

	err = WriteJson(directoryCache+fileMovieCategories, c.raw.movieCategories)
	if err != nil {
		return err
	}

	err = WriteJson(directoryCache+fileSeriesCategories, c.raw.seriesCategories)
	if err != nil {
		return err
	}

	// Export Streams
	err = WriteJson(directoryCache+fileLiveStreams, c.raw.livestreams)
	if err != nil {
		return err
	}

	err = WriteJson(directoryCache+fileMovies, c.raw.movies)
	if err != nil {
		return err
	}

	err = WriteJson(directoryCache+fileSeries, c.raw.series)
	if err != nil {
		return err
	}

	return nil
}

func (c *XtreamClient) ImportCache() error {
	c.old.liveCategories = map[int]Category{}
	c.old.movieCategories = map[int]Category{}
	c.old.seriesCategories = map[int]Category{}
	c.old.livestreams = map[int]LiveStream{}
	c.old.movies = map[int]Movie{}
	c.old.series = map[int]Series{}

	// Import Categories: LiveStreams
	file, err := os.ReadFile(directoryCache + fileLiveCategories)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if string(file) == "" {
		file = []byte("[]") // Create Empty JSON Array if file doesn't exist.
	}

	var categories []Category
	err = json.Unmarshal(file, &categories)
	if err != nil {
		return err
	}

	for _, category := range categories {
		id, err := strconv.Atoi(category.Id)
		if err != nil {
			return err
		}

		c.old.liveCategories[id] = category
	}

	// Import Categories: Movies
	file, err = os.ReadFile(directoryCache + fileMovieCategories)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if string(file) == "" {
		file = []byte("[]") // Create Empty JSON Array if file doesn't exist.
	}

	err = json.Unmarshal(file, &categories)
	if err != nil {
		return err
	}

	for _, category := range categories {
		id, err := strconv.Atoi(category.Id)
		if err != nil {
			return err
		}

		c.old.movieCategories[id] = category
	}

	// Import Categories: Series
	file, err = os.ReadFile(directoryCache + fileSeriesCategories)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if string(file) == "" {
		file = []byte("[]") // Create Empty JSON Array if file doesn't exist.
	}

	err = json.Unmarshal(file, &categories)
	if err != nil {
		return err
	}

	for _, category := range categories {
		id, err := strconv.Atoi(category.Id)
		if err != nil {
			return err
		}

		c.old.seriesCategories[id] = category
	}

	// Import Streams: LiveStreams
	var livestreams []LiveStream
	file, err = os.ReadFile(directoryCache + fileLiveStreams)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if string(file) == "" {
		file = []byte("[]") // Create Empty JSON Array if file doesn't exist.
	}

	err = json.Unmarshal(file, &livestreams)
	if err != nil {
		return err
	}

	for _, livestream := range livestreams {
		c.old.livestreams[livestream.Id] = livestream
	}

	// Import Streams: Movies
	var movies []Movie
	file, err = os.ReadFile(directoryCache + fileMovies)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if string(file) == "" {
		file = []byte("[]") // Create Empty JSON Array if file doesn't exist.
	}

	err = json.Unmarshal(file, &movies)
	if err != nil {
		return err
	}

	for _, movie := range movies {
		c.old.movies[movie.Id] = movie
	}

	// Import Streams: Series
	var series []Series
	file, err = os.ReadFile(directoryCache + fileSeries)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if string(file) == "" {
		file = []byte("[]") // Create Empty JSON Array if file doesn't exist.
	}

	err = json.Unmarshal(file, &series)
	if err != nil {
		return err
	}

	for _, show := range series {
		c.old.series[show.Id] = show
	}

	return nil
}
