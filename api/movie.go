package api

type Movie struct {
	Added        string `json:"added"`       // time.Time
	CategoryId   string `json:"category_id"` // int
	CategoryIds  []int  `json:"category_ids"`
	CystomSID    string `json:"custom_sid"`
	DirectSource string `json:"direct_source"`
	Extension    string `json:"container_extension"`
	Icon         string `json:"stream_icon"`
	Id           int    `json:"stream_id"`
	IsAdult      int    `json:"is_adult"` // bool
	Name         string `json:"name"`
	Number       int    `json:"num"`
	Rating       string `json:"rating"` // float64
	// Rating5Based string `json:"rating_5based"` // float64
	StreamType string `json:"stream_type"`
	// TMDB       string `json:"tmdb"` // int
	Trailer string `json:"trailer"`
}

type MovieInfo struct{}
