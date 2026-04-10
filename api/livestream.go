package api

import (
	"encoding/json"
	"fmt"
	"strings"
)

type LiveStream struct {
	Added        string `json:"added"`       // time.Time
	CategoryId   string `json:"category_id"` // int
	CategoryIds  []int  `json:"category_ids"`
	CustomSID    string `json:"custom_sid"`
	DirectSource string `json:"direct_source"`
	EpgId        string `json:"epg_channel_id"` // int
	Icon         string `json:"stream_icon"`
	Id           int    `json:"stream_id"`
	IsAdult      int    `json:"is_adult"` // bool
	Name         string `json:"name"`
	Number       int    `json:"num"`
	StreamType   string `json:"stream_type"`
	TvArchive    int    `json:"tv_archive"`
	// TvArchiveDuration string `json:"tv_archive_duration"` // int
	// CatchupDurationDays int  `json:"catchup_duration_days"`
	// HasCatchup          bool `json:"has_catchup"`
}

func (c *XtreamClient) GetLiveStreams() (map[int]LiveStream, error) {
	var livestreams []LiveStream
	query := fmt.Sprintf(queryApi, c.url, c.username, c.password, actionLivestreams)

	// Fetch LiveStreams
	resp, err := SendRequest(query)
	if err != nil {
		return map[int]LiveStream{}, err
	}

	err = json.Unmarshal(resp, &livestreams)
	if err != nil {
		return map[int]LiveStream{}, err
	}

	// Filter Banned LiveStreams
	for i := range livestreams {
		allowed := true
		livestream := livestreams[i]

		for j := range livestream.CategoryIds {
			id := livestream.CategoryIds[j]

			if _, ok := c.liveCategories[id]; !ok {
				allowed = false
			}
		}

		if allowed {
			c.livestreams[livestream.Id] = livestream
		}
	}

	return c.livestreams, nil
}

func (l LiveStream) Export(dir string, url string, enableImages bool) (int, int, error) {
	l.Name = strings.ReplaceAll(l.Name, "/", "_")

	pathDirectory := dir + l.Name
	pathFile := pathDirectory + "/" + l.Name + ".strm"
	pathImage := pathDirectory + "/cover" + GetImageExtension(l.Icon)

	// Write Stream to File
	updated_stream, err := WriteStream(pathDirectory, pathFile, url)
	if err != nil {
		return updated_stream, 0, err
	}

	// Write Image to File
	updated_image, err := WriteImage(pathDirectory, pathImage, l.Icon, enableImages)
	if err != nil {
		return updated_stream, updated_image, err
	}

	return updated_stream, updated_image, nil
}
