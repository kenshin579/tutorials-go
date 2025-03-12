package model

import "strconv"

type CustomError struct {
	httpCode     int
	errorCode    int
	errorMessage string
}

// Create a function Error() string and associate it to the struct.
func (error *CustomError) Error() string {
	return strconv.Itoa(error.httpCode) + ":" + strconv.Itoa(error.errorCode) + ":" + error.errorMessage
}

// Then create an error object using MyError struct.
func NewCustomError(httpCode int, errorCode int, errorMessage string) error {
	return &CustomError{
		httpCode:     httpCode,
		errorCode:    errorCode,
		errorMessage: errorMessage,
	}
}
