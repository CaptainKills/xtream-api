package main

import (
	"log"
	"strconv"
	"strings"

	"github.com/CaptainKills/xtream-api/api"
)

func FilterCategories(data *map[int]api.Category, bannedCategories []string, label string) error {
	total := len(*data)
	banned := 0

	if len(bannedCategories) == 1 && bannedCategories[0] == "" {
		return nil
	}

	for _, category := range *data {
		id, err := strconv.Atoi(category.Id)
		if err != nil {
			return err
		}

		for _, filter := range bannedCategories {
			if strings.Contains(category.Name, filter) {
				delete(*data, id)
				banned++
			}
		}
	}

	log.Printf("[INFO] (%s) Banned %6d out of %6d Categories\t(%6d Remaining)\n", label, banned, total, len(*data))

	return nil
}

func FilterStreams[T api.Stream](client *api.XtreamClient, streams *map[int]T, categories map[int]api.Category, cache map[int]T, state *map[int]*State, label string) error {
	count_updated := 0
	count_missing := 0
	count_image := 0
	count_nfo := 0

	count_banned := 0
	count_none := 0

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
			count_banned++
			continue
		}

		// If stream needs to be updated based on fetched data
		updated := isUpdated(stream, id, cache)
		if updated == true {
			count_updated++
			continue
		} else {
			// If stream did not change, check if it needs image or state update
			s, ok := (*state)[id]
			if !ok {
				// Stream does not have state, so must be updated
				count_missing++
				continue
			}
			needImageUpdate := client.Options.ImagesEnabled && !s.Image && strings.HasPrefix(stream.GetCover(), "http")
			needMetadataUpdate := client.Options.MetadataEnabled && !s.Nfo

			if needImageUpdate {
				count_image++
			}

			if needMetadataUpdate {
				count_nfo++
			}

			// If stream is not updated, and doesn't need image or nfo, do not export
			if !needImageUpdate && !needMetadataUpdate {
				delete(*streams, id)
				filtered++
				count_none++
				continue
			}
		}
	}

	log.Printf("[INFO] (%s) Filtered Out %6d out of %6d (%6d Remaining)\n", label, filtered, total, len(*streams))
	log.Printf("[DEBUG] (%s) UPD=%6d, MIS=%6d, IMG=%6d, NFO=%6d | BAN=%6d, NONE=%6d\n", label, count_updated, count_missing, count_image, count_nfo, count_banned, count_none)

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
		return !stream.Equals(old_stream)
	}

	return true // Didn't exist in cache, so stream is new
}
