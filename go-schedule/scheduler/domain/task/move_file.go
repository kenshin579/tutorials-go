package task

import (
	"encoding/json"
	"fmt"

	"github.com/kenshin579/tutorials-go/go-schedule/scheduler/common/utils"

	"github.com/labstack/gommon/log"
)

type HttpMoveFile struct {
	Url  string          `json:"url"`
	Body MoveFileRequest `json:"body"`
}

type MoveFileRequest struct {
	FileName    string `json:"fileName"`
	Destination string `json:"destination"`
}

type MoveFileResponse struct {
	Status string `json:"status"`
}

func (h *HttpMoveFile) Run() {
	log.Info(fmt.Sprintf("moving file %s -> %s", h.Body.FileName, h.Body.Destination))

	response, err := h.moveFile(h.Url, h.Body)
	if err != nil {
		log.Error(err)
	}

	log.Debug(response)
}

func (h *HttpMoveFile) moveFile(url string, body MoveFileRequest) (MoveFileResponse, error) {
	result := MoveFileResponse{}

	response, err := utils.HTTPPostMethod(url, body)
	if err != nil {
		log.Error(err)
		return result, err
	}
	err = json.Unmarshal(response, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
