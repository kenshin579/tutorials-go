package http

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

type httpSendStore struct {
	client *resty.Client
}

func NewHttpSendStore(client *resty.Client) *httpSendStore {
	return &httpSendStore{
		client: client,
	}
}

func (h *httpSendStore) Send(url string) ([]byte, error) {
	resp, err := h.client.R().Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return resp.Body(), nil
}
