package task

import "log"

type Print struct {
	Message string
}

type PrintRequest struct {
	Message string `json:"message"`
}

func (p *Print) Run() {
	log.Println(p.Message)
}
