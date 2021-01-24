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
	// some code
}
func New() *MailMan {
	return &MailMan{}
}

func SendWelcomeEmail(m EmailSender, to ...*mail.Address) {
	fmt.Println("SendWelcomeEmail")
}
