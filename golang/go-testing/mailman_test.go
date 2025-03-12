package go_testing

import (
	"net/mail"
	"testing"
)

// https://blog.fraixed.es/post/golang-slices-structs-or-pointers-to-structs-dilemma/
type testEmailSender struct {
	lastSubject string
	lastBody    string
	lastTo      []*mail.Address
}

// make sure it satisfies the interface
var _ EmailSender = (*testEmailSender)(nil)

func (t *testEmailSender) Send(subject, body string, to ...*mail.Address) {
	t.lastSubject = subject
	t.lastBody = body
	t.lastTo = to
}

// interface mocking하기
func TestSendWelcomeEmail(t *testing.T) {
	address := make([]*mail.Address, 1)
	address[0] = &mail.Address{
		Name:    "Sender",
		Address: "test3@test.com",
	}
	sender := &testEmailSender{
		lastSubject: "Welcome",
		lastBody:    "body",
		lastTo:      address,
	}
	to1 := &mail.Address{Name: "Receiver1", Address: "test1@test.com"}
	to2 := &mail.Address{Name: "Receiver2", Address: "test2@test.com"}

	SendWelcomeEmail(sender, to1, to2)
	if sender.lastSubject != "Welcome" {
		t.Error("Subject line was wrong")
	}

	//if sender.lastTo[0] != to1 && sender.To[1] != to2 {
	//	t.Error("Wrong recipients")
	//}
}
