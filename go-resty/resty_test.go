package go_resty

import (
	"context"
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

func Test_Context(t *testing.T) {
	server := mockHttpServer(t)
	defer server.Close()

	client := resty.New()
	client.SetBaseURL(server.URL)
	client.OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
		ctx := request.Context()
		value := ctx.Value("foo").(string)
		if value == "bar" {
			request.SetHeader("foo", "bar")
		}
		return nil
	})

	ctx := context.WithValue(context.Background(), "foo", "bar")
	resp, err := client.R().SetContext(ctx).Get("/employees")

	assert.NoError(t, err)
	assert.Contains(t, string(resp.Body()), "message")
	assert.Equal(t, "bar", resp.Request.Header.Get("foo"))

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

func Test_400Error(t *testing.T) {
	server := mockHttpServer(t)
	defer server.Close()

	client := resty.New()
	fmt.Println("server.URL", server.URL)
	client.SetBaseURL(server.URL)

	resp, err := client.R().Get("/bad-request")
	assert.NoError(t, err) // status 값이 BadRequest, InternalServerError이더라도 resty에서는 err를 반환하지 않음 (네트워크 오류가 발생할 때만 err를 반환하자는게 아닌가 싶음)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
}

func mockHttpServer(t *testing.T) *httptest.Server {
	// Create a mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set the response status code and headers
		w.Header().Set("Content-Type", "application/json")

		fmt.Printf("request: %+v\n", r)

		switch r.URL.Path {
		case "/employees":
			if r.Method == http.MethodGet {
				// Write the response body
				response := `{"message": "Hello, world!"}`
				w.Write([]byte(response))
			}
			w.WriteHeader(http.StatusOK)
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
			w.WriteHeader(http.StatusOK)
		case "/bad-request":
			if r.Method == http.MethodGet {
				w.WriteHeader(http.StatusBadRequest)
			}
		}

	}))

	return server
}
