package go_testing

import (
	"fmt"
	"net/mail"
)

type MailMan struct{}

type EmailSender interface {
	Send(subject, body string, to ...*mail.Address)
}

func (m *MailMan) Send(subject, body string, to ...*mail.Address) {
	fmt.Println("sending mail to -> ", to[0])
}
func New() *MailMan {
	return &MailMan{}
}

// func SendWelcomeEmail(m *MailMan, to ...*mail.Address) //인터페이스를 받도록 수정함
func SendWelcomeEmail(m EmailSender, to ...*mail.Address) {
	fmt.Println("m", m)
}
