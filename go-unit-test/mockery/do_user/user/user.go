package user

import "github.com/kenshin579/tutorials-go/go-unit-test/mockery/do_user/doer"

type User struct {
	Doer doer.Doer
}

func (u *User) Use() {
	u.Doer.Do(1, "abc")
}
