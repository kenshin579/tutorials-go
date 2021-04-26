package router

import (
	"fmt"
	"regexp"

	"github.com/kenshin579/tutorials-go/go-validation/article/exception"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo"
)

type Error struct {
	Code    int                    `json:"code,omitempty"`
	Message interface{}            `json:"message"`
	Errors  map[string]interface{} `json:"errors,omitempty"`
}

func customHTTPErrorHandler(err error, c echo.Context) {
	e := Error{}
	var httpCode int
	switch v := err.(type) {
	case *echo.HTTPError:
		fmt.Println("case HTTPError!!!", v)
		e.Message = v.Message
		httpCode = e.Code
	case validator.ValidationErrors:
		fmt.Printf("case ValidationErrors!!! : %v\n", v)
		e.Errors = make(map[string]interface{})

		for _, ev := range v {
			e.Errors["field"] = ev.StructNamespace()
			e.Errors["message"] = extractErrorMessage(ev.Error())
			break
		}
		e.Code = exception.ErrInvalidRequest.ErrorCode()
		e.Message = exception.ErrInvalidRequest.ErrorMessage()
		httpCode = exception.ErrInvalidRequest.HttpCode()
	case *exception.CustomError:
		fmt.Printf("case CustomError!!! : %v\n", v)
		e.Code = v.ErrorCode()
		e.Message = v.Error()
		httpCode = v.HttpCode()
	default:
		fmt.Printf("not predefined error type : %v\n", v)
	}

	c.Logger().Error(e)
	c.JSON(httpCode, e)
}

func extractErrorMessage(msg string) string {
	r := regexp.MustCompile("Error:(.*)")

	matchParts := r.FindStringSubmatch(msg)
	if len(matchParts) > 0 {
		return matchParts[0]
	}
	return ""
}
