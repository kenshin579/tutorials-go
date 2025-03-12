package job

import "fmt"

type PanicJob struct {
	Count int
}

func (p *PanicJob) Run() {
	p.Count++
	if p.Count == 1 {
		panic("oooooooooooooops!!!")
	}

	fmt.Println("hello world")
}
