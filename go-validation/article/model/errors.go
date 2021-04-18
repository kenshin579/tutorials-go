package model

type HTTPError struct {
	ErrorCode int         `json:"-"`
	Message   interface{} `json:"message"`
}

const ()
