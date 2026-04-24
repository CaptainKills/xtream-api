package api

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/CaptainKills/xtream-api/utils"
)

type LiveStream struct {
	Added               string `json:"added"`       // time.Time
	CategoryId          string `json:"category_id"` // int
	CategoryIds         []int  `json:"category_ids"`
	CustomSID           string `json:"custom_sid"`
	DirectSource        string `json:"direct_source"`
	EpgId               string `json:"epg_channel_id"`
	Icon                string `json:"stream_icon"`
	Id                  int    `json:"stream_id"`
	IsAdult             int    `json:"is_adult"` // bool
	Name                string `json:"name"`
	Number              int    `json:"num"`
	StreamType          string `json:"stream_type"`
	TvArchive           int    `json:"tv_archive"`
	CatchupDurationDays int    `json:"catchup_duration_days"`
	HasCatchup          bool   `json:"has_catchup"`
}

func (c *XtreamClient) GetLiveStreams() (map[int]LiveStream, error) {
	var livestreams []LiveStream
	livestream_map := map[int]LiveStream{}

	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionLivestreams)

	// Fetch LiveStreams
	resp, err := c.sendRequest(query)
	if err != nil {
		return map[int]LiveStream{}, err
	}

	// Unmarshal LiveStreams
	err = json.Unmarshal(resp, &livestreams)
	if err != nil {
		return map[int]LiveStream{}, err
	}

	// Map LiveStreams
	for _, livestream := range livestreams {
		livestream_map[livestream.Id] = livestream
	}

	return livestream_map, nil
}

func (l LiveStream) Export(c *XtreamClient, dir string) (int, int, int, error) {
	updated_stream := 0
	updated_image := 0
	updated_nfo := 0

	l.Name = strings.ReplaceAll(l.Name, "/", "_")

	pathDirectory := filepath.Join(dir, l.Name)
	pathStream := filepath.Join(pathDirectory, l.Name+".strm")
	pathImage := filepath.Join(pathDirectory, "cover"+utils.GetImageExtension(l.Icon))
	url := c.buildURL(l.StreamType, l.Id, c.account.UserInfo.AllowedOutputFormats[0])

	// Create Subdirectory
	err := os.Mkdir(pathDirectory, 0o755)
	if err != nil && !os.IsExist(err) {
		return updated_stream, updated_image, updated_nfo, err
	}

	// Write Stream to File
	updated_stream, err = utils.WriteFile(pathStream, url)
	if err != nil {
		return updated_stream, updated_image, updated_nfo, err
	}

	// Write Image to File
	if c.Options.ImagesEnabled && !utils.ImageExists(pathImage) && strings.HasPrefix(l.Icon, "http") {
		image, err := c.sendRequest(l.Icon)
		if err != nil {
			// Ignore error for image fetching
			// log.Printf("[WARNING] Failed to fetch Image: %v\n", err)
		} else {
			updated_image, err = utils.WriteImage(pathImage, image)
			if err != nil {
				return updated_stream, updated_image, updated_nfo, err
			}
		}
	}

	return updated_stream, updated_image, updated_nfo, nil
}
