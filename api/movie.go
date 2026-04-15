package api

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/CaptainKills/xtream-api/utils"
)

type Movie struct {
	Added        string `json:"added"`       // time.Time
	CategoryId   string `json:"category_id"` // int
	CategoryIds  []int  `json:"category_ids"`
	CystomSID    string `json:"custom_sid"`
	DirectSource string `json:"direct_source"`
	Extension    string `json:"container_extension"`
	Icon         string `json:"stream_icon"`
	Id           int    `json:"stream_id"`
	IsAdult      int    `json:"is_adult"` // bool
	Name         string `json:"name"`
	Number       int    `json:"num"`
	Rating       string `json:"rating"` // float64
	// Rating5Based string `json:"rating_5based"` // float64
	StreamType string `json:"stream_type"`
	// TMDB       string `json:"tmdb"` // int
	Trailer string `json:"trailer"`
}

type MovieInfo struct {
	Info ExtraMovieInfo `json:"info"`
	// Data MovieData      `json:"movie_data"`
}

type ExtraMovieInfo struct {
	Actors       string   `json:"actors"`
	Age          string   `json:"age"` // int
	BackdropPath []string `json:"backdrop_path"`
	Bitrate      int      `json:"bitrate"`
	Cast         string   `json:"cast"`
	Country      string   `json:"country"`
	Cover        string   `json:"cover_big"`
	Description  string   `json:"description"`
	Director     string   `json:"director"`
	Duration     string   `json:"duration"` // time.Time
	// DurationSecs string `json:"duration_secs"` // int
	Genre        string `json:"genre"`
	Image        string `json:"movie_image"`
	KinoUrl      string `json:"kinopoisk_url"`
	Name         string `json:"name"`
	OriginalName string `json:"o_name"`
	Plot         string `json:"plot"`
	// Rating string `json:"rating"` // float64
	ReleaseDate string `json:"releasedate"` // time.Time
	// RunTime string `json:"episode_run_time"` // int
	TMDB    int    `json:"tmdb"`
	Trailer string `json:"youtube_trailer"`
}

func (c *XtreamClient) GetMovies() (map[int]Movie, error) {
	var movies []Movie
	movie_map := map[int]Movie{}

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionMovies)

	// Fetch Movies
	resp, err := c.sendRequest(query)
	if err != nil {
		return map[int]Movie{}, err
	}

	// Unmarshal Movies
	err = json.Unmarshal(resp, &movies)
	if err != nil {
		return map[int]Movie{}, err
	}

	// Map Movies
	for _, movie := range movies {
		movie_map[movie.Id] = movie
	}

	return movie_map, nil
}

func (c *XtreamClient) GetMovieInfo(id int) (MovieInfo, error) {
	var info MovieInfo
	action := fmt.Sprintf(actionMovieInfo, id)
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, action)

	resp, err := c.sendRequest(query)
	if err != nil {
		return MovieInfo{}, err
	}

	err = json.Unmarshal(resp, &info)
	if err != nil {
		return MovieInfo{}, err
	}

	return info, nil
}

func (m Movie) Export(c *XtreamClient, dir string) (int, int, int, error) {
	updated_stream := 0
	updated_image := 0
	updated_nfo := 0

	m.Name = strings.ReplaceAll(m.Name, "/", "_")

	pathDirectory := dir + m.Name
	pathStream := pathDirectory + "/" + m.Name + ".strm"
	pathImage := pathDirectory + "/cover" + utils.GetImageExtension(m.Icon)
	pathNfo := pathDirectory + "/movie.nfo"
	url := c.buildURL(m.StreamType, m.Id, m.Extension)

	// Create Subdirectory
	err := os.Mkdir(pathDirectory, 0o750)
	if err != nil && !os.IsExist(err) {
		return updated_stream, updated_image, updated_nfo, err
	}

	// Write Stream to File
	updated_stream, err = utils.WriteFile(pathStream, url)
	if err != nil {
		return updated_stream, updated_image, updated_nfo, err
	}

	// Write Image to File
	if c.Options.ImagesEnabled && !utils.ImageExists(pathImage) && strings.HasPrefix(m.Icon, "http") {
		image, err := c.sendRequest(m.Icon)
		if err != nil {
			// Ignore error for image fetching
			// log.Printf("[WARNING] Failed to fetch Image: %v\n", err)
		} else {
			updated_image, err = utils.WriteImage(pathImage, image)
			if err != nil {
				return updated_stream, updated_image, updated_nfo, err
			}
		}
	}

	// Write NFO to File
	if c.Options.MetadataEnabled {
		info, err := c.GetMovieInfo(m.Id)
		if err != nil {
			// Ignore error for info fetching
			// log.Printf("[WARNING] Failed to fetch Movie Info: %v\n", err)
		} else {
			updated_nfo, err = utils.WriteFile(pathNfo, info.GenerateNfo())
			if err != nil {
				return updated_stream, updated_image, updated_nfo, err
			}
		}
	}

	return updated_stream, updated_image, updated_nfo, nil
}

func (i MovieInfo) GenerateNfo() string {
	builder := &strings.Builder{}

	builder.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	builder.WriteString("<movie>")

	fmt.Fprintf(builder, "<title>%s</title>", i.Info.Name)
	fmt.Fprintf(builder, "<originaltitle>%s</originaltitle>", i.Info.OriginalName)
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

	builder.WriteString("</movie>")

	return builder.String()
}
