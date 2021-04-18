package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

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

	req := httptest.NewRequest(echo.GET, "/api/articles", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	assert.NoError(t, h.GetArticle(c))

	if assert.Equal(t, http.StatusOK, rec.Code) {
		//var aa articleListResponse
		//err := json.Unmarshal(rec.Body.Bytes(), &aa)
		//assert.NoError(t, err)
		//assert.Equal(t, 2, aa.ArticlesCount)
	}

}
