package model

import "net/http"

var (
	/*
	   account : 10000
	   user    : 20000
	   system  : 30000
	*/

	ErrRequestUser  = NewCustomError(http.StatusBadRequest, 10100, "Request is invalid")
	ErrUserNotFound = NewCustomError(http.StatusNotFound, 20100, "User is not found")
)
