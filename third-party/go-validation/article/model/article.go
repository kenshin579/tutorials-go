package model

type PostStatus string

const (
	Unknown   PostStatus = "Unknown"
	None      PostStatus = "None"
	Draft     PostStatus = "Draft"
	Published PostStatus = "Published"
)

var AllPostStatus = []PostStatus{Unknown, None, Draft, Published}

type Article struct {
	ArticleID   string
	Title       string
	Description string
	Body        string
	PostStatus  PostStatus
}

type ArticleRequest struct {
	Title       string     `json:"title" validate:"required"`
	Description string     `json:"description" validate:"required"`
	Body        string     `json:"body" validate:"required"`
	PostStatus  PostStatus `json:"postStatus" validate:"postStatus"`
}

type ArticleResponse struct {
	ArticleID   string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

func NewArticleResponse(a *Article) *ArticleResponse {
	return &ArticleResponse{
		ArticleID:   a.ArticleID,
		Title:       a.Title,
		Description: a.Description,
		Body:        a.Body,
	}
}

type ArticleUsecase interface {
	GetArticle(string) (*ArticleResponse, error)
	CreateArticle(*ArticleRequest) (ArticleResponse, error)
	DeleteArticle(string) error
	ListArticle() []ArticleResponse
}

type ArticleStore interface {
	GetByID(string) (*Article, error)
	Create(*ArticleRequest) (ArticleResponse, error)
	Delete(string) error
	List() ([]Article, error)
}
