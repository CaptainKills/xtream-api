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

	Info SeriesInfo
}

type SeriesInfo struct {
	Episodes map[string][]Episode `json:"episodes"`
	Info     ExtraSeriesInfo      `json:"info"`
	Seasons  []Season             `json:"seasons"`
}

type ExtraSeriesInfo struct {
	// BackdropPath   []string `json:"backdrop_path"`
	Cast string `json:"cast"`
	// CategoryId     string `json:"category_id"` // int
	// CategoryIds    []int  `json:"category_ids"`
	Cover    string `json:"cover"`
	Director string `json:"director"`
	// EpisodeRunTime string `json:"episode_run_time"` // int
	Genre string `json:"genre"`
	// LastModified   string `json:"last_modified"` // time.Time
	Name string `json:"name"`
	Plot string `json:"plot"`
	// Rating         string `json:"rating"`        // float64
	// Rating5Based   string `json:"rating_5based"` // float64
	ReleaseDate string `json:"releaseDate"` // time.Time
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

func (c *XtreamClient) GetSeries() (map[int]Series, error) {
	c.series = map[int]Series{}
	var series []Series

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionSeries)

	// Fetch Series
	resp, err := SendRequest(query)
	if err != nil {
		return map[int]Series{}, err
	}

	err = json.Unmarshal(resp, &series)
	if err != nil {
		return map[int]Series{}, err
	}

	// Filter Banned Series
	for _, show := range series {
		allowed := true

		for _, id := range show.CategoryIds {
			if _, ok := c.seriesCategories[id]; !ok {
				allowed = false
			}
		}

		if allowed {
			c.series[show.Id] = show
		}
	}

	return c.series, nil
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

func (s Series) Export(dir string, ur string, enabledImages bool, enabledNfo bool) (int, int, error) {
	s.Name = strings.ReplaceAll(s.Name, "/", "_")

	pathImage := dir + "/cover" + GetImageExtension(s.Cover)
	pathNfo := dir + "/tvshow.nfo"

	// Write Image to File
	updated_image, err := WriteImage(dir, pathImage, s.Cover, enabledImages)
	if err != nil {
		return updated_image, 0, err
	}

	// Write NFO to File
	updated_nfo, err := WriteNfo(dir, pathNfo, GenerateSeriesNfo(s.Info), enabledNfo)
	if err != nil {
		return updated_image, updated_nfo, err
	}

	return updated_image, updated_nfo, nil
}

func (e Episode) Export(dir string, url string, enableNfo bool) (int, int, error) {
	e.Title = strings.ReplaceAll(e.Title, "/", "_")

	pathStream := dir + "/" + e.Title + ".strm"
	pathNfo := dir + "/" + e.Title + ".nfo"

	// Write Stream to File
	updated_stream, err := WriteStream(dir, pathStream, url)
	if err != nil {
		return updated_stream, 0, err
	}

	// Write NFO to File
	updated_nfo, err := WriteNfo(dir, pathNfo, GenerateEpisodeNfo(e), enableNfo)
	if err != nil {
		return updated_stream, updated_nfo, err
	}

	return updated_stream, updated_nfo, nil
}
