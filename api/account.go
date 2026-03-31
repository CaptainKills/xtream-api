package api

type Account struct {
	UserInfo   UserInfo   `json:"user_info"`
	ServerInfo ServerInfo `json:"server_info"`
}

type UserInfo struct {
	ActiveConnections    string   `json:"active_cons"` // int
	AllowedOutputFormats []string `json:"allowed_output_formats"`
	CreatedAt            string   `json:"created_at"`      // time.Time
	ExpiresAt            string   `json:"exp_date"`        // time.Time
	IsAuthorised         int      `json:"auth"`            // bool
	IsTrial              string   `json:"is_trial"`        // bool
	MaxConnections       string   `json:"max_connections"` // int
	Message              string   `json:"message"`
	Password             string   `json:"password"`
	Status               string   `json:"status"`
	Username             string   `json:"username"`
}

type ServerInfo struct {
	HttpPort     string `json:"port"`       // int
	HttpsPort    string `json:"https_port"` // int
	Process      bool   `json:"process"`
	Protocol     string `json:"server_protocol"`
	RtmpPort     string `json:"rtmp_port"`   // int
	TimeNow      string `json:"time_now"`    // time.Time
	TimestampNow string `json:timestamp_now` // time.Time
	Timezone     string `json:"timezone"`
	URL          string `json:"url"`
}
