package todo

import (
	"encoding/json"
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

func TestHandler_Create_Returns400OnInvalidPriority(t *testing.T) {
	t.Parallel()
	h, _ := newTestHandler(t)
	_, c, rec := newJSONRequest(http.MethodPost, "/api/todos", `{"title":"x","priority":"urgent"}`)
	assert.NoError(t, h.Create(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"code":"validation_failed"`)
	assert.Contains(t, rec.Body.String(), `"field":"priority"`)
}

func TestHandler_Create_201BodyShape(t *testing.T) {
	t.Parallel()
	h, _ := newTestHandler(t)
	_, c, rec := newJSONRequest(http.MethodPost, "/api/todos", `{"title":"buy milk","priority":"high"}`)
	assert.NoError(t, h.Create(c))
	assert.Equal(t, http.StatusCreated, rec.Code)

	var got Todo
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.NotEmpty(t, got.ID, "server must assign an id")
	assert.Equal(t, "buy milk", got.Title)
	assert.Equal(t, PriorityHigh, got.Priority)
	assert.False(t, got.Completed)
	assert.False(t, got.CreatedAt.IsZero())
	assert.Equal(t, got.CreatedAt, got.UpdatedAt)
}

func TestHandler_List_Defaults(t *testing.T) {
	t.Parallel()
	h, s := newTestHandler(t)
	s.Add(NewTodo{Title: "first"})
	s.Add(NewTodo{Title: "second"})

	_, c, rec := newJSONRequest(http.MethodGet, "/api/todos", "")
	assert.NoError(t, h.List(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, `"title":"first"`)
	assert.Contains(t, body, `"title":"second"`)
	assert.True(t, strings.HasPrefix(strings.TrimSpace(body), "["))
}

func TestHandler_List_FilterAndSort(t *testing.T) {
	t.Parallel()
	h, s := newTestHandler(t)
	a := s.Add(NewTodo{Title: "active"})
	b := s.Add(NewTodo{Title: "done"})
	completed := true
	if _, err := s.Update(b.ID, Patch{Completed: &completed}); err != nil {
		t.Fatalf("update: %v", err)
	}
	_ = a

	_, c, rec := newJSONRequest(http.MethodGet, "/api/todos?status=active", "")
	assert.NoError(t, h.List(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, `"title":"active"`)
	assert.NotContains(t, body, `"title":"done"`)
}

func TestHandler_List_InvalidQueryParam(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		query string
	}{
		{"bad status", "status=invalid"},
		{"bad sort", "sort=garbage"},
		{"bad order", "order=sideways"},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h, _ := newTestHandler(t)
			_, c, rec := newJSONRequest(http.MethodGet, "/api/todos?"+tc.query, "")
			assert.NoError(t, h.List(c))
			assert.Equal(t, http.StatusBadRequest, rec.Code)
			assert.Contains(t, rec.Body.String(), `"code":"validation_failed"`)
		})
	}
}

func TestHandler_Update_PatchSemantics(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantInBody string
		wantField  string
	}{
		{
			name:       "completed only",
			body:       `{"completed":true}`,
			wantStatus: http.StatusOK,
			wantInBody: `"completed":true`,
		},
		{
			name:       "clear duedate",
			body:       `{"dueDate":null}`,
			wantStatus: http.StatusOK,
			wantInBody: `"dueDate":null`,
		},
		{
			name:       "set duedate",
			body:       `{"dueDate":"2026-05-15T18:00:00Z"}`,
			wantStatus: http.StatusOK,
			wantInBody: `"dueDate":"2026-05-15T18:00:00Z"`,
		},
		{
			name:       "title only",
			body:       `{"title":"updated"}`,
			wantStatus: http.StatusOK,
			wantInBody: `"title":"updated"`,
		},
		{
			name:       "priority only",
			body:       `{"priority":"high"}`,
			wantStatus: http.StatusOK,
			wantInBody: `"priority":"high"`,
		},
		{
			name:       "empty body",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
			wantInBody: `"validation_failed"`,
		},
		{
			name:       "title null rejected",
			body:       `{"title":null}`,
			wantStatus: http.StatusBadRequest,
			wantInBody: `"validation_failed"`,
			wantField:  "title",
		},
		{
			name:       "completed null rejected",
			body:       `{"completed":null}`,
			wantStatus: http.StatusBadRequest,
			wantInBody: `"validation_failed"`,
			wantField:  "completed",
		},
		{
			name:       "priority null rejected",
			body:       `{"priority":null}`,
			wantStatus: http.StatusBadRequest,
			wantInBody: `"validation_failed"`,
			wantField:  "priority",
		},
		{
			name:       "invalid priority",
			body:       `{"priority":"urgent"}`,
			wantStatus: http.StatusBadRequest,
			wantInBody: `"validation_failed"`,
			wantField:  "priority",
		},
		{
			name:       "bad duedate string",
			body:       `{"dueDate":"not-a-date"}`,
			wantStatus: http.StatusBadRequest,
			wantInBody: `"validation_failed"`,
			wantField:  "dueDate",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			h, s := newTestHandler(t)
			added := s.Add(NewTodo{Title: "x"})

			e := echo.New()
			req := httptest.NewRequest(http.MethodPatch, "/api/todos/"+added.ID, strings.NewReader(tc.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(added.ID)

			assert.NoError(t, h.Update(c))
			assert.Equal(t, tc.wantStatus, rec.Code, "body: %s", rec.Body.String())
			assert.Contains(t, rec.Body.String(), tc.wantInBody)
			if tc.wantField != "" {
				assert.Contains(t, rec.Body.String(), `"field":"`+tc.wantField+`"`)
			}
		})
	}
}

func TestHandler_Update_NotFound(t *testing.T) {
	t.Parallel()
	h, _ := newTestHandler(t)
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/todos/nope", strings.NewReader(`{"completed":true}`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("nope")

	assert.NoError(t, h.Update(c))
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), `"code":"not_found"`)
}

func TestHandler_Update_InvalidJSON(t *testing.T) {
	t.Parallel()
	h, s := newTestHandler(t)
	added := s.Add(NewTodo{Title: "x"})
	e := echo.New()
	req := httptest.NewRequest(http.MethodPatch, "/api/todos/"+added.ID, strings.NewReader(`{not json`))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(added.ID)

	assert.NoError(t, h.Update(c))
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), `"code":"invalid_json"`)
}
