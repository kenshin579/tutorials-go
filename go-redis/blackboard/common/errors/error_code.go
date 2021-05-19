package errors

import (
	"net/http"
)

var (
	/*
		blackboard : 10000
		system   : 50000
	*/
	ErrNotFoundKey = NewBlackBoardError(http.StatusNotFound, 10100, "key is not found : {0}")

	ErrInvalidRequest = NewBlackBoardError(http.StatusBadRequest, 50100, "request is invalid")
	ErrBinding        = NewBlackBoardError(http.StatusBadRequest, 50200, "request binding error")
)
