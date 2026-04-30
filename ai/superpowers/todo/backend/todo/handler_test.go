package todo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func newTestHandler(t *testing.T) (*Handler, *Store) {
	t.Helper()
	s := NewStore()
	return NewHandler(s), s
}

func newJSONRequest(method, path, body string) (*echo.Echo, echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	return e, c, rec
}

func TestHandler_Create_Returns201WithBody(t *testing.T) {
	t.Parallel()
	h, _ := newTestHandler(t)
	_, c, rec := newJSONRequest(http.MethodPost, "/api/todos", `{"title":"buy milk"}`)
	assert.NoError(t, h.Create(c))
	assert.Equal(t, http.StatusCreated, rec.Code)
	assert.Contains(t, rec.Body.String(), `"title":"buy milk"`)
	assert.Contains(t, rec.Body.String(), `"priority":"medium"`)
}

func TestHandler_Create_Returns400OnEmptyTitle(t *testing.T) {
	t.Parallel()
	h, _ := newTestHandler(t)
	_, c, rec := newJSONRequest(http.MethodPost, "/api/todos", `{"title":"  "}`)
	assert.NoError(t, h.Create(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"code":"validation_failed"`)
	assert.Contains(t, rec.Body.String(), `"field":"title"`)
}

func TestHandler_Create_Returns400OnInvalidJSON(t *testing.T) {
	t.Parallel()
	h, _ := newTestHandler(t)
	_, c, rec := newJSONRequest(http.MethodPost, "/api/todos", `{not json`)
	assert.NoError(t, h.Create(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"code":"invalid_json"`)
}

func TestHandler_Delete_Returns204(t *testing.T) {
	t.Parallel()
	h, s := newTestHandler(t)
	added := s.Add(NewTodo{Title: "x"})

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/todos/"+added.ID, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(added.ID)

	assert.NoError(t, h.Delete(c))
	assert.Equal(t, http.StatusNoContent, rec.Code)
	_, ok := s.Get(added.ID)
	assert.False(t, ok)
}

func TestHandler_Delete_Returns404OnUnknownID(t *testing.T) {
	t.Parallel()
	h, _ := newTestHandler(t)

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/api/todos/nope", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("nope")

	assert.NoError(t, h.Delete(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), `"code":"not_found"`)
}
