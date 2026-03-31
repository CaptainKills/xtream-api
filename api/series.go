package api

type Series struct {
	// BackdropPath   []string `json:"backdrop_path"`
	Cast           string `json:"cast"`
	CategoryId     string `json:"category_id"` // int
	CategoryIds    []int  `json:"category_ids"`
	Cover          string `json:"cover"`
	Director       string `json:"director"`
	EpisodeRunTime string `json:"episode_run_time"` // int
	Genre          string `json:"genre"`
	Id             int    `json:"series_id"`
	LastModified   string `json:"last_modified"` // time.Time
	Name           string `json:"name"`
	Number         int    `json:"num"`
	Plot           string `json:"plot"`
	Rating         string `json:"rating"`        // float64
	Rating5Based   string `json:"rating_5based"` // float64
	ReleaseDate    string `json:"releaseDate"`   // time.Time
	ReleaseDate2   string `json:"release_date"`  // time.Time
	TMDB           string `json:"tmdb"`          // int
	Trailer        string `json:"youtube_trailer"`
}

type SeriesInfo struct{}
