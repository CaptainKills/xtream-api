package api

type Category struct {
	Id     string `json:"category_id"` // int
	Name   string `json:"category_name"`
	Parent int    `json:"parent_id"`
}
