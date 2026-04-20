package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/CaptainKills/xtream-api/api"
)

const (
	MODE_DIR  = 0o755
	MODE_FILE = 0o644

	debugPercent = 100
)

var (
	directoryRoot        = "media"
	directoryLivestreams = filepath.Join(directoryRoot, "live")
	directoryMovies      = filepath.Join(directoryRoot, "movies")
	directorySeries      = filepath.Join(directoryRoot, "series")
)


func Export[T api.Stream](client *api.XtreamClient, streams *map[int]T, metadata *map[int]*api.XtreamMetadata, dir string, label string) error {
	updated_streams := 0
	updated_images := 0
	updated_nfos := 0

	if len(*streams) == 0 {
		log.Printf("[INFO] (%s) No available streams for export\n", label)
		return nil
	}

	log.Printf("[INFO] (%s) Exporting Streams...\n", label)

	// Create Root Directory
	err := os.MkdirAll(dir, MODE_DIR)
	if err != nil {
		return err
	}

	client.ResetRateStats()

	length := len(*streams)
	i := 0
	for id, stream := range *streams {
		// Output Progress Information
		if length >= debugPercent && i%(length/debugPercent) == 0 || i == length-1 {
			percentage := float64(i) / float64(length) * 100
			log.Printf("[DEBUG] (%s) Export Progress: %6d / %6d (%6.2f%%, %7.2f req/min)\tSTRM: %6d, IMG: %6d, NFO: %6d\n", label, i+1, length, percentage, client.GetRequestRate(), updated_streams, updated_images, updated_nfos)
		}

		updated_stream, updated_image, updated_nfo, err := stream.Export(client, dir)
		if err != nil {
			log.Printf("[ERROR] (%s) Failed to export Streams (%d): %v\n", label, id, err)
			continue
		}

		// Metadata Analysis
		md, ok := (*metadata)[id]
		if !ok {
			md = &api.XtreamMetadata{}
		}

		if updated_stream > 0 {
			md.Strm = true
		}

		if updated_image > 0 {
			md.Image = true
		}

		switch any(stream).(type) {
		case api.LiveStream:
			md.Nfo = true // LiveStream does not have NFO, so set to true to disable update check
		default:
			if updated_nfo > 0 {
				md.Nfo = true
			}
		}

		(*metadata)[id] = md

		// Update Counters
		updated_streams += updated_stream
		updated_images += updated_image
		updated_nfos += updated_nfo
		i++
	}

	log.Printf("[INFO] (%s) Streams Processed: %d, Streams Updated: %d, Images Updated: %d, Metadata Updated: %d\n", label, length, updated_streams, updated_images, updated_nfos)

	return nil
}
