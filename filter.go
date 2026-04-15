package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/CaptainKills/xtream-api/api"
)

func FilterCategories(options *api.XtreamOptions, live *map[int]api.Category, movie *map[int]api.Category, series *map[int]api.Category) error {
	// Filter Banned Categories: LiveStreams
	banned, total, err := filterCategory(live, options.BannedLiveStreams)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Banned %6d out of %6d livestream Categories\t(%6d Remaining)\n", banned, total, len(*live))

	// Filter Banned Categories: Movies
	banned, total, err = filterCategory(movie, options.BannedMovies)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Banned %6d out of %6d Movie Categories\t\t(%6d Remaining)\n", banned, total, len(*movie))

	// Filter Banned Categories: Series
	banned, total, err = filterCategory(series, options.BannedSeries)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Banned %6d out of %6d Series Categories\t\t(%6d Remaining)\n", banned, total, len(*series))

	return nil
}

func filterCategory(data *map[int]api.Category, bannedCategories []string) (int, int, error) {
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

func FilterLiveStreams(client *api.XtreamClient, livestreams *map[int]api.LiveStream, categories map[int]api.Category) error {
	total := len(*livestreams)
	filtered := 0

	for id, livestream := range *livestreams {
		// if old_stream, ok := old[id]; ok {
		// 	if reflect.DeepEqual(livestream, old_stream) {
		// 		delete(*livestreams, id)
		// 		filtered++
		// 		continue
		// 	}
		// }

		for _, catId := range livestream.CategoryIds {
			allowed := false

			for _, category := range categories {
				checkId, err := strconv.Atoi(category.Id)
				if err != nil {
					return err
				}

				if catId == checkId {
					allowed = true
					break
				}
			}

			if !allowed {
				delete(*livestreams, id)
				filtered++
			}
		}
	}

	log.Printf("[INFO] Filtered %6d out of %6d Livestreams\t(%6d Remaining)\n", filtered, total, len(*livestreams))

	return nil
}

func FilterMovies(client *api.XtreamClient, movies *map[int]api.Movie, categories map[int]api.Category) error {
	total := len(*movies)
	filtered := 0

	for id, movie := range *movies {
		// if old_stream, ok := old[id]; ok {
		// 	if reflect.DeepEqual(movie, old_stream) {
		// 		delete(*Data, id)
		// 		filtered++
		// 		continue
		// 	}
		// }

		for _, catId := range movie.CategoryIds {
			allowed := false

			for _, category := range categories {
				checkId, err := strconv.Atoi(category.Id)
				if err != nil {
					return err
				}

				if catId == checkId {
					allowed = true
					break
				}
			}

			if !allowed {
				delete(*movies, id)
				filtered++
			}
		}
	}

	log.Printf("[INFO] Filtered %6d out of %6d Movies\t\t(%6d Remaining)\n", filtered, total, len(*movies))

	return nil
}

func FilterSeries(client *api.XtreamClient, series *map[int]api.Series, categories map[int]api.Category) error {
	total := len(*series)
	filtered := 0

	for id, show := range *series {
		// if old_stream, ok := old[id]; ok {
		// 	if reflect.DeepEqual(show, old_stream) {
		// 		delete(*series, id)
		// 		filtered++
		// 		continue
		// 	}
		// }

		for _, catId := range show.CategoryIds {
			allowed := false

			for _, category := range categories {
				checkId, err := strconv.Atoi(category.Id)
				if err != nil {
					return err
				}

				if catId == checkId {
					allowed = true
					break
				}
			}

			if !allowed {
				delete(*series, id)
				filtered++
			}
		}
	}

	log.Printf("[INFO] Filtered %6d out of %6d Series\t\t(%6d Remaining)\n", filtered, total, len(*series))

	return nil
}
