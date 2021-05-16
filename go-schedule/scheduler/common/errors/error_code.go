package errors

import (
	"net/http"
)

var (
	/*
		schedule : 10000
		system   : 50000
	*/
	ErrNotFoundJob = NewScheduleError(http.StatusNotFound, 10100, "job is not found : {0}")

	ErrInvalidRequest = NewScheduleError(http.StatusBadRequest, 50100, "request is invalid")
	ErrBinding        = NewScheduleError(http.StatusBadRequest, 50200, "request binding error")
)
