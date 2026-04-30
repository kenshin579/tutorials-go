package todo

import (
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

// writeError maps domain errors to JSON error responses with appropriate HTTP status codes.
// *ValidationError → 400 validation_failed (with field detail), ErrNotFound → 404 not_found,
// all else → 500 internal_error (logged).
func writeError(c echo.Context, err error) error {
	var verr *ValidationError
	switch {
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
