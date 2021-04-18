package model

//https://golangbyexample.com/error-in-golang/
type InputError struct {
	Message      string
	MissingField string
}

//inputError가 Error() 메서드 구현하고 있어서 type error로 인식된다
func (i *InputError) Error() string {
	return i.Message
}

//추가 정보를 얻어오는 메서드도 추가할 수 있다
func (i *InputError) GetMissingField() string {
	return i.MissingField
}
