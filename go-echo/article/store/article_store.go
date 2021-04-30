package store

import (
	"errors"

	"github.com/google/uuid"
	"github.com/kenshin579/tutorials-go/go-echo/article/model"
)

var (
	ErrNotFound = errors.New("Article not found")
)

type articleStore struct {
	storeList []model.Article
}

func NewArticleStore() *articleStore {
	return &articleStore{
		storeList: make([]model.Article, 0),
	}
}

func (as *articleStore) Create(request *model.ArticleRequest) error {
	a := model.Article{
		ArticleID:   uuid.New().String(),
		Title:       request.Title,
		Description: request.Description,
		Body:        request.Body,
	}
	as.storeList = append(as.storeList, a)
	return nil
}

func (as *articleStore) Delete(articleID string) error {
	temp := as.storeList[:0]
	for _, article := range as.storeList {
		if article.ArticleID != articleID {
			temp = append(temp, article)
		}
	}
	as.storeList = temp
	return ErrNotFound
}

func (as *articleStore) GetByID(articleID string) (*model.Article, error) {
	for _, article := range as.storeList {
		if article.ArticleID == articleID {
			return &article, nil
		}
	}
	return nil, ErrNotFound
}

func (as *articleStore) List() ([]model.Article, error) {
	return as.storeList, nil
}
