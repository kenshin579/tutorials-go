package go_resty

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func Test_Get(t *testing.T) {
	server := mockHttpServer(t)
	defer server.Close()

	client := resty.New()
	resp, err := client.R().Get(server.URL)
	assert.NoError(t, err)
	assert.Contains(t, string(resp.Body()), "message")
}

type User struct {
	Name string `json:"name"`
}

func Test_Post(t *testing.T) {
	server := mockHttpServer(t)
	defer server.Close()

	user := User{
		Name: "frank",
	}

	jsonUser, _ := json.Marshal(user)

	client := resty.New()
	resp, err := client.R().
		SetBody(jsonUser).
		Post(server.URL)
	assert.NoError(t, err)
	assert.Contains(t, string(resp.Body()), "message")
}

func Test_Execute(t *testing.T) {
	server := mockHttpServer(t)
	defer server.Close()

	user := User{
		Name: "frank",
	}

	jsonUser, _ := json.Marshal(user)

	client := resty.New()
	resp, err := client.R().
		SetBody(jsonUser).
		Execute(resty.MethodPost, server.URL)
	assert.NoError(t, err)
	assert.Contains(t, string(resp.Body()), "message")
}

func mockHttpServer(t *testing.T) *httptest.Server {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the response status code and headers
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {
			// Write the response body
			response := `{"message": "Hello, world!"}`
			w.Write([]byte(response))
		} else if r.Method == http.MethodPost {
			// Read the request body
			requestBody, err := ioutil.ReadAll(r.Body)
			fmt.Println("requestBody", string(requestBody))

			if err != nil {
				t.Fatalf("Failed to read request body: %v", err)
			}

			// Parse the request body as JSON
			var requestBodyData map[string]interface{}
			err = json.Unmarshal(requestBody, &requestBodyData)
			if err != nil {
				t.Fatalf("Failed to parse request body as JSON: %v", err)
			}

			// Write the response body
			response := `{"message": "User created successfully"}`
			w.Write([]byte(response))
		}
	}))

	return server
}
