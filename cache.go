package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/CaptainKills/xtream-api/api"
)

const (
	directoryCache = "cache/"

	cacheFileLivestreams = "livestreams.json"
	cacheFileMovies      = "movies.json"
	cacheFileSeries      = "series.json"

	metadataFileLivestreams = "livestreams_metadata.json"
	metadataFileMovies      = "movies_metadata.json"
	metadataFileSeries      = "series_metadata.json"
)

func ImportCache[T api.Stream](streams *map[int]T, file string, label string) error {
	data, err := os.ReadFile(directoryCache + file)
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

	log.Printf("[INFO] Imported %6d Entries from %s Cache\n", len(*streams), label)

	return nil
}

func ExportCache[T api.Stream](streams *map[int]T, cache *map[int]T, metadata *map[int]*api.XtreamMetadata, file string, label string) error {
	// Create Root Directory
	err := os.MkdirAll(directoryCache, MODE_DIR)
	if err != nil {
		return err
	}

	// Copy Streams to Cache
	for id, stream := range *streams {
		md, ok := (*metadata)[id]
		if !ok {
			// No metadata, so doesn't need to be exported to cache
			continue
		}

		if md.Strm == true || md.Image == true || md.Nfo == true {
			(*cache)[id] = stream
		}
	}

	// Export Cache
	data, err := json.MarshalIndent(cache, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(directoryCache+file, data, MODE_FILE)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Exported %6d Entries to %s Cache\n", len(*cache), label)

	return nil
}

func ImportMetadata(file string, label string) (map[int]*api.XtreamMetadata, error) {
	metadata := map[int]*api.XtreamMetadata{}

	data, err := os.ReadFile(directoryCache + file)
	if err != nil && !os.IsNotExist(err) {
		return map[int]*api.XtreamMetadata{}, err
	}

	if data == nil {
		data = []byte("{}") // Create Empty JSON Object if file doesn't exist.
	}

	err = json.Unmarshal(data, &metadata)
	if err != nil {
		return map[int]*api.XtreamMetadata{}, err
	}

	log.Printf("[INFO] Imported %6d Entries from %s Metadata\n", len(metadata), label)

	return metadata, nil
}

func ExportMetadata(metadata *map[int]*api.XtreamMetadata, file string, label string) error {
	// Create Root Directory
	err := os.MkdirAll(directoryCache, MODE_DIR)
	if err != nil {
		return err
	}

	// Export Metadata
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(directoryCache+file, data, MODE_FILE)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Exported %6d Entries to %s Metadata\n", len(*metadata), label)

	return nil
}
