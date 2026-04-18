package main

import (
	"log"
	"os"

	"github.com/CaptainKills/xtream-api/api"
)

func Export[T api.Stream](client *api.XtreamClient, streams *map[int]T, metadata *map[int]*api.XtreamMetadata, dir string, label string) error {
	updated_streams := 0
	updated_images := 0
	updated_nfos := 0

	if len(*streams) == 0 {
		log.Printf("[INFO] No available '%s' for export.\n", label)
		return nil
	}

	log.Printf("[INFO] Exporting '%s'...\n", label)

	// Create Root Directory
	err := os.MkdirAll(dir, MODE_DIR)
	if err != nil {
		return err
	}

	length := len(*streams)
	i := 0
	for id, stream := range *streams {
		// Output Progress Information
		if length >= debugPercent && i%(length/debugPercent) == 0 || i == length-1 {
			percentage := float64(i) / float64(length) * 100
			log.Printf("[DEBUG] (%s) Export Progress: %6d / %6d (%6.2f%%)\tSTRM: %6d, IMG: %6d, NFO: %6d\n", label, i+1, length, percentage, updated_streams, updated_images, updated_nfos)
		}

		updated_stream, updated_image, updated_nfo, err := stream.Export(client, dir)
		if err != nil {
			log.Printf("[ERROR] Failed to export '%s' (%d): %v\n", label, id, err)
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
