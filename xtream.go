package main

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"io"
// 	"log"
// 	"net/http"
// 	"strconv"
// )
//
//
//
// type (
// 	MovieInfo  struct{}
// 	SeriesInfo struct{}
// )
//
//
// func GetMovieCategories() []Category {
// 	var data []Category
//
// 	if IsFileOld(PATH_MOVIE_CAT) {
// 		action := "get_vod_categories"
//
// 		query := fmt.Sprintf(QUERY, URL, USERNAME, PASSWORD, action)
// 		resp := sendRequest(query)
//
// 		format := formatJSON(resp)
// 		WriteFile(PATH_MOVIE_CAT, format)
// 		// fmt.Println(format)
//
// 		err := json.Unmarshal(resp, &data)
// 		if err != nil {
// 			log.Fatal("Error parsing JSON:", err)
// 		}
// 	} else {
// 		fmt.Println("Movie Categories: reading from cache file...")
// 		resp := ReadFile(PATH_MOVIE_CAT)
//
// 		err := json.Unmarshal([]byte(resp), &data)
// 		if err != nil {
// 			log.Fatal("Error parsing JSON:", err)
// 		}
// 	}
//
// 	return data
// }
//
// func GetSerieCategories() []Category {
// 	var data []Category
//
// 	if IsFileOld(PATH_SERIE_CAT) {
// 		action := "get_series_categories"
//
// 		query := fmt.Sprintf(QUERY, URL, USERNAME, PASSWORD, action)
// 		resp := sendRequest(query)
//
// 		format := formatJSON(resp)
// 		WriteFile(PATH_SERIE_CAT, format)
// 		// fmt.Println(format)
//
// 		err := json.Unmarshal(resp, &data)
// 		if err != nil {
// 			log.Fatal("Error parsing JSON:", err)
// 		}
// 	} else {
// 		fmt.Println("Serie Categories: reading from cache file...")
// 		resp := ReadFile(PATH_SERIE_CAT)
//
// 		err := json.Unmarshal([]byte(resp), &data)
// 		if err != nil {
// 			log.Fatal("Error parsing JSON:", err)
// 		}
// 	}
//
// 	return data
// }
//
// func GetLiveStreams() []LiveStream {
// 	var data []LiveStream
//
// 	if IsFileOld(PATH_LIVESTREAMS) {
// 		action := "get_live_streams"
//
// 		query := fmt.Sprintf(QUERY, URL, USERNAME, PASSWORD, action)
// 		resp := sendRequest(query)
//
// 		format := formatJSON(resp)
// 		WriteFile(PATH_LIVESTREAMS, format)
// 		// fmt.Println(format)
//
// 		err := json.Unmarshal(resp, &data)
// 		if err != nil {
// 			log.Fatal("Error parsing JSON:", err)
// 		}
// 	} else {
// 		fmt.Println("Livestreams: reading from cache file...")
// 		resp := ReadFile(PATH_LIVESTREAMS)
//
// 		err := json.Unmarshal([]byte(resp), &data)
// 		if err != nil {
// 			log.Fatal("Error parsing JSON:", err)
// 		}
// 	}
//
// 	return data
// }
//
// func GetMovies() []Movie {
// 	var data []Movie
//
// 	if IsFileOld(PATH_MOVIES) {
// 		action := "get_vod_streams"
//
// 		query := fmt.Sprintf(QUERY, URL, USERNAME, PASSWORD, action)
// 		resp := sendRequest(query)
//
// 		format := formatJSON(resp)
// 		WriteFile(PATH_MOVIES, format)
// 		// fmt.Println(format)
//
// 		err := json.Unmarshal(resp, &data)
// 		if err != nil {
// 			log.Fatal("Error parsing JSON:", err)
// 		}
// 	} else {
// 		fmt.Println("Movies: reading from cache file...")
// 		resp := ReadFile(PATH_MOVIES)
//
// 		err := json.Unmarshal([]byte(resp), &data)
// 		if err != nil {
// 			log.Fatal("Error parsing JSON:", err)
// 		}
// 	}
//
// 	return data
// }
//
// func GetSeries() []Serie {
// 	var data []Serie
//
// 	if IsFileOld(PATH_SERIES) {
// 		action := "get_series"
//
// 		query := fmt.Sprintf(QUERY, URL, USERNAME, PASSWORD, action)
// 		resp := sendRequest(query)
//
// 		format := formatJSON(resp)
// 		WriteFile(PATH_SERIES, format)
// 		// fmt.Println(format)
//
// 		err := json.Unmarshal(resp, &data)
// 		if err != nil {
// 			log.Fatal("Error parsing JSON:", err)
// 		}
// 	} else {
// 		fmt.Println("Series: reading from cache file...")
// 		resp := ReadFile(PATH_SERIES)
//
// 		err := json.Unmarshal([]byte(resp), &data)
// 		if err != nil {
// 			log.Fatal("Error parsing JSON:", err)
// 		}
// 	}
//
// 	return data
// }
//
// func GetMovieInfo(id int) (MovieInfo, string) {
// 	action := "get_vod_info&vod_id=" + strconv.Itoa(id)
//
// 	query := fmt.Sprintf(QUERY, URL, USERNAME, PASSWORD, action)
// 	resp := sendRequest(query)
// 	format := formatJSON(resp)
//
// 	var data MovieInfo
// 	err := json.Unmarshal(resp, &data)
// 	if err != nil {
// 		log.Fatal("Error parsing JSON:", err)
// 	}
//
// 	return data, format
// }
//
// func GetSeriesInfo(id int) (SeriesInfo, string) {
// 	action := "get_series_info&series_id=" + strconv.Itoa(id)
//
// 	query := fmt.Sprintf(QUERY, URL, USERNAME, PASSWORD, action)
// 	resp := sendRequest(query)
// 	format := formatJSON(resp)
//
// 	var data SeriesInfo
// 	err := json.Unmarshal(resp, &data)
// 	if err != nil {
// 		log.Fatal("Error parsing JSON:", err)
// 	}
//
// 	return data, format
// }
