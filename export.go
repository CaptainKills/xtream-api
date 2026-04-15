package main

import (
	"log"
	"os"

	"github.com/CaptainKills/xtream-api/api"
)

func Export[T api.Stream](client *api.XtreamClient, streams *map[int]T, dir string, label string) error {
	updated_streams := 0
	updated_images := 0
	updated_nfos := 0

	if len(*streams) == 0 {
		log.Printf("[INFO] No available '%s' for export.\n", label)
		return nil
	}

	log.Printf("[INFO] Exporting '%s'...\n", label)

	// Create Root Directory
	err := os.MkdirAll(dir, 0o750)
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

		updated_streams += updated_stream
		updated_images += updated_image
		updated_nfos += updated_nfo
		i++
	}

	log.Printf("[INFO] (%s) Streams Processed: %d, Streams Updated: %d, Images Updated: %d, Metadata Updated: %d\n", label, length, updated_streams, updated_images, updated_nfos)

	return nil
}
