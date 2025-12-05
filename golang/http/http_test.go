package go_http2

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

const (
	url = "http://localhost:8080"
)

func Test_Http_Client_With_Timeout(t *testing.T) {
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	response, err := client.Get(url)
	assert.Empty(t, response)
	fmt.Println("err", err)
	assert.Error(t, err)
}

func Test_Http_Client(t *testing.T) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	response, err := client.Get(url)
	require.NoError(t, err)
	fmt.Println("response", response)

	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	require.NoError(t, err)
	assert.Equal(t, "Hello World", string(body))
}
