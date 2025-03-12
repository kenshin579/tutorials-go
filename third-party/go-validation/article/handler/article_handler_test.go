package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/kenshin579/tutorials-go/go-validation/article/exception"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/kenshin579/tutorials-go/go-validation/article/store"

	"github.com/kenshin579/tutorials-go/go-validation/article/router"
	"github.com/kenshin579/tutorials-go/go-validation/article/usecase"

	"github.com/kenshin579/tutorials-go/go-validation/article/model"
	"github.com/labstack/echo"
)

var (
	h  *Handler
	e  *echo.Echo
	au model.ArticleUsecase
)

func setup() {
	e = router.New()
	as := store.NewArticleStore()
	au = usecase.NewArticleUsecase(as)
	h = NewHandler(au)
}

func teardown() {
}

func TestCreateArticle(t *testing.T) {
	setup()
	defer teardown()

	articleRequest := model.ArticleRequest{
		Title:       "title1",
		Description: "description1",
		Body:        "this is a body",
	}
	pbytes, _ := json.Marshal(articleRequest)
	buff := bytes.NewBuffer(pbytes)

	req := httptest.NewRequest(echo.POST, "/api/articles", buff)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.CreateArticle(c))

	if assert.Equal(t, http.StatusCreated, rec.Code) {
		var aa model.ArticleResponse
		err := json.Unmarshal(rec.Body.Bytes(), &aa)
		assert.NoError(t, err)
		assert.Equal(t, "title1", aa.Title)
	}
}

func TestCreateArticle_필요한_Request가_없는_경우_Err를_반환한다(t *testing.T) {
	setup()
	defer teardown()

	articleRequest := model.ArticleRequest{
		Title: "title1",
	}

	pbytes, _ := json.Marshal(articleRequest)
	buff := bytes.NewBuffer(pbytes)

	req := httptest.NewRequest(echo.POST, "/api/articles", buff)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err := h.CreateArticle(c)
	assert.Error(t, err)

	if validationError, ok := err.(validator.ValidationErrors); ok {
		assert.Contains(t, validationError.Error(), "ArticleRequest.Description")
	}
}

func TestGetArticle_존재하지_않는_데이터_조회시_Err를_반환한다(t *testing.T) {
	setup()
	defer teardown()

	articleID := 1

	req := httptest.NewRequest(echo.GET, "/api/articles"+strconv.Itoa(articleID), strings.NewReader(""))

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("articles/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(articleID))

	err := h.GetArticle(c)

	assert.Error(t, err)

	if customError, ok := err.(*exception.CustomError); ok {
		assert.Equal(t, http.StatusNotFound, customError.HttpCode())
		assert.Equal(t, exception.ErrArticleNotFound.Error(), customError.Error())
	}
}
