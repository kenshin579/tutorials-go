package handler

import (
	"fmt"
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

	request := httptest.NewRequest(echo.GET, "/api/articles", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	responseRecorder := httptest.NewRecorder()
	context := e.NewContext(request, responseRecorder)
	assert.NoError(t, h.GetArticle(context))
	fmt.Println("responseRecorder", responseRecorder.Body)
	if assert.Equal(t, http.StatusOK, responseRecorder.Code) {
		//var aa articleListResponse
		//err := json.Unmarshal(responseRecorder.Body.Bytes(), &aa)
		//assert.NoError(t, err)
		//assert.Equal(t, 2, aa.ArticlesCount)
	}

}
