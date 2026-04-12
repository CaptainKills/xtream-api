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
	directoryCache       = directoryRoot + "cache/"

	debugPercent = 100
)

type XtreamClient struct {
	url      string
	username string
	password string

	account Account
	data    XtreamData
	old     XtreamData

	options XtreamOptions
	raw     XtreamRaw
}

type XtreamData struct {
	liveCategories   map[int]Category
	movieCategories  map[int]Category
	seriesCategories map[int]Category
	livestreams      map[int]LiveStream
	movies           map[int]Movie
	series           map[int]Series
}

type XtreamOptions struct {
	ImagesEnabled bool
	NfoEnabled    bool

	RequestPerMinute time.Duration
	RequestTimeout   time.Duration

	BannedLiveStreams []string
	BannedMovies      []string
	BannedSeries      []string
}

func NewClient(url string, username string, password string, options XtreamOptions) *XtreamClient {
	// Initialise Request Handler
	InitRequest(options.RequestPerMinute, options.RequestTimeout)

	return &XtreamClient{
		url:      url,
		username: username,
		password: password,
		options:  options,
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

	if len(c.data.livestreams) == 0 {
		return errors.New("No available LiveStreams for export!")
	}

	log.Printf("[INFO] Exporting LiveStreams...")

	// Create Root Directory
	err := os.MkdirAll(directoryLivestreams, 0o750)
	if err != nil {
		return err
	}

	length := len(c.data.livestreams)
	i := 0
	for _, livestream := range c.data.livestreams {
		// Output Progress Information
		if i%(length/debugPercent) == 0 || i == length-1 {
			percentage := float64(i) / float64(length) * 100
			log.Printf("[DEBUG] (LiveStream) Export Progress: %6d / %6d (%6.2f%%), STRM: %6d, IMG: %6d\n", i+1, length, percentage, updated_streams, updated_images)
		}

		url, err := c.buildURL(livestream.StreamType, livestream.Id, c.account.UserInfo.AllowedOutputFormats[0])
		if err != nil {
			return fmt.Errorf("(%d) %w", livestream.Id, err)
		}

		updated_stream, updated_image, err := livestream.Export(directoryLivestreams, url, c.options.ImagesEnabled)
		if err != nil {
			return fmt.Errorf("(%d) %w", livestream.Id, err)
		}
		updated_streams += updated_stream
		updated_images += updated_image
		i++
	}

	log.Printf("[INFO] LiveStreams Processed: %d, Streams Updated: %d, Images Updated: %d\n", length, updated_streams, updated_images)

	return nil
}

func (c *XtreamClient) ExportMovies() error {
	updated_streams := 0
	updated_images := 0
	updated_nfos := 0

	if len(c.data.movies) == 0 {
		return errors.New("No available Movies for export!")
	}

	log.Printf("[INFO] Exporting Movies...")

	// Create Root Directory
	err := os.MkdirAll(directoryMovies, 0o750)
	if err != nil {
		return err
	}

	length := len(c.data.movies)
	i := 0
	for _, movie := range c.data.movies {
		// Output Progress Information
		if i%(length/debugPercent) == 0 || i == length-1 {
			percentage := float64(i) / float64(length) * 100
			log.Printf("[DEBUG] (  Movies  ) Export Progress: %6d / %6d (%6.2f%%), STRM: %6d, IMG: %6d, NFO: %6d\n", i+1, length, percentage, updated_streams, updated_images, updated_nfos)
		}

		// info, err := c.GetMovieInfo(movie.Id)
		// if err != nil {
		// 	log.Printf("[ERROR] Failed to get Movie Info (%d): %q\n", movie.Id, err)
		// 	continue
		// }
		// movie.Info = info

		url, err := c.buildURL(movie.StreamType, movie.Id, movie.Extension)
		if err != nil {
			return fmt.Errorf("(%d) %w", movie.Id, err)
		}

		updated_stream, updated_image, updated_nfo, err := movie.Export(directoryMovies, url, c.options.ImagesEnabled, c.options.NfoEnabled)
		if err != nil {
			return fmt.Errorf("(%d) %w", movie.Id, err)
		}

		updated_streams += updated_stream
		updated_images += updated_image
		updated_nfos += updated_nfo
		i++
	}

	log.Printf("[INFO] Movies Processed: %d, Streams Updated: %d, Images Updated: %d, Metadata Updated: %d\n", length, updated_streams, updated_images, updated_nfos)

	return nil
}

func (c *XtreamClient) ExportSeries() error {
	updated_streams := 0
	updated_images := 0
	updated_nfos := 0

	if len(c.data.series) == 0 {
		return errors.New("No available Series for export!")
	}

	log.Printf("[INFO] Exporting Series...")

	err := os.MkdirAll(directorySeries, 0o750)
	if err != nil {
		return err
	}

	length := len(c.data.series)
	i := 0
	for _, show := range c.data.series {
		// Output Progress Information
		if i%(length/debugPercent) == 0 || i == length-1 {
			percentage := float64(i) / float64(length) * 100
			log.Printf("[DEBUG] (  Series  ) Export Progress: %6d / %6d (%6.2f%%), STRM: %6d, IMG: %6d, NFO: %6d\n", i+1, length, percentage, updated_streams, updated_images, updated_nfos)
		}

		info, err := c.GetSeriesInfo(show.Id)
		if err != nil {
			log.Printf("[ERROR] Failed to get Series Info (%d): %q\n", show.Id, err)
			continue
		}
		show.Info = info

		pathDirectory := directorySeries + show.Name
		err = os.MkdirAll(pathDirectory, 0o750)
		if err != nil {
			return fmt.Errorf("(%d) %w", show.Id, err)
		}

		updated_image, updated_nfo, err := show.Export(pathDirectory, show.Cover, c.options.ImagesEnabled, c.options.NfoEnabled)
		if err != nil {
			return fmt.Errorf("(%d) %w", show.Id, err)
		}
		updated_images += updated_image
		updated_nfos += updated_nfo

		for season, episodes := range info.Episodes {
			directory := pathDirectory + "/Season " + season
			err = os.MkdirAll(directory, 0o750)
			if err != nil {
				return fmt.Errorf("(%d) %w", show.Id, err)
			}

			for _, episode := range episodes {
				id, err := strconv.Atoi(episode.Id)
				if err != nil {
					return fmt.Errorf("(%d) %w", show.Id, err)
				}

				url, err := c.buildURL("series", id, episode.Extension)
				if err != nil {
					return fmt.Errorf("(%d) %w", show.Id, err)
				}

				updated_stream, updated_nfo, err := episode.Export(directory, url, c.options.NfoEnabled)
				if err != nil {
					return fmt.Errorf("(%d) %w", show.Id, err)
				}

				updated_streams += updated_stream
				updated_nfos += updated_nfo
			}
		}

		i++
	}

	log.Printf("[INFO] Series Processed: %d, Streams Updated: %d, Images Updated: %d, Metadata Updated: %d\n", length, updated_streams, updated_images, updated_nfos)

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
