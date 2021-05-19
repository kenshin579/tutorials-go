package model

type PostStatus string

const (
	Unknown   PostStatus = "Unknown"
	None      PostStatus = "None"
	Draft     PostStatus = "Draft"
	Published PostStatus = "Published"
)

type ArticleWithOneofTag struct {
	Title      string `json:"title" validate:"required"`
	Body       string `json:"body" validate:"required"`
	PostStatus string `json:"postStatus" validate:"required,oneof=None Draft Published"`
}
