package go_testing

import (
	"net/mail"
	"testing"
)

type testEmailSender struct {
	lastSubject string
	lastBody    string
	lastTo      []*mail.Address //todo : 이건 어떻게 사용할 수 있나?
}

// make sure it satisfies the interface
var _ EmailSender = (*testEmailSender)(nil)

func (t *testEmailSender) Send(subject, body string, to ...*mail.Address) {
	t.lastSubject = subject
	t.lastBody = body
	t.lastTo = to
}

//todo : 이 부분 다시 확인해보기
func TestSendWelcomeEmail(t *testing.T) {
	sender := &testEmailSender{
		lastSubject: "subject",
		lastBody:    "body",
		//lastTo: &[]mail.Address{"name1", "test3@test.com"},

	}
	to1 := &mail.Address{Name: "Frank", Address: "test1@test.com"}
	to2 := &mail.Address{Name: "Frank", Address: "test2@test.com"}

	SendWelcomeEmail(sender, to1, to2)
	if sender.lastSubject != "Welcome" {
		t.Error("Subject line was wrong")
	}

	//if sender.To[0] != to1 && sender.To[1] != to2 {
	//	t.Error("Wrong recipients")
	//}
}
