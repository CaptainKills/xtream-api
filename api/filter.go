package api

import (
	"log"
	"reflect"
	"strconv"
	"strings"
)

func (c *XtreamClient) Filter() error {
	// Filter Banned Categories: LiveStreams
	banned, total, err := filterCategories(&c.data.liveCategories, c.options.BannedLiveStreams)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Banned %6d out of %6d Livestream Categories\t(%6d Remaining)\n", banned, total, len(c.data.liveCategories))

	// Filter Banned Categories: Movies
	banned, total, err = filterCategories(&c.data.movieCategories, c.options.BannedMovies)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Banned %6d out of %6d Movie Categories\t\t(%6d Remaining)\n", banned, total, len(c.data.movieCategories))

	// Filter Banned Categories: Series
	banned, total, err = filterCategories(&c.data.seriesCategories, c.options.BannedSeries)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Banned %6d out of %6d Series Categories\t(%6d Remaining)\n", banned, total, len(c.data.seriesCategories))

	// Filter Banned & Cached Livestreams
	filtered, total, err := filterLiveStreams(&c.data.livestreams, c.old.livestreams, c.data.liveCategories)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Filtered %6d out of %6d Livestreams\t(%6d Remaining)\n", filtered, total, len(c.data.livestreams))

	// Filter Banned & Cached Movies
	filtered, total, err = filterMovies(&c.data.movies, c.old.movies, c.data.movieCategories)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Filtered %6d out of %6d Movies\t\t(%6d Remaining)\n", filtered, total, len(c.data.movies))

	// Filter Banned & Cached Series
	filtered, total, err = filterSeries(&c.data.series, c.old.series, c.data.seriesCategories)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Filtered %6d out of %6d Series\t\t(%6d Remaining)\n", filtered, total, len(c.data.series))

	return nil
}

func filterCategories(data *map[int]Category, bannedCategories []string) (int, int, error) {
	total := len(*data)
	banned := 0

	if len(bannedCategories) == 1 && bannedCategories[0] == "" {
		return banned, total, nil
	}

	for _, category := range *data {
		id, err := strconv.Atoi(category.Id)
		if err != nil {
			return 0, 0, err
		}

		for _, filter := range bannedCategories {
			if strings.Contains(category.Name, filter) {
				delete(*data, id)
				banned++
			}
		}
	}

	return banned, total, nil
}

func filterLiveStreams(data *map[int]LiveStream, old map[int]LiveStream, categories map[int]Category) (int, int, error) {
	total := len(*data)
	filtered := 0

	for id, livestream := range *data {
		if old_stream, ok := old[id]; ok {
			if reflect.DeepEqual(livestream, old_stream) {
				delete(*data, id)
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
				delete(*data, id)
				filtered++
			}
		}
	}

	return filtered, total, nil
}

func filterMovies(data *map[int]Movie, old map[int]Movie, categories map[int]Category) (int, int, error) {
	total := len(*data)
	filtered := 0

	for id, movie := range *data {
		if old_stream, ok := old[id]; ok {
			if reflect.DeepEqual(movie, old_stream) {
				delete(*data, id)
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
				delete(*data, id)
				filtered++
			}
		}
	}

	return filtered, total, nil
}

func filterSeries(data *map[int]Series, old map[int]Series, categories map[int]Category) (int, int, error) {
	total := len(*data)
	filtered := 0

	for id, show := range *data {
		if old_stream, ok := old[id]; ok {
			if reflect.DeepEqual(show, old_stream) {
				delete(*data, id)
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
				delete(*data, id)
				filtered++
			}
		}
	}

	return filtered, total, nil
}
