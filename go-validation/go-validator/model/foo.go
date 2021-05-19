package model

type Foo struct {
	A int `validate:"required"`
	B int `validate:"rfe=A:1"` //A의 값이 1인지 체크
}
