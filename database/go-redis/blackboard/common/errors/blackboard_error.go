package errors

import (
	"fmt"
	"regexp"
)

type BlackBoardError struct {
	httpCode     int
	errorCode    int
	errorMessage string
}

func (be *BlackBoardError) Error() string {
	return be.errorMessage
}

func (be *BlackBoardError) ErrorCode() int {
	return be.errorCode
}

func (be *BlackBoardError) HttpCode() int {
	return be.httpCode
}

func NewBlackBoardError(httpCode int, errorCode int, errorMessage string) *BlackBoardError {
	return &BlackBoardError{
		httpCode:     httpCode,
		errorCode:    errorCode,
		errorMessage: errorMessage,
	}
}

func (be *BlackBoardError) WithParams(params ...string) *BlackBoardError {
	var newMessage = be.Error()
	for i, param := range params {
		re := regexp.MustCompile(fmt.Sprintf("\\{%d\\}", i))
		newMessage = re.ReplaceAllString(newMessage, param)
	}
	be.errorMessage = newMessage
	return be
}
