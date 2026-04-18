package main

import (
	"log"
	"reflect"
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

	log.Printf("[INFO] Banned %6d out of %6d Series Categories\t(%6d Remaining)\n", banned, total, len(*series))

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

func FilterStreams[T api.Stream](client *api.XtreamClient, streams *map[int]T, categories map[int]api.Category, cache map[int]T, metadata *map[int]*api.XtreamMetadata, label string) error {
	total := len(*streams)
	filtered := 0

	for id, stream := range *streams {
		// If stream belongs to banned category, do not export
		banned, err := isBanned(stream, categories)
		if err != nil {
			return err
		}

		if banned == true {
			delete(*streams, id)
			filtered++
			continue
		}

		// If stream needs to be updated based on fetched data
		updated := isUpdated(stream, id, cache)
		if updated == true {
			continue
		} else {
			// If stream did not change, check if it needs image or metadata update
			md, ok := (*metadata)[id]
			if !ok {
				// Stream does not have metadata, so must be updated
				continue
			}
			needImageUpdate := client.Options.ImagesEnabled && !md.Image
			needMetadataUpdate := client.Options.MetadataEnabled && !md.Nfo

			// If stream is not updated, and doesn't need image or nfo, do not export
			if !needImageUpdate && !needMetadataUpdate {
				delete(*streams, id)
				filtered++
				continue
			}
		}
	}

	log.Printf("[INFO] Filtered Out %6d out of %6d %s (%6d Remaining)\n", filtered, total, label, len(*streams))

	return nil
}

func isBanned[T api.Stream](stream T, categories map[int]api.Category) (bool, error) {
	for _, catId := range stream.GetCategoryIds() {
		for _, category := range categories {
			checkId, err := strconv.Atoi(category.Id)
			if err != nil {
				return false, err
			}

			if catId == checkId {
				return false, nil
			}
		}
	}

	return true, nil
}

func isUpdated[T api.Stream](stream T, id int, cache map[int]T) bool {
	if old_stream, ok := cache[id]; ok {
		if reflect.DeepEqual(stream, old_stream) {
			return false
		} else {
			return true
		}
	}

	return true // Didn't exist in cache, so stream is new
}
