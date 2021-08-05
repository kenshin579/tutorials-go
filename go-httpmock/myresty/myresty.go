package myresty

import (
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	Timeout    = 15000
	RetryCount = 3
)

func New() *resty.Client {
	return resty.New().
		SetTimeout(Timeout * time.Millisecond).
		SetRetryCount(RetryCount - 1)
}
