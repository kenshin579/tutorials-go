package utils

import (
	"fmt"

	"github.com/labstack/gommon/log"

	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

func NewError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	switch v := err.(type) {
	case *echo.HTTPError:
		e.Errors["body"] = v.Message
	default:
		e.Errors["body"] = v.Error()
	}
	return e
}

func NewValidatorError(err error) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	errs := err.(validator.ValidationErrors)

	for _, v := range errs {
		/*
			{
			  "time": "2021-04-18T18:20:23.050955+09:00",
			  "level": "ERROR",
			  "prefix": "-",
			  "file": "errors.go",
			  "line": "36",
			  "message": "vKey: 'ArticleRequest.Description' Error:Field validation for 'Description' failed on the 'required' tag"
			}
		*/
		log.Error("v", v)
		e.Errors[v.Field()] = fmt.Sprintf("%v", v.Tag())
	}
	return e
}

func NotFound() Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["body"] = "resource not found"
	return e
}
