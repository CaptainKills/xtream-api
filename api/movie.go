package api

import (
	"encoding/json"
	"fmt"
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

	Info MovieInfo
}

type MovieInfo struct {
	Info ExtraMovieInfo `json:"info"`
	// Data MovieData      `json:"movie_data"`
}

type ExtraMovieInfo struct {
	// BackdropPath []string `json:"backdrop_path"`
	Cast string `json:"cast"`
	// Cover        string   `json:"cover_big"`
	Director string `json:"director"`
	// Duration     string   `json:"duration"` // time.Time
	// DurationSecs int      `json:"duration_secs"`
	Genre        string `json:"genre"`
	Name         string `json:"name"`
	OriginalName string `json:"o_name"`
	Plot         string `json:"plot"`
	// Rating       string   `json:"rating"`      // float64
	ReleaseDate string `json:"releasedate"` // time.Time
	// TMDB         int      `json:"tmdb"`
	// Trailer      string   `json:"youtube_trailer"`
}

func (c *XtreamClient) GetMovies() (map[int]Movie, error) {
	c.Data.Movies = map[int]Movie{}
	var movies []Movie

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionMovies)

	// Fetch Movies
	resp, err := utils.SendRequest(query)
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
		c.Data.Movies[movie.Id] = movie
	}

	return c.Data.Movies, nil
}

func (c *XtreamClient) GetMovieInfo(id int) (MovieInfo, error) {
	var info MovieInfo
	action := fmt.Sprintf(actionMovieInfo, id)
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, action)

	resp, err := utils.SendRequest(query)
	if err != nil {
		return MovieInfo{}, err
	}

	err = json.Unmarshal(resp, &info)
	if err != nil {
		return MovieInfo{}, err
	}

	return info, nil
}

func (m Movie) Export(dir string, url string, enableImages bool, enableNfo bool) (int, int, int, error) {
	m.Name = strings.ReplaceAll(m.Name, "/", "_")

	pathDirectory := dir + m.Name
	pathStream := pathDirectory + "/" + m.Name + ".strm"
	pathImage := pathDirectory + "/cover" + utils.GetImageExtension(m.Icon)
	pathNfo := pathDirectory + "/movie.nfo"

	// Write Stream to File
	updated_stream, err := utils.WriteStream(pathDirectory, pathStream, url)
	if err != nil {
		return updated_stream, 0, 0, err
	}

	// Write Image to File
	updated_image, err := utils.WriteImage(pathDirectory, pathImage, m.Icon, enableImages)
	if err != nil {
		return updated_stream, updated_image, 0, err
	}

	// Write NFO to File
	updated_nfo, err := utils.WriteNfo(pathDirectory, pathNfo, m.Info.GenerateNfo(), enableNfo)
	if err != nil {
		return updated_stream, updated_image, updated_nfo, err
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
