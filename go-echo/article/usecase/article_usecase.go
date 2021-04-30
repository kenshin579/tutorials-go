package usecase

import (
	"github.com/kenshin579/tutorials-go/go-echo/article/model"
	"github.com/labstack/gommon/log"
)

type articleUsecase struct {
	articleStore model.ArticleStore
}

func NewArticleUsecase(as model.ArticleStore) model.ArticleUsecase {
	return &articleUsecase{
		articleStore: as,
	}
}

func (a *articleUsecase) CreateArticle(request *model.ArticleRequest) error {
	return a.articleStore.Create(request)
}

func (a *articleUsecase) GetArticle(articleID string) (*model.ArticleResponse, error) {
	article, err := a.articleStore.GetByID(articleID)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return model.NewArticleResponse(article), nil
}

func (a *articleUsecase) DeleteArticle(articleID string) error {
	return a.articleStore.Delete(articleID)
}

func (a *articleUsecase) ListArticle() []model.ArticleResponse {
	result := make([]model.ArticleResponse, 0)
	articleList, err := a.articleStore.List()
	if err != nil {
		log.Error(err)
	}

	for _, article := range articleList {
		response := model.NewArticleResponse(&article)
		result = append(result, *response)
	}

	return result
}
