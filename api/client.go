package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
)

const (
	queryApi = "%s/player_api.php?username=%s&password=%s&action=%s"
	queryEpg = "%s/xmltv.php?username=%s&password=%s"
	queryUrl = "%s://%s:%s/%s/%s/%s/%d.%s"

	actionAccountInfo      = ""
	actionLiveCategories   = "get_live_categories"
	actionMovieCategories  = "get_vod_categories"
	actionSeriesCategories = "get_series_categories"
	actionLivestreams      = "get_live_streams"
	actionMovies           = "get_vod_streams"
	actionSeries           = "get_series"
	actionMovieInfo        = "get_vod_info&vod_id=%d"
	actionSeriesInfo       = "get_series_info&series_id=%d"

	directoryRoot        = "media/"
	directoryLivestreams = directoryRoot + "live/"
	directoryMovies      = directoryRoot + "movies/"
	directorySeries      = directoryRoot + "series/"

	debugPercent = 100
)

type XtreamClient struct {
	url      string
	username string
	password string

	account          Account
	liveCategories   []Category
	movieCategories  []Category
	seriesCategories []Category
	livestreams      []LiveStream
	movies           []Movie
	series           []Series

	rawAccount          []byte
	rawLiveCategories   []byte
	rawMovieCategories  []byte
	rawSeriesCategories []byte
	rawLiveStreams      []byte
	rawMovies           []byte
	rawSeries           []byte
}

func NewClient(url string, username string, password string) *XtreamClient {
	return &XtreamClient{
		url:      url,
		username: username,
		password: password,

		account:          Account{},
		liveCategories:   []Category{},
		movieCategories:  []Category{},
		seriesCategories: []Category{},
		livestreams:      []LiveStream{},
		movies:           []Movie{},
		series:           []Series{},

		rawAccount:          []byte{},
		rawLiveCategories:   []byte{},
		rawMovieCategories:  []byte{},
		rawSeriesCategories: []byte{},
		rawLiveStreams:      []byte{},
		rawMovies:           []byte{},
		rawSeries:           []byte{},
	}
}

func (c *XtreamClient) buildURL(stream string, id int, ext string) (string, error) {
	protocol := c.account.ServerInfo.Protocol
	domain := c.account.ServerInfo.URL

	var port string
	switch protocol {
	case "http":
		port = c.account.ServerInfo.HttpPort
	case "https":
		port = c.account.ServerInfo.HttpsPort
	default:
		err := fmt.Sprintf("Unknown Server Protocol. Expected http/https, Got %s", protocol)
		return "", errors.New(err)
	}

	username := c.username
	password := c.password

	return fmt.Sprintf(queryUrl, protocol, domain, port, stream, username, password, id, ext), nil
}

func (c *XtreamClient) GetAccountInfo() (Account, error) {
	var account Account
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionAccountInfo)

	resp, err := SendRequest(query)
	if err != nil {
		return Account{}, err
	}

	err = json.Unmarshal(resp, &account)
	if err != nil {
		return Account{}, err
	}

	c.account = account
	c.rawAccount = resp
	return account, nil
}

func (c *XtreamClient) GetLiveStreamCategories() ([]Category, error) {
	var categories []Category
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionLiveCategories)

	resp, err := SendRequest(query)
	if err != nil {
		return []Category{}, err
	}

	err = json.Unmarshal(resp, &categories)
	if err != nil {
		return []Category{}, err
	}

	c.liveCategories = categories
	c.rawLiveCategories = resp
	return categories, nil
}

func (c *XtreamClient) GetMovieCategories() ([]Category, error) {
	var categories []Category
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionMovieCategories)

	resp, err := SendRequest(query)
	if err != nil {
		return []Category{}, err
	}

	err = json.Unmarshal(resp, &categories)
	if err != nil {
		return []Category{}, err
	}

	c.movieCategories = categories
	c.rawMovieCategories = resp
	return categories, nil
}

func (c *XtreamClient) GetSeriesCategories() ([]Category, error) {
	var categories []Category
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionSeriesCategories)

	resp, err := SendRequest(query)
	if err != nil {
		return []Category{}, err
	}

	err = json.Unmarshal(resp, &categories)
	if err != nil {
		return []Category{}, err
	}

	c.seriesCategories = categories
	c.rawSeriesCategories = resp
	return categories, nil
}

func (c *XtreamClient) GetLiveStreams() ([]LiveStream, error) {
	var livestreams []LiveStream
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionLivestreams)

	resp, err := SendRequest(query)
	if err != nil {
		return []LiveStream{}, err
	}

	err = json.Unmarshal(resp, &livestreams)
	if err != nil {
		return []LiveStream{}, err
	}

	c.livestreams = livestreams
	c.rawLiveStreams = resp
	return livestreams, nil
}

func (c *XtreamClient) GetMovies() ([]Movie, error) {
	var movies []Movie
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionMovies)

	resp, err := SendRequest(query)
	if err != nil {
		return []Movie{}, err
	}

	err = json.Unmarshal(resp, &movies)
	if err != nil {
		return []Movie{}, err
	}

	c.movies = movies
	c.rawMovies = resp
	return movies, nil
}

func (c *XtreamClient) GetSeries() ([]Series, error) {
	var series []Series
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionSeries)

	resp, err := SendRequest(query)
	if err != nil {
		return []Series{}, err
	}

	err = json.Unmarshal(resp, &series)
	if err != nil {
		return []Series{}, err
	}

	c.series = series
	c.rawSeries = resp
	return series, nil
}

func (c *XtreamClient) GetMovieInfo(id int) (MovieInfo, error) {
	var info MovieInfo
	action := fmt.Sprintf(actionMovieInfo, id)
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, action)

	resp, err := SendRequest(query)
	if err != nil {
		return MovieInfo{}, err
	}

	err = json.Unmarshal(resp, &info)
	if err != nil {
		return MovieInfo{}, err
	}

	return info, nil
}

func (c *XtreamClient) GetSeriesInfo(id int) (SeriesInfo, error) {
	var info SeriesInfo
	action := fmt.Sprintf(actionSeriesInfo, id)
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, action)

	resp, err := SendRequest(query)
	if err != nil {
		return SeriesInfo{}, err
	}

	err = json.Unmarshal(resp, &info)
	if err != nil {
		return SeriesInfo{}, err
	}

	return info, nil
}

func (c *XtreamClient) ExportLiveStreams() error {
	updated_streams := 0
	updated_images := 0

	if len(c.livestreams) == 0 {
		return errors.New("No available LiveStreams for export!")
	}

	log.Printf("[INFO] Exporting LiveStreams...")

	// Create Root Directory
	err := os.MkdirAll(directoryLivestreams, 0o750)
	if err != nil {
		return err
	}

	length := len(c.livestreams)
	for i := range c.livestreams {
		livestream := c.livestreams[i]

		url, err := c.buildURL(livestream.StreamType, livestream.Id, c.account.UserInfo.AllowedOutputFormats[0])
		if err != nil {
			return err
		}

		updated_stream, updated_image, err := livestream.Export(directoryLivestreams, url)
		if err != nil {
			return err
		}
		updated_streams += updated_stream
		updated_streams += updated_image

		// Output Progress Information
		if i%(length/debugPercent) == 0 || i == length-1 {
			percentage := float64(i) / float64(length) * 100
			log.Printf("[DEBUG] (LiveStream) Export Progress: %6d / %6d (%6.2f%%)\n", i+1, length, percentage)
		}
	}

	log.Printf("[INFO] LiveStreams Processed: %d, LiveStreams Updated: %d, Images Skipped: %d\n", length, updated_streams, updated_images)

	return nil
}

func (c *XtreamClient) ExportMovies() error {
	updated_streams := 0
	updated_images := 0

	if len(c.movies) == 0 {
		return errors.New("No available Movies for export!")
	}

	log.Printf("[INFO] Exporting Movies...")

	// Create Root Directory
	err := os.MkdirAll(directoryMovies, 0o750)
	if err != nil {
		return err
	}

	length := len(c.movies)
	for i := range c.movies {
		movie := c.movies[i]

		url, err := c.buildURL(movie.StreamType, movie.Id, movie.Extension)
		if err != nil {
			return err
		}

		updated_stream, updated_image, err := movie.Export(directoryMovies, url)
		if err != nil {
			return err
		}
		updated_streams += updated_stream
		updated_images += updated_image

		// Output Progress Information
		if i%(length/debugPercent) == 0 || i == length-1 {
			percentage := float64(i) / float64(length) * 100
			log.Printf("[DEBUG] (  Movies  ) Export Progress: %6d / %6d (%6.2f%%)\n", i+1, length, percentage)
		}
	}

	log.Printf("[INFO] Movies Processed: %d, Movies Updated: %d, Images Updated: %d\n", length, updated_streams, updated_images)

	return nil
}

func (c *XtreamClient) ExportSeries() error {
	updated_streams := 0
	updated_images := 0

	if len(c.series) == 0 {
		return errors.New("No available Series for export!")
	}

	log.Printf("[INFO] Exporting Series...")

	err := os.MkdirAll(directorySeries, 0o750)
	if err != nil {
		return err
	}

	length := len(c.series)
	for i := range c.series {
		show := c.series[i]

		url, err := c.buildURL("series", show.Id, "")
		if err != nil {
			return err
		}

		updated_stream, updated_image, err := show.Export(directorySeries, url)
		if err != nil {
			return err
		}
		updated_streams += updated_stream
		updated_images += updated_image

		// Output Progress Information
		if i%(length/debugPercent) == 0 || i == length-1 {
			percentage := float64(i) / float64(length) * 100
			log.Printf("[DEBUG] (  Series  ) Export Progress: %6d / %6d (%6.2f%%)\n", i+1, length, percentage)
		}
	}

	log.Printf("[INFO] Series Processed: %d, Series Updated: %d, Images Skipped: %d\n", length, updated_streams, updated_images)

	return nil
}
