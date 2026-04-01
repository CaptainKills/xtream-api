package api

import "strings"

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

func (l LiveStream) Export(dir string, url string) (int, error) {
	l.Name = strings.ReplaceAll(l.Name, "/", "_")

	pathDirectory := dir + l.Name
	pathFile := dir + l.Name + "/" + l.Name + ".strm"

	updated, err := WriteStream(pathDirectory, pathFile, url)
	if err != nil {
		return updated, err
	}

	return updated, nil
}
