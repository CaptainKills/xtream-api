package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golang.org/x/time/rate"
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
)

var xmlHeader = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>` + "\n"

type XtreamClient struct {
	url      string
	username string
	password string

	account Account
	Options XtreamOptions

	limiter    *rate.Limiter
	httpClient http.Client

	requestCount  int
	totalWaitTime time.Duration
	startTime     time.Time
}

type XtreamOptions struct {
	LaunchTime time.Time

	ImagesEnabled   bool
	MetadataEnabled bool

	RequestPerMinute time.Duration
	RequestTimeout   time.Duration

	BannedLiveStreams []string
	BannedMovies      []string
	BannedSeries      []string
}

type XtreamMetadata struct {
	Strm  bool `json:"strm"`
	Nfo   bool `json:"nfo"`
	Image bool `json:"image"`
}

func NewClient(url string, username string, password string, options XtreamOptions) *XtreamClient {
	var limiter *rate.Limiter
	if options.RequestPerMinute != 0 {
		every := rate.Every(60 * time.Second / options.RequestPerMinute)
		limiter = rate.NewLimiter(every, 1)
	} else {
		limiter = rate.NewLimiter(rate.Inf, 1)
	}

	var httpClient http.Client
	if options.RequestTimeout != 0 {
		httpClient = http.Client{Timeout: options.RequestTimeout * time.Second}
	} else {
		httpClient = *http.DefaultClient
	}

	return &XtreamClient{
		url:      url,
		username: username,
		password: password,
		Options:  options,

		limiter:    limiter,
		httpClient: httpClient,
	}
}

func (c *XtreamClient) buildURL(stream string, id int, ext string) string {
	protocol := c.account.ServerInfo.Protocol
	domain := c.account.ServerInfo.URL

	var port string
	switch protocol {
	case "http":
		port = c.account.ServerInfo.HttpPort
	case "https":
		port = c.account.ServerInfo.HttpsPort
	default:
		log.Printf("[WARNING] Unknown Server Protocol. Expected http/https, Got %s", protocol)
		port = "https"
	}

	username := c.username
	password := c.password

	return fmt.Sprintf(queryUrl, protocol, domain, port, stream, username, password, id, ext)
}

func (c *XtreamClient) sendRequest(query string) ([]byte, error) {
	waitStart := time.Now()
	err := c.limiter.Wait(context.Background())
	if err != nil {
		return []byte{}, err
	}
	c.totalWaitTime += time.Since(waitStart)
	c.requestCount++

	req, err := http.NewRequest(http.MethodGet, query, nil)
	if err != nil {
		return []byte{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return []byte{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return []byte{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}

func (c *XtreamClient) ResetRateStats() {
	c.startTime = time.Now()
	c.requestCount = 0
	c.totalWaitTime = time.Duration(0)
}

func (c *XtreamClient) GetRequestRate() float64 {
	elapsed := time.Since(c.startTime)
	if elapsed == 0 || c.requestCount == 0 {
		return 0.0
	}

	actualRPM := float64(c.requestCount) / elapsed.Minutes()

	return actualRPM
}
