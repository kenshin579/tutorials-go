package go_resty

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
)

func Test_Get(t *testing.T) {
	server := mockHttpServer(t)
	defer server.Close()

	client := resty.New()
	resp, err := client.R().Get(server.URL + "/employees")
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
		Post(server.URL + "/users")
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

func Test_BaseUrl(t *testing.T) {
	server := mockHttpServer(t)
	defer server.Close()

	client := resty.New()
	fmt.Println("server.URL", server.URL)
	client.SetBaseURL(server.URL)

	resp, err := client.R().Get("/employees")
	assert.NoError(t, err)
	assert.Contains(t, string(resp.Body()), "message")
}

func Test_Middleware(t *testing.T) {
	server := mockHttpServer(t)
	defer server.Close()

	client := resty.New()
	client.SetBaseURL(server.URL)
	client.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		log.Infof("resty.OnBeforeRequest(). url:%s, body:%s", request.URL, request.Body)
		return nil
	})

	client.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		log.Infof("resty.OnAfterResponse(). url:%s, body:%s", response.Request.URL, response.Body())
		return nil
	})

	resp, err := client.R().Get("/employees")
	assert.NoError(t, err)
	assert.Contains(t, string(resp.Body()), "message")
}

func mockHttpServer(t *testing.T) *httptest.Server {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the response status code and headers
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		fmt.Printf("request: %+v\n", r)

		switch r.URL.Path {
		case "/employees":
			if r.Method == http.MethodGet {
				// Write the response body
				response := `{"message": "Hello, world!"}`
				w.Write([]byte(response))
			}
		case "/users":
			if r.Method == http.MethodPost {
				// Read the request body
				requestBody, err := io.ReadAll(r.Body)
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
		}

	}))

	return server
}
