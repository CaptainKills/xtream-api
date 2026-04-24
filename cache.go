package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/CaptainKills/xtream-api/api"
)

var (
	directoryCache = "cache"

	cacheFileLivestreams = "livestreams.json"
	cacheFileMovies      = "movies.json"
	cacheFileSeries      = "series.json"

	metadataFileLivestreams = "livestreams_metadata.json"
	metadataFileMovies      = "movies_metadata.json"
	metadataFileSeries      = "series_metadata.json"
)

func ImportCache[T api.Stream](streams *map[int]T, file string, label string) error {
	data, err := os.ReadFile(filepath.Join(directoryCache, file))
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if data == nil {
		data = []byte("{}") // Create Empty JSON Object if file doesn't exist.
	}

	err = json.Unmarshal(data, streams)
	if err != nil {
		return err
	}

	log.Printf("[INFO] (%s) Imported %6d Entries from Cache\n", label, len(*streams))

	return nil
}

func ExportCache[T api.Stream](streams *map[int]T, cache *map[int]T, state *map[int]*State, file string, label string) error {
	// Create Root Directory
	err := os.MkdirAll(directoryCache, MODE_DIR)
	if err != nil {
		return err
	}

	// Copy Streams to Cache
	for id, stream := range *streams {
		s, ok := (*state)[id]
		if !ok {
			// No state, so doesn't need to be exported to cache
			continue
		}

		if s.Strm == true || s.Image == true || s.Nfo == true {
			(*cache)[id] = stream
		}
	}

	// Export Cache
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(directoryCache, file), data, MODE_FILE)
	if err != nil {
		return err
	}

	log.Printf("[INFO] (%s) Exported %6d Entries to Cache\n", label, len(*cache))

	return nil
}

func LoadState(file string, label string) (map[int]*State, error) {
	state := map[int]*State{}

	data, err := os.ReadFile(filepath.Join(directoryCache, file))
	if err != nil && !os.IsNotExist(err) {
		return map[int]*State{}, err
	}

	if data == nil {
		data = []byte("{}") // Create Empty JSON Object if file doesn't exist.
	}

	err = json.Unmarshal(data, &state)
	if err != nil {
		return map[int]*State{}, err
	}

	log.Printf("[INFO] (%s) Loaded State of %6d Entries\n", label, len(state))

	return state, nil
}

func SaveState(state *map[int]*State, file string, label string) error {
	// Create Root Directory
	err := os.MkdirAll(directoryCache, MODE_DIR)
	if err != nil {
		return err
	}

	// Export Metadata
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath.Join(directoryCache, file), data, MODE_FILE)
	if err != nil {
		return err
	}

	log.Printf("[INFO] (%s) Saved State of %6d Entries\n", label, len(*state))

	return nil
}
