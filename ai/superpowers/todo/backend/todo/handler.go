package todo

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// Handler exposes the todo store as JSON HTTP endpoints.
// It is responsible for HTTP I/O concerns only — request parsing,
// validation shape, response serialization, status code mapping.
// Domain logic lives in Store.
type Handler struct {
	store *Store
}

// NewHandler returns a Handler bound to the given store.
func NewHandler(s *Store) *Handler {
	return &Handler{store: s}
}

// createBody is the JSON body shape accepted by POST /api/todos.
type createBody struct {
	Title    string     `json:"title"`
	Priority string     `json:"priority"`
	DueDate  *time.Time `json:"dueDate"`
}

// Create handles POST /api/todos.
// Body: { "title": "...", "priority"?: "...", "dueDate"?: "RFC3339" }.
// Returns 201 with the created Todo, or 400 on validation/JSON errors.
func (h *Handler) Create(c echo.Context) error {
	var body createBody
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return c.JSON(http.StatusBadRequest, errBody("invalid_json", "request body is not valid JSON", nil))
	}
	input := NewTodo{
		Title:    body.Title,
		Priority: Priority(body.Priority),
		DueDate:  body.DueDate,
	}
	if err := input.Validate(); err != nil {
		return writeError(c, err)
	}
	created := h.store.Add(input)
	return c.JSON(http.StatusCreated, created)
}

// Delete handles DELETE /api/todos/{id}.
// Returns 204 on success or 404 when the id is unknown.
func (h *Handler) Delete(c echo.Context) error {
	id := c.Param("id")
	if err := h.store.Delete(id); err != nil {
		return writeError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// errBody builds the standard error envelope { "error": { code, message, details? } }.
func errBody(code, msg string, details map[string]any) echo.Map {
	body := echo.Map{"code": code, "message": msg}
	if details != nil {
		body["details"] = details
	}
	return echo.Map{"error": body}
}

// List handles GET /api/todos with optional status/sort/order query params.
// Empty params apply defaults (status=all, sort=createdAt, order=desc).
func (h *Handler) List(c echo.Context) error {
	q, err := parseQuery(c)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, h.store.List(q))
}

// parseQuery validates and converts URL query params into a Query.
// Returns *ValidationError when any param has an invalid value.
func parseQuery(c echo.Context) (Query, error) {
	q := Query{}
	if v := c.QueryParam("status"); v != "" {
		s := StatusFilter(v)
		if !s.IsValid() {
			return Query{}, &ValidationError{Field: "status", Message: "status must be one of: all, active, completed"}
		}
		q.Status = s
	}
	if v := c.QueryParam("sort"); v != "" {
		s := SortKey(v)
		if !s.IsValid() {
			return Query{}, &ValidationError{Field: "sort", Message: "sort must be one of: createdAt, dueDate, priority"}
		}
		q.Sort = s
	}
	if v := c.QueryParam("order"); v != "" {
		o := OrderDir(v)
		if !o.IsValid() {
			return Query{}, &ValidationError{Field: "order", Message: "order must be one of: asc, desc"}
		}
		q.Order = o
	}
	return q, nil
}

// Update handles PATCH /api/todos/{id}. The body must be a non-empty JSON object.
// Only dueDate accepts JSON null (clears the value); other fields with null are rejected.
func (h *Handler) Update(c echo.Context) error {
	patch, err := parsePatch(c)
	if err != nil {
		return writeError(c, err)
	}
	id := c.Param("id")
	updated, err := h.store.Update(id, patch)
	if err != nil {
		return writeError(c, err)
	}
	return c.JSON(http.StatusOK, updated)
}

// parsePatch reads the request body as a partial Todo update.
// Uses map[string]json.RawMessage for first-pass decoding so we can
// distinguish "key absent" from "key present with null".
// Only dueDate accepts JSON null (clears the value).
func parsePatch(c echo.Context) (Patch, error) {
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(c.Request().Body).Decode(&raw); err != nil {
		return Patch{}, &jsonDecodeError{}
	}

	patch := Patch{}

	if v, ok := raw["title"]; ok {
		if isJSONNull(v) {
			return Patch{}, &ValidationError{Field: "title", Message: "title cannot be null"}
		}
		var s string
		if err := json.Unmarshal(v, &s); err != nil {
			return Patch{}, &ValidationError{Field: "title", Message: "title must be a string"}
		}
		patch.Title = &s
	}
	if v, ok := raw["completed"]; ok {
		if isJSONNull(v) {
			return Patch{}, &ValidationError{Field: "completed", Message: "completed cannot be null"}
		}
		var b bool
		if err := json.Unmarshal(v, &b); err != nil {
			return Patch{}, &ValidationError{Field: "completed", Message: "completed must be a boolean"}
		}
		patch.Completed = &b
	}
	if v, ok := raw["priority"]; ok {
		if isJSONNull(v) {
			return Patch{}, &ValidationError{Field: "priority", Message: "priority cannot be null"}
		}
		var s string
		if err := json.Unmarshal(v, &s); err != nil {
			return Patch{}, &ValidationError{Field: "priority", Message: "priority must be a string"}
		}
		p := Priority(s)
		patch.Priority = &p
	}
	if v, ok := raw["dueDate"]; ok {
		if isJSONNull(v) {
			patch.ClearDueDate = true
		} else {
			var t time.Time
			if err := json.Unmarshal(v, &t); err != nil {
				return Patch{}, &ValidationError{Field: "dueDate", Message: "dueDate must be RFC3339 timestamp or null"}
			}
			patch.DueDate = &t
		}
	}
	return patch, nil
}

// isJSONNull reports whether the raw message is the literal "null" (modulo whitespace).
func isJSONNull(b json.RawMessage) bool {
	return string(bytes.TrimSpace(b)) == "null"
}

// jsonDecodeError is a sentinel-style type wrapping JSON parse failures
// so writeError emits invalid_json instead of validation_failed.
type jsonDecodeError struct{}

func (jsonDecodeError) Error() string { return "invalid_json" }

// writeError maps domain errors to JSON error responses with appropriate HTTP status codes.
// *jsonDecodeError → 400 invalid_json, *ValidationError → 400 validation_failed (with field detail),
// ErrNotFound → 404 not_found, all else → 500 internal_error (logged).
func writeError(c echo.Context, err error) error {
	var verr *ValidationError
	var jerr *jsonDecodeError
	switch {
	case errors.As(err, &jerr):
		return c.JSON(http.StatusBadRequest, errBody("invalid_json", "request body is not valid JSON", nil))
	case errors.As(err, &verr):
		var details map[string]any
		if verr.Field != "" {
			details = map[string]any{"field": verr.Field}
		}
		return c.JSON(http.StatusBadRequest, errBody("validation_failed", verr.Message, details))
	case errors.Is(err, ErrNotFound):
		return c.JSON(http.StatusNotFound, errBody("not_found", err.Error(), nil))
	default:
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, errBody("internal_error", "internal server error", nil))
	}
}
