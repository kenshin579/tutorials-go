package errors

import (
	"fmt"
	"regexp"
)

type ScheduleError struct {
	httpCode     int
	errorCode    int
	errorMessage string
}

func (se *ScheduleError) Error() string {
	return se.errorMessage
}

func (se *ScheduleError) ErrorCode() int {
	return se.errorCode
}

func (se *ScheduleError) HttpCode() int {
	return se.httpCode
}

func NewScheduleError(httpCode int, errorCode int, errorMessage string) *ScheduleError {
	return &ScheduleError{
		httpCode:     httpCode,
		errorCode:    errorCode,
		errorMessage: errorMessage,
	}
}

func (se *ScheduleError) WithParams(params ...string) *ScheduleError {
	var newMessage = se.Error()
	for i, param := range params {
		re := regexp.MustCompile(fmt.Sprintf("\\{%d\\}", i))
		newMessage = re.ReplaceAllString(newMessage, param)
	}
	se.errorMessage = newMessage
	return se
}
