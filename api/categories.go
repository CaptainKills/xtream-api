package api

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Category struct {
	Id     string `json:"category_id"`
	Name   string `json:"category_name"`
	Parent int    `json:"parent_id"`
}

func (c *XtreamClient) GetLiveStreamCategories() (map[int]Category, error) {
	var categories []Category
	category_map := map[int]Category{}

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionLiveCategories)

	// Fetch Categories
	resp, err := c.sendRequest(query)
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

		category_map[id] = category
	}

	return category_map, nil
}

func (c *XtreamClient) GetMovieCategories() (map[int]Category, error) {
	var categories []Category
	category_map := map[int]Category{}

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionMovieCategories)

	// Fetch Categories
	resp, err := c.sendRequest(query)
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

		category_map[id] = category
	}

	return category_map, nil
}

func (c *XtreamClient) GetSeriesCategories() (map[int]Category, error) {
	var categories []Category
	category_map := map[int]Category{}

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionSeriesCategories)

	// Fetch Categories
	resp, err := c.sendRequest(query)
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

		category_map[id] = category
	}

	return category_map, nil
}
