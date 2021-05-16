package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo"
)

func HTTPPostMethod(url string, request interface{}) ([]byte, error) {
	body, _ := json.Marshal(&request)
	response, err := http.Post(url, echo.MIMEApplicationJSON, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(response.Body)
}
