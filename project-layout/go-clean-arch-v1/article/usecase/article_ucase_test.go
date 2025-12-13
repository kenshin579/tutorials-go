package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	ucase "github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/article/usecase"
	"github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/domain"
	"github.com/kenshin579/tutorials-go/project-layout/go-clean-arch-v1/domain/mocks"
)

func TestFetch(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title:   "Hello",
		Content: "Content",
	}

	mockListArtilce := make([]domain.Article, 0)
	mockListArtilce = append(mockListArtilce, mockArticle)

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("FetchArticle", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(mockListArtilce, "next-cursor", nil).Once()
		mockAuthor := domain.Author{
			ID:   1,
			Name: "Iman Tumorang",
		}
		mockAuthorrepo := new(mocks.AuthorRepository)
		mockAuthorrepo.On("GetArticleByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.FetchArticle(context.TODO(), cursor, num)
		cursorExpected := "next-cursor"
		assert.Equal(t, cursorExpected, nextCursor)
		assert.NotEmpty(t, nextCursor)
		assert.NoError(t, err)
		assert.Len(t, list, len(mockListArtilce))

		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepo.On("FetchArticle", mock.Anything, mock.AnythingOfType("string"),
			mock.AnythingOfType("int64")).Return(nil, "", errors.New("Unexpexted Error")).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)
		num := int64(1)
		cursor := "12"
		list, nextCursor, err := u.FetchArticle(context.TODO(), cursor, num)

		assert.Empty(t, nextCursor)
		assert.Error(t, err)
		assert.Len(t, list, 0)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestGetByID(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title:   "Hello",
		Content: "Content",
	}
	mockAuthor := domain.Author{
		ID:   1,
		Name: "Iman Tumorang",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("GetArticleByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockArticle, nil).Once()
		mockAuthorrepo := new(mocks.AuthorRepository)
		mockAuthorrepo.On("GetArticleByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)
		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

		a, err := u.GetArticleByID(context.TODO(), mockArticle.ID)

		assert.NoError(t, err)
		assert.NotNil(t, a)

		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("error-failed", func(t *testing.T) {
		mockArticleRepo.On("GetArticleByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Article{}, errors.New("Unexpected")).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

		a, err := u.GetArticleByID(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		assert.Equal(t, domain.Article{}, a)

		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestStore(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		tempMockArticle := mockArticle
		tempMockArticle.ID = 0
		mockArticleRepo.On("GetArticleByTitle", mock.Anything, mock.AnythingOfType("string")).Return(domain.Article{}, domain.ErrNotFound).Once()
		mockArticleRepo.On("StoreArticle", mock.Anything, mock.AnythingOfType("*domain.Article")).Return(nil).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

		err := u.StoreArticle(context.TODO(), &tempMockArticle)

		assert.NoError(t, err)
		assert.Equal(t, mockArticle.Title, tempMockArticle.Title)
		mockArticleRepo.AssertExpectations(t)
	})
	t.Run("existing-title", func(t *testing.T) {
		existingArticle := mockArticle
		mockArticleRepo.On("GetArticleByTitle", mock.Anything, mock.AnythingOfType("string")).Return(existingArticle, nil).Once()
		mockAuthor := domain.Author{
			ID:   1,
			Name: "Iman Tumorang",
		}
		mockAuthorrepo := new(mocks.AuthorRepository)
		mockAuthorrepo.On("GetArticleByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockAuthor, nil)

		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

		err := u.StoreArticle(context.TODO(), &mockArticle)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestDelete(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title:   "Hello",
		Content: "Content",
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("GetArticleByID", mock.Anything, mock.AnythingOfType("int64")).Return(mockArticle, nil).Once()

		mockArticleRepo.On("DeleteArticleByID", mock.Anything, mock.AnythingOfType("int64")).Return(nil).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

		err := u.DeleteArticleByID(context.TODO(), mockArticle.ID)

		assert.NoError(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("article-is-not-exist", func(t *testing.T) {
		mockArticleRepo.On("GetArticleByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Article{}, nil).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

		err := u.DeleteArticleByID(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})
	t.Run("error-happens-in-db", func(t *testing.T) {
		mockArticleRepo.On("GetArticleByID", mock.Anything, mock.AnythingOfType("int64")).Return(domain.Article{}, errors.New("Unexpected Error")).Once()

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

		err := u.DeleteArticleByID(context.TODO(), mockArticle.ID)

		assert.Error(t, err)
		mockArticleRepo.AssertExpectations(t)
		mockAuthorrepo.AssertExpectations(t)
	})

}

func TestUpdate(t *testing.T) {
	mockArticleRepo := new(mocks.ArticleRepository)
	mockArticle := domain.Article{
		Title:   "Hello",
		Content: "Content",
		ID:      23,
	}

	t.Run("success", func(t *testing.T) {
		mockArticleRepo.On("UpdateArticle", mock.Anything, &mockArticle).Once().Return(nil)

		mockAuthorrepo := new(mocks.AuthorRepository)
		u := ucase.NewArticleUsecase(mockArticleRepo, mockAuthorrepo, time.Second*2)

		err := u.UpdateArticle(context.TODO(), &mockArticle)
		assert.NoError(t, err)
		mockArticleRepo.AssertExpectations(t)
	})
}
