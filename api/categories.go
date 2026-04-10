package api

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Category struct {
	Id     string `json:"category_id"` // int
	Name   string `json:"category_name"`
	Parent int    `json:"parent_id"`
}

func (c *XtreamClient) GetLiveStreamCategories() (map[int]Category, error) {
	c.liveCategories = map[int]Category{}
	var categories []Category

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionLiveCategories)

	// Fetch Categories
	resp, err := SendRequest(query)
	if err != nil {
		return map[int]Category{}, err
	}

	err = json.Unmarshal(resp, &categories)
	if err != nil {
		return map[int]Category{}, err
	}

	// Filter Banned Categories
	for _, category := range categories {
		allowed := true

		for _, filter := range c.bannedLivestreams {
			if strings.Contains(category.Name, filter) {
				allowed = false
			}
		}

		if allowed {
			id, err := strconv.Atoi(category.Id)
			if err != nil {
				return map[int]Category{}, err
			}
			c.liveCategories[id] = category
		}
	}

	return c.liveCategories, nil
}

func (c *XtreamClient) GetMovieCategories() (map[int]Category, error) {
	c.movieCategories = map[int]Category{}
	var categories []Category

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionMovieCategories)

	// Fetch Categories
	resp, err := SendRequest(query)
	if err != nil {
		return map[int]Category{}, err
	}

	err = json.Unmarshal(resp, &categories)
	if err != nil {
		return map[int]Category{}, err
	}

	// Filter Banned Categories
	for _, category := range categories {
		allowed := true

		for _, filter := range c.bannedMovies {
			if strings.Contains(category.Name, filter) {
				allowed = false
			}
		}

		if allowed {
			id, err := strconv.Atoi(category.Id)
			if err != nil {
				return map[int]Category{}, err
			}
			c.movieCategories[id] = category
		}
	}

	return c.movieCategories, nil
}

func (c *XtreamClient) GetSeriesCategories() (map[int]Category, error) {
	c.seriesCategories = map[int]Category{}
	var categories []Category

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionSeriesCategories)

	// Fetch Categories
	resp, err := SendRequest(query)
	if err != nil {
		return map[int]Category{}, err
	}

	err = json.Unmarshal(resp, &categories)
	if err != nil {
		return map[int]Category{}, err
	}

	// Filter Banned Categories
	for _, category := range categories {
		allowed := true

		for _, filter := range c.bannedSeries {
			if strings.Contains(category.Name, filter) {
				allowed = false
			}
		}

		if allowed {
			id, err := strconv.Atoi(category.Id)
			if err != nil {
				return map[int]Category{}, err
			}
			c.seriesCategories[id] = category
		}
	}

	return c.seriesCategories, nil
}
