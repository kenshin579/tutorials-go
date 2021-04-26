package exception

type CustomError struct {
	httpCode     int
	errorCode    int
	errorMessage string
}

func (ce *CustomError) Error() string {
	return ce.errorMessage
}

func (ce *CustomError) ErrorCode() int {
	return ce.errorCode
}

func (ce *CustomError) HttpCode() int {
	return ce.httpCode
}

func (ce *CustomError) ErrorMessage() string {
	return ce.errorMessage
}

func NewCustomError(httpCode int, errorCode int, errorMessage string) *CustomError {
	return &CustomError{
		httpCode:     httpCode,
		errorCode:    errorCode,
		errorMessage: errorMessage,
	}
}
