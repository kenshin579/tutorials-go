package resty

import (
	"fmt"
	"testing"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/labstack/echo/v4"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-resty/resty/v2"
)

const (
	url = "http://localhost:8080"
)

//https://github.com/go-resty/resty/blob/master/retry_test.go
func Test_Resty_Timeout(t *testing.T) {
	client := resty.New().SetTimeout(1*time.Second).
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
		OnError(func(request *resty.Request, err error) {
			log.Errorf("request:%v, err:%v\n", request, err)
		})
	resp, err := client.R().Get(url)

	fmt.Println("err", err)
	fmt.Println("resp", resp)
	assert.Error(t, err)
}

func Test_Resty_Success(t *testing.T) {
	client := resty.New().SetTimeout(5*time.Second).
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp, err := client.R().Get(url)
	require.NoError(t, err)
	assert.Equal(t, "Hello World", string(resp.Body()))
}

func Test_Resty_Timeout_Retry(t *testing.T) {
	attempts := 3
	client := resty.New().SetTimeout(1*time.Second).
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
		SetRetryCount(attempts)
	resp, err := client.R().Get(url)

	fmt.Println("err", err)
	fmt.Println("resp", resp)
	assert.Error(t, err)
}

func Test_Resty_Middleware(t *testing.T) {
	client := resty.New().SetTimeout(5*time.Second).
		SetHeader(echo.HeaderContentType, echo.MIMEApplicationJSON).
		OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
			log.Infof("request: %+v\n", request)
			return nil
		}).
		OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
			log.Infof("response: %+v\n", response)
			return nil
		})

	resp, err := client.R().Get(url)
	require.NoError(t, err)
	assert.Equal(t, "Hello World", string(resp.Body()))
}
