package api

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	queryApi = "%s/player_api.php?username=%s&password=%s&action=%s"
	queryEpg = "%s/xmltv.php?username=%s&password=%s"
	queryUrl = "%s://%s:%s/%s/%s/%s/%d.%s"

	actionAccountInfo      = ""
	actionLiveCategories   = "get_live_categories"
	actionMovieCategories  = "get_vod_categories"
	actionSeriesCategories = "get_series_categories"
	actionLivestreams      = "get_live_streams"
	actionMovies           = "get_vod_streams"
	actionSeries           = "get_series"
	actionMovieInfo        = "get_vod_info&vod_id=%d"
	actionSeriesInfo       = "get_series_info&series_id=%d"

	directoryRoot        = "media/"
	directoryLivestreams = directoryRoot + "live/"
	directoryMovies      = directoryRoot + "movies/"
	directorySeries      = directoryRoot + "series/"

	debugPercent = 100
)

type XtreamClient struct {
	url      string
	username string
	password string

	account          Account
	liveCategories   []Category
	movieCategories  []Category
	seriesCategories []Category
	livestreams      []LiveStream
	movies           []Movie
	series           []Series

	enableImages bool
}

func NewClient(url string, username string, password string, enableImages bool) *XtreamClient {
	return &XtreamClient{
		url:      url,
		username: username,
		password: password,

		account:          Account{},
		liveCategories:   []Category{},
		movieCategories:  []Category{},
		seriesCategories: []Category{},
		livestreams:      []LiveStream{},
		movies:           []Movie{},
		series:           []Series{},

		enableImages: enableImages,
	}
}

func (c *XtreamClient) buildURL(stream string, id int, ext string) (string, error) {
	protocol := c.account.ServerInfo.Protocol
	domain := c.account.ServerInfo.URL

	var port string
	switch protocol {
	case "http":
		port = c.account.ServerInfo.HttpPort
	case "https":
		port = c.account.ServerInfo.HttpsPort
	default:
		err := fmt.Sprintf("Unknown Server Protocol. Expected http/https, Got %s", protocol)
		return "", errors.New(err)
	}

	username := c.username
	password := c.password

	return fmt.Sprintf(queryUrl, protocol, domain, port, stream, username, password, id, ext), nil
}

func (c *XtreamClient) ExportLiveStreams() error {
	updated_streams := 0
	updated_images := 0

	if len(c.livestreams) == 0 {
		return errors.New("No available LiveStreams for export!")
	}

	log.Printf("[INFO] Exporting LiveStreams...")

	// Create Root Directory
	err := os.MkdirAll(directoryLivestreams, 0o750)
	if err != nil {
		return err
	}

	length := len(c.livestreams)
	for i := range c.livestreams {
		// Output Progress Information
		if i%(length/debugPercent) == 0 || i == length-1 {
			percentage := float64(i) / float64(length) * 100
			log.Printf("[DEBUG] (LiveStream) Export Progress: %6d / %6d (%6.2f%%), STRM: %6d, IMG: %6d\n", i+1, length, percentage, updated_streams, updated_images)
		}

		livestream := c.livestreams[i]

		url, err := c.buildURL(livestream.StreamType, livestream.Id, c.account.UserInfo.AllowedOutputFormats[0])
		if err != nil {
			return err
		}

		updated_stream, updated_image, err := livestream.Export(directoryLivestreams, url, c.enableImages)
		if err != nil {
			return err
		}
		updated_streams += updated_stream
		updated_images += updated_image
	}

	log.Printf("[INFO] LiveStreams Processed: %d, LiveStreams Updated: %d, Images Updated: %d\n", length, updated_streams, updated_images)

	return nil
}

func (c *XtreamClient) ExportMovies() error {
	updated_streams := 0
	updated_images := 0

	if len(c.movies) == 0 {
		return errors.New("No available Movies for export!")
	}

	log.Printf("[INFO] Exporting Movies...")

	// Create Root Directory
	err := os.MkdirAll(directoryMovies, 0o750)
	if err != nil {
		return err
	}

	length := len(c.movies)
	for i := range c.movies {
		// Output Progress Information
		if i%(length/debugPercent) == 0 || i == length-1 {
			percentage := float64(i) / float64(length) * 100
			log.Printf("[DEBUG] (  Movies  ) Export Progress: %6d / %6d (%6.2f%%), STRM: %6d, IMG: %6d\n", i+1, length, percentage, updated_streams, updated_images)
		}

		movie := c.movies[i]

		url, err := c.buildURL(movie.StreamType, movie.Id, movie.Extension)
		if err != nil {
			return err
		}

		updated_stream, updated_image, err := movie.Export(directoryMovies, url, c.enableImages)
		if err != nil {
			return err
		}
		updated_streams += updated_stream
		updated_images += updated_image
	}

	log.Printf("[INFO] Movies Processed: %d, Movies Updated: %d, Images Updated: %d\n", length, updated_streams, updated_images)

	return nil
}

func (c *XtreamClient) ExportSeries() error {
	updated_streams := 0
	updated_images := 0

	if len(c.series) == 0 {
		return errors.New("No available Series for export!")
	}

	log.Printf("[INFO] Exporting Series...")

	err := os.MkdirAll(directorySeries, 0o750)
	if err != nil {
		return err
	}

	length := len(c.series)
	for i := range c.series {
		// Output Progress Information
		if i%(length/debugPercent) == 0 || i == length-1 {
			percentage := float64(i) / float64(length) * 100
			log.Printf("[DEBUG] (  Series  ) Export Progress: %6d / %6d (%6.2f%%), STRM: %6d, IMG: %6d\n", i+1, length, percentage, updated_streams, updated_images)
		}

		show := c.series[i]
		info, err := c.GetSeriesInfo(show.Id)
		if os.IsTimeout(err) {
			log.Println("[WARNING] Idleing for 10 Minutes...")
			time.Sleep(10 * time.Minute)

			info, err = c.GetSeriesInfo(show.Id)
			if err != nil {
				return err
			}
		} else if err != nil {
			log.Printf("[ERROR] Failed to get Series Info: %q\n", err)
			continue
		}

		pathDirectory := directorySeries + show.Name
		err = os.MkdirAll(pathDirectory, 0o750)
		if err != nil {
			return err
		}

		updated_image, err := show.Export(pathDirectory, show.Cover, c.enableImages)
		if err != nil {
			return err
		}
		updated_images += updated_image

		for j := range info.Seasons {
			season := info.Seasons[j]
			directory := pathDirectory + "/" + season.Name
			err := os.MkdirAll(directory, 0o750)
			if err != nil {
				return err
			}

			index := strconv.Itoa(season.Number)
			episodes := info.Episodes[index]
			for k := range episodes {
				episode := episodes[k]

				id, err := strconv.Atoi(episode.Id)
				if err != nil {
					return err
				}

				url, err := c.buildURL("series", id, episode.Extension)
				if err != nil {
					return err
				}

				updated_stream, err := episode.Export(directory, url)
				if err != nil {
					return err
				}

				updated_streams += updated_stream
			}
		}
	}

	log.Printf("[INFO] Series Processed: %d, Series Updated: %d, Images Updated: %d\n", length, updated_streams, updated_images)

	return nil
}

func (c *XtreamClient) ValidateLiveStreams() error {
	total_streams := 0
	total_covers := 0

	log.Printf("[INFO] Validating LiveStreams...")

	root, err := os.ReadDir(directoryLivestreams)
	if err != nil {
		return err
	}

	for i := range root {
		dir := root[i]

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

		for j := range subdir {
			file := subdir[j]

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
			err := c.cleanDirectory(directoryLivestreams+dir.Name(), subdir)
			if err != nil {
				return err
			}
		}
	}

	log.Printf("[INFO] (LiveStreams) Validated %6d Streams, %6d Covers\n", total_streams, total_covers)

	return nil
}

func (c *XtreamClient) ValidateMovies() error {
	total_streams := 0
	total_covers := 0

	log.Printf("[INFO] Validating Movies...")

	root, err := os.ReadDir(directoryMovies)
	if err != nil {
		return err
	}

	for i := range root {
		dir := root[i]

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

		for j := range subdir {
			file := subdir[j]

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
			err := c.cleanDirectory(directoryMovies+dir.Name(), subdir)
			if err != nil {
				return err
			}
		}
	}

	log.Printf("[INFO] (  Movies  ) Validated %6d Streams, %6d Covers\n", total_streams, total_covers)

	return nil
}

func (c *XtreamClient) ValidateSeries() error {
	total_streams := 0
	total_covers := 0

	log.Printf("[INFO] Validating Series...")

	root, err := os.ReadDir(directorySeries)
	if err != nil {
		return err
	}

	for i := range root {
		dir := root[i]

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

		for j := range subdir {
			file := subdir[j]

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
			err := c.cleanDirectory(directorySeries+dir.Name(), subdir)
			if err != nil {
				return err
			}
		}
	}

	log.Printf("[INFO] (  Series  ) Validated %6d Streams, %6d Covers\n", total_streams, total_covers)

	return nil
}

func (c *XtreamClient) cleanDirectory(root string, dir []os.DirEntry) error {
	streams := []os.DirEntry{}
	images := []os.DirEntry{}

	for i := range dir {
		file := dir[i]

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

		for i := range streams {
			stream := streams[i]

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

		for i := range images {
			image := images[i]

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
