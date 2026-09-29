package metrics

import (
	"time"

	"github.com/CaptainKills/xtream-api/api"
)

type Statistics struct {
	App     AppStatistics     `json:"app"`
	Options api.XtreamOptions `json:"options"`
	Time    TimeStatistics    `json:"time"`

	Live   StreamStatistics `json:"live"`
	Movies StreamStatistics `json:"movies"`
	Series StreamStatistics `json:"series"`
}

type AppStatistics struct {
	Version  string `json:"version"`
	Commit   string `json:"commit"`
	Branch   string `json:"branch"`
	Date     string `json:"date"`
	Modified string `json:"modified"`
}

type TimeStatistics struct {
	Start    time.Time     `json:"start"`
	End      time.Time     `json:"end"`
	Duration time.Duration `json:"duration"`
	Next     time.Time     `json:"next"`
}

type StreamStatistics struct {
	// Program Statistics
	Streams    int `json:"streams"`
	Categories int `json:"categories"`
	Cached     int `json:"cached"`
	Metadata   int `json:"metadata"`

	FilteredStreams    int `json:"filtered_streams"`
	FilteredCategories int `json:"filtered_categories"`

	// Detailed Statistics
	UpdatedStreams  int `json:"updated_streams"`
	UpdatedImages   int `json:"updated_images"`
	UpdatedMetadata int `json:"updated_metadata"`

	RequestRate float64 `json:"request_rate"`

	ValidatedStreams  int `json:"validated_streams"`
	ValidatedImages   int `json:"validated_images"`
	ValidatedMetadata int `json:"validated_metadata"`
}

var statistics = Statistics{}

func UpdateApp(app AppStatistics) {
	statistics.App = app
}

func UpdateOptions(options api.XtreamOptions) {
	statistics.Options = options
}

func UpdateTime(time TimeStatistics) {
	statistics.Time = time
}

func UpdateLivestreams(stats StreamStatistics) {
	statistics.Live = stats
}

func UpdateMovies(stats StreamStatistics) {
	statistics.Movies = stats
}

func UpdateSeries(stats StreamStatistics) {
	statistics.Series = stats
}
