package store

import (
	"github.com/google/uuid"
	"github.com/kenshin579/tutorials-go/go-validation/article/model"
)

type ArticleStore struct {
	storeList []model.Article
}

func NewArticleStore() *ArticleStore { //todo: 어떻게 interface로 반환을 해야 하나?
	return &ArticleStore{
		storeList: make([]model.Article, 0),
	}
}

func (as *ArticleStore) Create(request *model.ArticleRequest) (model.ArticleResponse, error) {
	a := model.Article{
		ArticleID:   uuid.New().String(),
		Title:       request.Title,
		Description: request.Description,
		Body:        request.Body,
	}
	as.storeList = append(as.storeList, a)
	response := model.ArticleResponse{
		ArticleID:   a.ArticleID,
		Title:       a.Title,
		Description: a.Description,
		Body:        a.Body,
	}
	return response, nil
}

func (as *ArticleStore) Delete(articleID string) error {
	temp := as.storeList[:0]
	for _, article := range as.storeList {
		if article.ArticleID != articleID {
			temp = append(temp, article)
		}
	}
	as.storeList = temp
	return model.ErrArticleNotFound
}

func (as *ArticleStore) GetByID(articleID string) (*model.Article, error) {
	for _, article := range as.storeList {
		if article.ArticleID == articleID {
			return &article, nil
		}
	}
	return nil, model.ErrArticleNotFound
}

func (as *ArticleStore) List() ([]model.Article, error) {
	return as.storeList, nil
}
