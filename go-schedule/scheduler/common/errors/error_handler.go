package errors

import (
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo"
)

type Error struct {
	Code    int                    `json:"code,omitempty"`
	Message interface{}            `json:"message"`
	Errors  map[string]interface{} `json:"errors,omitempty"`
}

func ScheduleHTTPErrorHandler(err error, c echo.Context) {
	e := Error{}
	var httpCode int
	switch v := err.(type) {
	case *echo.HTTPError:
		e.Message = v.Message
		httpCode = v.Code
	case validator.ValidationErrors:
		e.Errors = make(map[string]interface{})

		for _, ev := range v {
			e.Errors["field"] = ev.StructNamespace()
			e.Errors["message"] = extractErrorMessagePart(ev.Error())
		}

		e.Code = ErrInvalidRequest.ErrorCode()
		e.Message = ErrInvalidRequest.Error()
		httpCode = ErrInvalidRequest.HttpCode()
	case *ScheduleError:
		e.Code = v.ErrorCode()
		e.Message = v.Error()
		httpCode = v.HttpCode()
	default:
		e.Message = v.Error()
		httpCode = http.StatusInternalServerError
	}

	c.Logger().Error(e)
	c.JSON(httpCode, e)
}

func extractErrorMessagePart(msg string) string {
	r := regexp.MustCompile("Error:(.*)")

	matchParts := r.FindStringSubmatch(msg)
	if len(matchParts) > 0 {
		return matchParts[0]
	}
	return ""
}
