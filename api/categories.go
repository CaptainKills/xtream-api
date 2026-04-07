package api

import (
	"encoding/json"
	"fmt"
)

type Category struct {
	Id     string `json:"category_id"` // int
	Name   string `json:"category_name"`
	Parent int    `json:"parent_id"`
}

func (c *XtreamClient) GetLiveStreamCategories() ([]Category, error) {
	var categories []Category
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionLiveCategories)

	resp, err := SendRequest(query)
	if err != nil {
		return []Category{}, err
	}

	err = json.Unmarshal(resp, &categories)
	if err != nil {
		return []Category{}, err
	}

	c.liveCategories = categories
	return categories, nil
}

func (c *XtreamClient) GetMovieCategories() ([]Category, error) {
	var categories []Category
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionMovieCategories)

	resp, err := SendRequest(query)
	if err != nil {
		return []Category{}, err
	}

	err = json.Unmarshal(resp, &categories)
	if err != nil {
		return []Category{}, err
	}

	c.movieCategories = categories
	return categories, nil
}

func (c *XtreamClient) GetSeriesCategories() ([]Category, error) {
	var categories []Category
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionSeriesCategories)

	resp, err := SendRequest(query)
	if err != nil {
		return []Category{}, err
	}

	err = json.Unmarshal(resp, &categories)
	if err != nil {
		return []Category{}, err
	}

	c.seriesCategories = categories
	return categories, nil
}
