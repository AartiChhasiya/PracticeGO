package entity

type Post struct {
	ID    int64  `json:"id"`
	Title string `json:"Title"`
	Text  string `json:"Text"`
}
