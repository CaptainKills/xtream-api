package main

import (
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/CaptainKills/xtream-api/api"
)

func Filter(c *api.XtreamClient) error {
	// Filter Banned Categories: LiveStreams
	banned, total, err := filterCategories(&c.Data.LiveCategories, c.Options.BannedLiveStreams)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Banned %6d out of %6d livestream Categories\t(%6d Remaining)\n", banned, total, len(c.Data.LiveCategories))

	// Filter Banned Categories: Movies
	banned, total, err = filterCategories(&c.Data.MovieCategories, c.Options.BannedMovies)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Banned %6d out of %6d Movie Categories\t\t(%6d Remaining)\n", banned, total, len(c.Data.MovieCategories))

	// Filter Banned Categories: Series
	banned, total, err = filterCategories(&c.Data.SeriesCategories, c.Options.BannedSeries)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Banned %6d out of %6d Series Categories\t(%6d Remaining)\n", banned, total, len(c.Data.SeriesCategories))

	// Filter Banned & Cached Livestreams
	filtered, total, err := filterLiveStreams(&c.Data.Livestreams, c.Old.Livestreams, c.Data.LiveCategories)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Filtered %6d out of %6d Livestreams\t(%6d Remaining)\n", filtered, total, len(c.Data.Livestreams))

	// Filter Banned & Cached Movies
	filtered, total, err = filterMovies(&c.Data.Movies, c.Old.Movies, c.Data.MovieCategories)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Filtered %6d out of %6d Movies\t\t(%6d Remaining)\n", filtered, total, len(c.Data.Movies))

	// Filter Banned & Cached Series
	filtered, total, err = filterSeries(&c.Data.Series, c.Old.Series, c.Data.SeriesCategories)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Filtered %6d out of %6d Series\t\t(%6d Remaining)\n", filtered, total, len(c.Data.Series))

	return nil
}

func filterCategories(Data *map[int]api.Category, bannedCategories []string) (int, int, error) {
	total := len(*Data)
	banned := 0

	if len(bannedCategories) == 1 && bannedCategories[0] == "" {
		return banned, total, nil
	}

	for _, category := range *Data {
		id, err := strconv.Atoi(category.Id)
		if err != nil {
			return 0, 0, err
		}

		for _, filter := range bannedCategories {
			if strings.Contains(category.Name, filter) {
				delete(*Data, id)
				banned++
			}
		}
	}

	return banned, total, nil
}

func filterLiveStreams(Data *map[int]api.LiveStream, old map[int]api.LiveStream, categories map[int]api.Category) (int, int, error) {
	total := len(*Data)
	filtered := 0

	for id, livestream := range *Data {
		if old_stream, ok := old[id]; ok {
			if reflect.DeepEqual(livestream, old_stream) {
				delete(*Data, id)
				filtered++
				continue
			}
		}

		for _, catId := range livestream.CategoryIds {
			allowed := false

			for _, category := range categories {
				checkId, err := strconv.Atoi(category.Id)
				if err != nil {
					return filtered, total, err
				}

				if catId == checkId {
					allowed = true
					break
				}
			}

			if !allowed {
				delete(*Data, id)
				filtered++
			}
		}
	}

	return filtered, total, nil
}

func filterMovies(Data *map[int]api.Movie, old map[int]api.Movie, categories map[int]api.Category) (int, int, error) {
	total := len(*Data)
	filtered := 0

	for id, movie := range *Data {
		if old_stream, ok := old[id]; ok {
			if reflect.DeepEqual(movie, old_stream) {
				delete(*Data, id)
				filtered++
				continue
			}
		}

		for _, catId := range movie.CategoryIds {
			allowed := false

			for _, category := range categories {
				checkId, err := strconv.Atoi(category.Id)
				if err != nil {
					return filtered, total, err
				}

				if catId == checkId {
					allowed = true
					break
				}
			}

			if !allowed {
				delete(*Data, id)
				filtered++
			}
		}
	}

	return filtered, total, nil
}

func filterSeries(Data *map[int]api.Series, old map[int]api.Series, categories map[int]api.Category) (int, int, error) {
	total := len(*Data)
	filtered := 0

	for id, show := range *Data {
		if old_stream, ok := old[id]; ok {
			if reflect.DeepEqual(show, old_stream) {
				delete(*Data, id)
				filtered++
				continue
			}
		}

		for _, catId := range show.CategoryIds {
			allowed := false

			for _, category := range categories {
				checkId, err := strconv.Atoi(category.Id)
				if err != nil {
					return filtered, total, err
				}

				if catId == checkId {
					allowed = true
					break
				}
			}

			if !allowed {
				delete(*Data, id)
				filtered++
			}
		}
	}

	return filtered, total, nil
}
