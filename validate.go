package main

import (
	"log"
	"os"
	"strings"
)

func ValidateLiveStreams() error {
	total_streams := 0
	total_covers := 0

	log.Printf("[INFO] Validating LiveStreams...")

	root, err := os.ReadDir(directoryLivestreams)
	if err != nil {
		return err
	}

	for _, dir := range root {
		if !dir.IsDir() {
			log.Printf("[WARNING] Found File in Root Directory: %s\n", dir.Name())

			err := os.Remove(directorySeries + dir.Name())
			if err != nil {
				return err
			}
			continue
		}

		subdir, err := os.ReadDir(directoryLivestreams + dir.Name())
		if err != nil {
			return err
		}

		nr_of_streams := 0
		nr_of_covers := 0

		for _, file := range subdir {
			if strings.HasSuffix(file.Name(), ".strm") {
				nr_of_streams++
				total_streams++
			}

			if strings.HasPrefix(file.Name(), "cover") {
				nr_of_covers++
				total_covers++
			}
		}

		if nr_of_streams == 0 || nr_of_streams > 1 || nr_of_covers > 1 {
			log.Printf("[WARNING] Unexpected Number of Files: STRM=%d, IMG=%d | %s\n", nr_of_streams, nr_of_covers, dir.Name())
			err := cleanDirectory(directoryLivestreams+dir.Name(), subdir)
			if err != nil {
				return err
			}
		}
	}

	log.Printf("[INFO] (LiveStreams) Validated %6d Streams, %6d Covers\n", total_streams, total_covers)

	return nil
}

func ValidateMovies() error {
	total_streams := 0
	total_covers := 0

	log.Printf("[INFO] Validating Movies...")

	root, err := os.ReadDir(directoryMovies)
	if err != nil {
		return err
	}

	for _, dir := range root {
		if !dir.IsDir() {
			log.Printf("[WARNING] Found File in Root Directory: %s\n", dir.Name())

			err := os.Remove(directorySeries + dir.Name())
			if err != nil {
				return err
			}
			continue
		}

		subdir, err := os.ReadDir(directoryMovies + dir.Name())
		if err != nil {
			return err
		}

		nr_of_streams := 0
		nr_of_covers := 0

		for _, file := range subdir {
			if strings.HasSuffix(file.Name(), ".strm") {
				nr_of_streams++
				total_streams++
			}

			if strings.HasPrefix(file.Name(), "cover") {
				nr_of_covers++
				total_covers++
			}
		}

		if nr_of_streams == 0 || nr_of_streams > 1 || nr_of_covers > 1 {
			log.Printf("[WARNING] Unexpected Number of Files: STRM=%d, IMG=%d | %s\n", nr_of_streams, nr_of_covers, dir.Name())
			err := cleanDirectory(directoryMovies+dir.Name(), subdir)
			if err != nil {
				return err
			}
		}
	}

	log.Printf("[INFO] (  Movies  ) Validated %6d Streams, %6d Covers\n", total_streams, total_covers)

	return nil
}

func ValidateSeries() error {
	total_streams := 0
	total_covers := 0

	log.Printf("[INFO] Validating Series...")

	root, err := os.ReadDir(directorySeries)
	if err != nil {
		return err
	}

	for _, dir := range root {
		if !dir.IsDir() {
			log.Printf("[WARNING] Found File in Root Directory: %s\n", dir.Name())

			err := os.Remove(directorySeries + dir.Name())
			if err != nil {
				return err
			}
			continue
		}

		subdir, err := os.ReadDir(directorySeries + dir.Name())
		if err != nil {
			return err
		}

		nr_of_streams := 0
		nr_of_covers := 0

		for _, file := range subdir {
			// if strings.HasSuffix(file.Name(), ".strm") {
			// 	nr_of_streams++
			// 	total_streams++
			// }

			if strings.HasPrefix(file.Name(), "cover") {
				nr_of_covers++
				total_covers++
			}
		}

		if nr_of_covers > 1 {
			log.Printf("[WARNING] Unexpected Number of Files: STRM=%d, IMG=%d | %s\n", nr_of_streams, nr_of_covers, dir.Name())
			err := cleanDirectory(directorySeries+dir.Name(), subdir)
			if err != nil {
				return err
			}
		}
	}

	log.Printf("[INFO] (  Series  ) Validated %6d Streams, %6d Covers\n", total_streams, total_covers)

	return nil
}

func cleanDirectory(root string, dir []os.DirEntry) error {
	streams := []os.DirEntry{}
	images := []os.DirEntry{}

	for _, file := range dir {
		if strings.HasSuffix(file.Name(), ".strm") {
			streams = append(streams, file)
		}

		if strings.HasPrefix(file.Name(), "cover") {
			images = append(images, file)
		}
	}

	deletedStreams := 0
	if len(streams) > 1 {
		newestStream := images[0]
		newestInfo, err := os.Stat(root + "/" + newestStream.Name())
		if err != nil {
			return err
		}
		newestModtime := newestInfo.ModTime()

		for _, stream := range streams {
			info, err := os.Stat(root + "/" + stream.Name())
			if err != nil {
				return err
			}
			modtime := info.ModTime()

			if modtime.After(newestModtime) {
				err := os.Remove(root + "/" + newestStream.Name())
				if err != nil {
					return err
				}

				newestStream = stream
				newestModtime = modtime
				deletedStreams++
			}
		}
	}

	deletedImages := 0
	if len(images) > 1 {
		newestImage := images[0]
		newestInfo, err := os.Stat(root + "/" + newestImage.Name())
		if err != nil {
			return err
		}
		newestModtime := newestInfo.ModTime()

		for _, image := range images {
			info, err := os.Stat(root + "/" + image.Name())
			if err != nil {
				return err
			}
			modtime := info.ModTime()

			if modtime.After(newestModtime) {
				err := os.Remove(root + "/" + newestImage.Name())
				if err != nil {
					return err
				}

				newestImage = image
				newestModtime = modtime
				deletedImages++
			} else if modtime.Before(newestModtime) {
				err := os.Remove(root + "/" + image.Name())
				if err != nil {
					return err
				}
				deletedImages++
			}
		}
	}

	log.Printf("[DEBUG] Cleaned %d Streams, %d Images\n", deletedStreams, deletedImages)

	return nil
}
