package api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Category struct {
	Id     string `json:"category_id"` // int
	Name   string `json:"category_name"`
	Parent int    `json:"parent_id"`
}

func (c *XtreamClient) GetLiveStreamCategories() (map[int]Category, error) {
	c.data.liveCategories = map[int]Category{}
	var categories []Category

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionLiveCategories)

	// Fetch Categories
	resp, err := SendRequest(query)
	if err != nil {
		return map[int]Category{}, err
	}
	c.raw.liveCategories = resp

	// Unmarshal Categories
	err = json.Unmarshal(resp, &categories)
	if err != nil {
		return map[int]Category{}, err
	}

	// Map Categories
	for _, category := range categories {
		id, err := strconv.Atoi(category.Id)
		if err != nil {
			return map[int]Category{}, err
		}

		c.data.liveCategories[id] = category
	}

	return c.data.liveCategories, nil
}

func (c *XtreamClient) GetMovieCategories() (map[int]Category, error) {
	c.data.movieCategories = map[int]Category{}
	var categories []Category

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionMovieCategories)

	// Fetch Categories
	resp, err := SendRequest(query)
	if err != nil {
		return map[int]Category{}, err
	}
	c.raw.movieCategories = resp

	// Unmarshal Categories
	err = json.Unmarshal(resp, &categories)
	if err != nil {
		return map[int]Category{}, err
	}

	// Map Categories
	for _, category := range categories {
		id, err := strconv.Atoi(category.Id)
		if err != nil {
			return map[int]Category{}, err
		}

		c.data.movieCategories[id] = category
	}

	return c.data.movieCategories, nil
}

func (c *XtreamClient) GetSeriesCategories() (map[int]Category, error) {
	c.data.seriesCategories = map[int]Category{}
	var categories []Category

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionSeriesCategories)

	// Fetch Categories
	resp, err := SendRequest(query)
	if err != nil {
		return map[int]Category{}, err
	}
	c.raw.seriesCategories = resp

	// Unmarshal Categories
	err = json.Unmarshal(resp, &categories)
	if err != nil {
		return map[int]Category{}, err
	}

	// Map Categories
	for _, category := range categories {
		id, err := strconv.Atoi(category.Id)
		if err != nil {
			return map[int]Category{}, err
		}

		c.data.seriesCategories[id] = category
	}

	return c.data.seriesCategories, nil
}
