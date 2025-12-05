package router

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/kenshin579/tutorials-go/go-validation/article/model"
)

func TestValidatePostStatus(t *testing.T) {
	validator := NewValidator()

	request := model.ArticleRequest{
		Title:       "title1",
		Description: "d1",
		Body:        "this is a body",
		PostStatus:  model.Published,
	}

	err := validator.Validate(request)
	assert.NoError(t, err)
}

func TestValidatePostStatus_Invalidate한_값인_경우_Err를_반환한다(t *testing.T) {
	validator := NewValidator()

	request := model.ArticleRequest{
		Title:       "title1",
		Description: "d1",
		Body:        "this is a body",
		PostStatus:  "TestStatus",
	}

	err := validator.Validate(request)
	assert.Error(t, err)
}
