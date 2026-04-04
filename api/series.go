package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Series struct {
	// BackdropPath   []string `json:"backdrop_path"`
	Cast           string `json:"cast"`
	CategoryId     string `json:"category_id"` // int
	CategoryIds    []int  `json:"category_ids"`
	Cover          string `json:"cover"`
	Director       string `json:"director"`
	EpisodeRunTime string `json:"episode_run_time"` // int
	Genre          string `json:"genre"`
	Id             int    `json:"series_id"`
	LastModified   string `json:"last_modified"` // time.Time
	Name           string `json:"name"`
	Number         int    `json:"num"`
	Plot           string `json:"plot"`
	Rating         string `json:"rating"`        // float64
	Rating5Based   string `json:"rating_5based"` // float64
	ReleaseDate    string `json:"releaseDate"`   // time.Time
	ReleaseDate2   string `json:"release_date"`  // time.Time
	TMDB           string `json:"tmdb"`          // int
	Trailer        string `json:"youtube_trailer"`
}

type SeriesInfo struct {
	Episodes map[string][]Episode `json:"episodes"`
	Info     ExtraInfo            `json:"info"`
	Seasons  []Season             `json:"seasons"`
}

type ExtraInfo struct {
	// BackdropPath   []string `json:"backdrop_path"`
	// Cast           string `json:"cast"`
	// CategoryId     string `json:"category_id"` // int
	// CategoryIds    []int  `json:"category_ids"`
	Cover string `json:"cover"`
	// Director       string `json:"director"`
	// EpisodeRunTime string `json:"episode_run_time"` // int
	// Genre          string `json:"genre"`
	// LastModified   string `json:"last_modified"` // time.Time
	Name string `json:"name"`
	// Plot           string `json:"plot"`
	// Rating         string `json:"rating"`        // float64
	// Rating5Based   string `json:"rating_5based"` // float64
	// ReleaseDate    string `json:"releaseDate"`   // time.Time
	// ReleaseDate2   string `json:"release_date"`  // time.Time
	// TMDB           string `json:"tmdb"`
	// Trailer        string `json:"youtube_trailer"`
}

type Season struct {
	// AirDate      string `json:"air_date"` // time.Time
	// Cover        string `json:"cover"`
	// CoverBig     string `json:"cover_big"`
	// CoverTMDB    string `json:"cover_tmdb"`
	// Duration     string `json:"duration"`      // int
	// EpisodeCount string `json:"episode_count"` // int
	Name   string `json:"name"`
	Number int    `json:"season_number"`
	// Overview     string `json:"overview"`
	// ReleaseDate  string `json:"releaseDate"` // time.Time
}

type Episode struct {
	// Added        string `json:"added"` // time.Time
	// CustomSID    string `json:"custom_sid"`
	// DirectSource string `json:"direct_source"`
	Extension string `json:"container_extension"`
	Id        string `json:"id"` // int
	// Info         EpisodeInfo `json:"info"`
	// Number int    `json:"episode_num"`
	// Season int    `json:"season"` // int
	Title string `json:"title"`
}

type EpisodeInfo struct{}

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

func (s Series) Export(dir string, ur string, enabledImages bool) (int, error) {
	s.Name = strings.ReplaceAll(s.Name, "/", "_")

	pathImage := dir + "/cover" + GetImageExtension(s.Cover)

	// Write Image to File
	updated_image, err := WriteImage(dir, pathImage, s.Cover, enabledImages)
	if err != nil {
		return updated_image, err
	}

	return updated_image, nil
}

func (e Episode) Export(dir string, url string) (int, error) {
	e.Title = strings.ReplaceAll(e.Title, "/", "_")

	pathFile := dir + "/" + e.Title + ".strm"

	// Write Stream to File
	updated_stream, err := WriteStream(dir, pathFile, url)
	if err != nil {
		return 0, err
	}

	return updated_stream, nil
}
