package main

import (
	"errors"
	"log"
	"os"

	"github.com/CaptainKills/xtream-api/api"
)

func ExportLiveStreams(client *api.XtreamClient, livestreams *map[int]api.LiveStream) error {
	updated_streams := 0
	updated_images := 0

	if len(*livestreams) == 0 {
		return errors.New("No available LiveStreams for export!")
	}

	log.Printf("[INFO] Exporting LiveStreams...")

	// Create Root Directory
	err := os.MkdirAll(directoryLivestreams, 0o750)
	if err != nil {
		return err
	}

	length := len(*livestreams)
	i := 0
	for _, livestream := range *livestreams {
		// Output Progress Information
		if i%(length/debugPercent) == 0 || i == length-1 {
			percentage := float64(i) / float64(length) * 100
			log.Printf("[DEBUG] (LiveStream) Export Progress: %6d / %6d (%6.2f%%), STRM: %6d, IMG: %6d\n", i+1, length, percentage, updated_streams, updated_images)
		}

		updated_stream, updated_image, err := livestream.Export(client, directoryLivestreams)
		if err != nil {
			log.Printf("[ERROR] Failed to export Livestream (%d): %v\n", livestream.Id, err)
			continue
		}

		updated_streams += updated_stream
		updated_images += updated_image
		i++
	}

	log.Printf("[INFO] LiveStreams Processed: %d, Streams Updated: %d, Images Updated: %d\n", length, updated_streams, updated_images)

	return nil
}

func ExportMovies(client *api.XtreamClient, movies *map[int]api.Movie) error {
	updated_streams := 0
	updated_images := 0
	updated_nfos := 0

	if len(*movies) == 0 {
		return errors.New("No available Movies for export!")
	}

	log.Printf("[INFO] Exporting Movies...")

	// Create Root Directory
	err := os.MkdirAll(directoryMovies, 0o750)
	if err != nil {
		return err
	}

	length := len(*movies)
	i := 0
	for _, movie := range *movies {
		// Output Progress Information
		if i%(length/debugPercent) == 0 || i == length-1 {
			percentage := float64(i) / float64(length) * 100
			log.Printf("[DEBUG] (  Movies  ) Export Progress: %6d / %6d (%6.2f%%), STRM: %6d, IMG: %6d, NFO: %6d\n", i+1, length, percentage, updated_streams, updated_images, updated_nfos)
		}

		updated_stream, updated_image, updated_nfo, err := movie.Export(client, directoryMovies)
		if err != nil {
			log.Printf("[ERROR] Failed to export Movie (%d): %v\n", movie.Id, err)
			continue
		}

		updated_streams += updated_stream
		updated_images += updated_image
		updated_nfos += updated_nfo
		i++
	}

	log.Printf("[INFO] Movies Processed: %d, Streams Updated: %d, Images Updated: %d, Metadata Updated: %d\n", length, updated_streams, updated_images, updated_nfos)

	return nil
}

func ExportSeries(client *api.XtreamClient, series *map[int]api.Series) error {
	updated_streams := 0
	updated_images := 0
	updated_nfos := 0

	if len(*series) == 0 {
		return errors.New("No available Series for export!")
	}

	log.Printf("[INFO] Exporting Series...")

	// Create Root Directory
	err := os.MkdirAll(directorySeries, 0o750)
	if err != nil {
		return err
	}

	length := len(*series)
	i := 0
	for _, show := range *series {
		// Output Progress Information
		if i%(length/debugPercent) == 0 || i == length-1 {
			percentage := float64(i) / float64(length) * 100
			log.Printf("[DEBUG] (  Series  ) Export Progress: %6d / %6d (%6.2f%%), STRM: %6d, IMG: %6d, NFO: %6d\n", i+1, length, percentage, updated_streams, updated_images, updated_nfos)
		}

		updated_stream, updated_image, updated_nfo, err := show.Export(client, directorySeries)
		if err != nil {
			log.Printf("[ERROR] Failed to export Series (%d): %v\n", show.Id, err)
			continue
		}

		updated_streams += updated_stream
		updated_images += updated_image
		updated_nfos += updated_nfo
		i++
	}

	log.Printf("[INFO] Series Processed: %d, Streams Updated: %d, Images Updated: %d, Metadata Updated: %d\n", length, updated_streams, updated_images, updated_nfos)

	return nil
}
