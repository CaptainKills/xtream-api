package api

import (
	"encoding/json"
	"fmt"
	"strings"
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

type MovieInfo struct{}

func (c *XtreamClient) GetMovies() (map[int]Movie, error) {
	c.movies = map[int]Movie{}
	var movies []Movie

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionMovies)

	// Fetch Movies
	resp, err := SendRequest(query)
	if err != nil {
		return map[int]Movie{}, err
	}

	err = json.Unmarshal(resp, &movies)
	if err != nil {
		return map[int]Movie{}, err
	}

	// Filter Banned Movies
	for _, movie := range movies {
		allowed := true

		for _, id := range movie.CategoryIds {
			if _, ok := c.movieCategories[id]; !ok {
				allowed = false
			}
		}

		if allowed {
			c.movies[movie.Id] = movie
		}
	}

	return c.movies, nil
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

func (m Movie) Export(dir string, url string, enableImages bool) (int, int, error) {
	m.Name = strings.ReplaceAll(m.Name, "/", "_")

	pathDirectory := dir + m.Name
	pathFile := pathDirectory + "/" + m.Name + ".strm"
	pathImage := pathDirectory + "/cover" + GetImageExtension(m.Icon)

	// Write Stream to File
	updated_stream, err := WriteStream(pathDirectory, pathFile, url)
	if err != nil {
		return updated_stream, 0, err
	}

	// Write Image to File
	updated_image, err := WriteImage(pathDirectory, pathImage, m.Icon, enableImages)
	if err != nil {
		return updated_stream, updated_image, err
	}

	return updated_stream, updated_image, nil
}
