package api

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/CaptainKills/xtream-api/utils"
)

type Category struct {
	Id     string `json:"category_id"` // int
	Name   string `json:"category_name"`
	Parent int    `json:"parent_id"`
}

func (c *XtreamClient) GetLiveStreamCategories() (map[int]Category, error) {
	c.Data.LiveCategories = map[int]Category{}
	var categories []Category

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionLiveCategories)

	// Fetch Categories
	resp, err := utils.SendRequest(query)
	if err != nil {
		return map[int]Category{}, err
	}

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

		c.Data.LiveCategories[id] = category
	}

	return c.Data.LiveCategories, nil
}

func (c *XtreamClient) GetMovieCategories() (map[int]Category, error) {
	c.Data.MovieCategories = map[int]Category{}
	var categories []Category

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionMovieCategories)

	// Fetch Categories
	resp, err := utils.SendRequest(query)
	if err != nil {
		return map[int]Category{}, err
	}

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

		c.Data.MovieCategories[id] = category
	}

	return c.Data.MovieCategories, nil
}

func (c *XtreamClient) GetSeriesCategories() (map[int]Category, error) {
	c.Data.SeriesCategories = map[int]Category{}
	var categories []Category

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionSeriesCategories)

	// Fetch Categories
	resp, err := utils.SendRequest(query)
	if err != nil {
		return map[int]Category{}, err
	}

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

		c.Data.SeriesCategories[id] = category
	}

	return c.Data.SeriesCategories, nil
}
