package models

type Bookmark struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Tags    []*Tag `json:"tags"`
	Color   *Color `json:"color"`
}
