package job

import "log"

type Hello struct {
	Name string
}

func (h Hello) Run() {
	log.Println("Hi " + h.Name)
}
