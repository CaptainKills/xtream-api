package main

import (
	"log"
	"os"
	"strings"
)

func Validate(dir string, label string) error {
	total_streams := 0
	total_covers := 0
	total_metadata := 0

	log.Printf("[INFO] Validating '%s'...\n", label)

	root, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, subdir := range root {
		if !subdir.IsDir() {
			log.Printf("[WARNING] Found File in Root Directory: %s\n", subdir.Name())

			err := os.Remove(dir + subdir.Name())
			if err != nil {
				return err
			}
			continue
		}

		content, err := os.ReadDir(dir + subdir.Name())
		if err != nil {
			return err
		}

		nr_of_streams := 0
		nr_of_covers := 0
		nr_of_metadata := 0

		for _, file := range content {
			if strings.HasSuffix(file.Name(), ".strm") {
				nr_of_streams++
				total_streams++
			}

			if strings.HasPrefix(file.Name(), "cover") {
				nr_of_covers++
				total_covers++
			}

			if strings.HasSuffix(file.Name(), ".nfo") {
				nr_of_metadata++
				total_metadata++
			}

			if file.IsDir() {
				season, err := os.ReadDir(dir + subdir.Name() + "/" + file.Name())
				if err != nil {
					return err
				}

				for _, subfile := range season {
					if strings.HasSuffix(subfile.Name(), ".strm") {
						nr_of_streams++
						total_streams++
					}

					if strings.HasSuffix(subfile.Name(), ".nfo") {
						nr_of_metadata++
						total_metadata++
					}
				}
			}
		}

		switch label {
		case "Livestreams":
			if nr_of_streams == 0 || nr_of_streams > 1 || nr_of_covers > 1 || nr_of_metadata > 1 {
				log.Printf("[WARNING] (%s) Unexpected Number of Files: STRM=%2d, IMG=%2d, NFO=%2d | %s\n", label, nr_of_streams, nr_of_covers, nr_of_metadata, subdir.Name())
				err := cleanDirectory(directoryLivestreams+subdir.Name(), content)
				if err != nil {
					return err
				}
			}
		case "  Movies   ":
			if nr_of_streams == 0 || nr_of_streams > 1 || nr_of_covers > 1 || nr_of_metadata > 1 {
				log.Printf("[WARNING] (%s) Unexpected Number of Files: STRM=%2d, IMG=%2d, NFO=%2d | %s\n", label, nr_of_streams, nr_of_covers, nr_of_metadata, subdir.Name())
				err := cleanDirectory(directoryLivestreams+subdir.Name(), content)
				if err != nil {
					return err
				}
			}
		case "  Series   ":
			if nr_of_streams == 0 || nr_of_covers > 1 {
				log.Printf("[WARNING] (%s) Unexpected Number of Files: STRM=%2d, IMG=%2d, NFO=%2d | %s\n", label, nr_of_streams, nr_of_covers, nr_of_metadata, subdir.Name())
			}
		}
	}

	log.Printf("[INFO] (%s) Validated %6d Streams, %6d Covers, %6d Metadata\n", label, total_streams, total_covers, total_metadata)

	return nil
}

func cleanDirectory(root string, dir []os.DirEntry) error {
	streams := []os.DirEntry{}
	images := []os.DirEntry{}
	metadata := []os.DirEntry{}

	for _, file := range dir {
		if strings.HasSuffix(file.Name(), ".strm") {
			streams = append(streams, file)
		}

		if strings.HasPrefix(file.Name(), "cover") {
			images = append(images, file)
		}

		if strings.HasSuffix(file.Name(), ".nfo") {
			metadata = append(metadata, file)
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

	deletedMetadata := 0
	if len(metadata) > 1 {
		newestMetadata := metadata[0]
		newestInfo, err := os.Stat(root + "/" + newestMetadata.Name())
		if err != nil {
			return err
		}
		newestModtime := newestInfo.ModTime()

		for _, nfo := range metadata {
			info, err := os.Stat(root + "/" + nfo.Name())
			if err != nil {
				return err
			}
			modtime := info.ModTime()

			if modtime.After(newestModtime) {
				err := os.Remove(root + "/" + newestMetadata.Name())
				if err != nil {
					return err
				}

				newestMetadata = nfo
				newestModtime = modtime
				deletedMetadata++
			} else if modtime.Before(newestModtime) {
				err := os.Remove(root + "/" + nfo.Name())
				if err != nil {
					return err
				}
				deletedMetadata++
			}
		}
	}

	log.Printf("[DEBUG] Cleaned %d Streams, %d Images, %d Metadata\n", deletedStreams, deletedImages, deletedMetadata)

	return nil
}
