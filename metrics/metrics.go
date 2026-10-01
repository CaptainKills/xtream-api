package metrics

import (
	"time"

	"github.com/CaptainKills/xtream-api/api"
)

type Metrics struct {
	app     AppMetrics
	options api.XtreamOptions
	time    TimeMetrics

	live   StreamMetrics
	movies StreamMetrics
	series StreamMetrics
}

type AppMetrics struct {
	Version  string
	Commit   string
	Branch   string
	Date     string
	Modified string
}

type TimeMetrics struct {
	Start    time.Time
	End      time.Time
	Duration time.Duration
	Next     time.Time
}

type StreamMetrics struct {
	// Program Statistics
	Streams    int
	Categories int
	Cached     int
	Metadata   int

	FilteredStreams    int
	FilteredCategories int

	// Detailed Statistics
	UpdatedStreams  int
	UpdatedImages   int
	UpdatedMetadata int

	RequestRate float64

	ValidatedStreams  int
	ValidatedImages   int
	ValidatedMetadata int
}

var metrics = Metrics{}

func UpdateApp(app AppMetrics) {
	metrics.app = app
}

func UpdateOptions(options api.XtreamOptions) {
	metrics.options = options
}

func UpdateTime(time TimeMetrics) {
	metrics.time = time
}

func UpdateLivestreams(m StreamMetrics) {
	metrics.live = m
}

func UpdateMovies(m StreamMetrics) {
	metrics.movies = m
}

func UpdateSeries(m StreamMetrics) {
	metrics.series = m
}
