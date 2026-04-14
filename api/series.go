package api

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/CaptainKills/xtream-api/utils"
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
	var series []Series
	series_map := map[int]Series{}

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionSeries)

	// Fetch Series
	resp, err := c.sendRequest(query)
	if err != nil {
		return map[int]Series{}, err
	}

	// Unmarshal Series
	err = json.Unmarshal(resp, &series)
	if err != nil {
		return map[int]Series{}, err
	}

	// Map Series
	for _, show := range series {
		series_map[show.Id] = show
	}

	return series_map, nil
}

func (c *XtreamClient) GetSeriesInfo(id int) (SeriesInfo, error) {
	var info SeriesInfo
	action := fmt.Sprintf(actionSeriesInfo, id)
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, action)

	resp, err := c.sendRequest(query)
	if err != nil {
		return SeriesInfo{}, err
	}

	err = json.Unmarshal(resp, &info)
	if err != nil {
		return SeriesInfo{}, err
	}

	return info, nil
}

func (s Series) Export(c *XtreamClient, dir string) (int, int, int, error) {
	updated_streams := 0
	updated_images := 0
	updated_nfos := 0

	s.Name = strings.ReplaceAll(s.Name, "/", "_")

	pathDirectory := dir + s.Name
	pathImage := dir + "/cover" + utils.GetImageExtension(s.Cover)
	pathNfo := dir + "/tvshow.nfo"

	// Create Subdirectory
	err := os.Mkdir(pathDirectory, 0o750)
	if err != nil && !os.IsExist(err) {
		return updated_streams, updated_images, updated_nfos, err
	}

	// Fetch Series Info
	info, err := c.GetSeriesInfo(s.Id)
	if err != nil {
		return updated_streams, updated_images, updated_nfos, err
	}

	// Write Episodes to File
	for season, episodes := range info.Episodes {
		pathSeason := pathDirectory + "/Season " + season

		// Create Season Subdirectory
		err := os.Mkdir(pathSeason, 0o750)
		if err != nil && !os.IsExist(err) {
			return updated_streams, updated_images, updated_nfos, err
		}

		for _, episode := range episodes {
			updated_stream, updated_nfo, err := episode.Export(c, pathSeason)
			if err != nil {
				return updated_streams, updated_images, updated_nfos, err
			}

			updated_streams += updated_stream
			updated_nfos += updated_nfo
		}
	}

	// Write Image to File
	if c.Options.ImagesEnabled && !utils.ImageExists(pathImage) {
		image, err := c.sendRequest(s.Cover)
		if err != nil {
			return updated_streams, updated_images, updated_nfos, err
		}

		updated_image, err := utils.WriteImage(pathImage, image)
		if err != nil {
			return updated_streams, updated_images, updated_nfos, err
		}

		updated_images += updated_image
	}

	// Write NFO to File
	if c.Options.NfoEnabled {
		updated_nfo, err := utils.WriteFile(pathNfo, info.GenerateNfo())
		if err != nil {
			return updated_streams, updated_images, updated_nfos, err
		}

		updated_nfos += updated_nfo
	}

	return updated_streams, updated_images, updated_nfos, nil
}

func (i SeriesInfo) GenerateNfo() string {
	builder := &strings.Builder{}

	builder.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	builder.WriteString("<tvshow>")

	fmt.Fprintf(builder, "<title>%s</title>", i.Info.Name)
	fmt.Fprintf(builder, "<plot>%s</plot>", i.Info.Plot)
	fmt.Fprintf(builder, "<releasedate>%s</releasedate>", i.Info.ReleaseDate)

	genres := strings.SplitSeq(i.Info.Genre, ", ")
	for genre := range genres {
		fmt.Fprintf(builder, "<genre>%s</genre>", genre)
	}

	directors := strings.SplitSeq(i.Info.Director, ", ")
	for director := range directors {
		fmt.Fprintf(builder, "<director>%s</director>", director)
	}

	actors := strings.Split(i.Info.Cast, ", ")
	for index, actor := range actors {
		fmt.Fprintf(builder, "<actor>")
		fmt.Fprintf(builder, "<name>%s</name>", actor)
		fmt.Fprintf(builder, "<order>%d</order>", index)
		fmt.Fprintf(builder, "</actor>")
	}

	builder.WriteString("</tvshow>")

	return builder.String()
}

func (e Episode) Export(c *XtreamClient, dir string) (int, int, error) {
	updated_stream := 0
	updated_nfo := 0

	e.Title = strings.ReplaceAll(e.Title, "/", "_")

	pathStream := dir + "/" + e.Title + ".strm"
	pathNfo := dir + "/" + e.Title + ".nfo"

	id, err := strconv.Atoi(e.Id)
	if err != nil {
		return updated_stream, updated_nfo, err
	}
	url := c.buildURL("series", id, e.Extension)

	// Write Stream to File
	updated_stream, err = utils.WriteFile(pathStream, url)
	if err != nil {
		return updated_stream, updated_nfo, err
	}

	// Write NFO to File
	if c.Options.NfoEnabled {
		updated_nfo, err = utils.WriteFile(pathNfo, e.GenerateNfo())
		if err != nil {
			return updated_stream, updated_nfo, err
		}
	}

	return updated_stream, updated_nfo, nil
}

func (e Episode) GenerateNfo() string {
	builder := &strings.Builder{}
	e.Title = strings.ReplaceAll(e.Title, "&", "&amp;")

	var title string
	sections := strings.Split(e.Title, " - ")
	for index, part := range sections {
		if index == 0 {
			continue
		}
		title += part
		if index != len(sections)-1 {
			title += " "
		}
	}

	re := regexp.MustCompile(`[sS][0-9]*[eE][0-9]*\s`)
	title = re.ReplaceAllString(title, "")

	builder.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	builder.WriteString("<episodedetails>")
	fmt.Fprintf(builder, "<title>%s</title>", title)
	builder.WriteString("</episodedetails>")

	return builder.String()
}
