package go_validator

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/kenshin579/tutorials-go/go-validation/go-validator/model"
	"github.com/stretchr/testify/assert"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	"github.com/go-playground/validator/v10"
)

func TestValidation_Oneof_Tag(t *testing.T) {
	v := validator.New()
	a := model.ArticleWithOneofTag{
		Title:      "title1",
		Body:       "body1",
		PostStatus: "TestStatus",
	}

	err := v.Struct(a)
	assert.Error(t, err)

	for _, e := range err.(validator.ValidationErrors) {
		fmt.Println(e)
	}
}

func TestValidation_Email_Validate(t *testing.T) {
	v := validator.New()
	a := model.User{
		Email: "a",
	}
	err := v.Struct(a)

	for _, e := range err.(validator.ValidationErrors) {
		fmt.Println(e)
	}
}

func TestValidation_Custom_Validator_with_Custom_Tag(t *testing.T) {
	v := validator.New()

	//custom tag에 대한 validate 함수를 등록한다
	_ = v.RegisterValidation("passwd", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > 6
	})

	a := model.CustomUser{
		Email:    "something@gmail.com",
		Name:     "A girl has no name",
		Password: "1234",
	}
	err := v.Struct(a)

	for _, e := range err.(validator.ValidationErrors) {
		fmt.Println(e)
	}
}

func TestValidation_Custom_Message(t *testing.T) {
	translator := en.New()
	uni := ut.New(translator, translator)

	// this is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	//validate tag에 대해서 사용자 정의 메시지를 등록할 수 있다
	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("passwd", trans, func(ut ut.Translator) error {
		//return ut.Add("passwd", "{0} is not strong enough", true) // see universal-translator for details
		return ut.Add("passwd", "{0} is not strong enough", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("passwd", fe.Field())
		return t
	})

	_ = v.RegisterValidation("passwd", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > 6
	})

	a := model.CustomUser{
		Email:    "a",
		Password: "1234",
	}
	err := v.Struct(a)

	for _, e := range err.(validator.ValidationErrors) {
		fmt.Println(e.Translate(trans))
	}
}

//https://github.com/go-playground/validator/issues/494
func TestValidation_Check_BasedOn_Another_Field_Value(t *testing.T) {
	foo := model.Foo{
		A: 1,
		B: 0,
	}

	v := validator.New()

	// register custom validation: rfe(Required if Field is Equal to some value).
	v.RegisterValidation(`rfe`, func(fl validator.FieldLevel) bool {
		param := strings.Split(fl.Param(), `:`) //A:1
		paramField := param[0]                  //A
		paramValue := param[1]                  //1

		if paramField == `` {
			return true
		}

		// param field reflect.Value.
		var paramFieldValue reflect.Value

		//fl.Parent()는 변수 Foo.A를 가리킨다
		if fl.Parent().Kind() == reflect.Ptr {
			paramFieldValue = fl.Parent().Elem().FieldByName(paramField)
		} else {
			paramFieldValue = fl.Parent().FieldByName(paramField)
		}

		//A의 값과 paramValue를 비교하는 로직
		if isEq(paramFieldValue, paramValue) == false {
			return true
		}

		return hasValue(fl)
	})

	err := v.Struct(&foo)
	assert.Error(t, err)
	assert.Equal(t, "Key: 'Foo.B' Error:Field validation for 'B' failed on the 'rfe' tag", err.Error())
}

func hasValue(fl validator.FieldLevel) bool {
	return requireCheckFieldKind(fl, "")
}

//todo: 다음에 다시 스터디하는 걸로 함
//필드 값이 존재하는지 체크하는 로직으로 판단됨
func requireCheckFieldKind(fl validator.FieldLevel, param string) bool {
	field := fl.Field()
	if len(param) > 0 {
		if fl.Parent().Kind() == reflect.Ptr {
			field = fl.Parent().Elem().FieldByName(param)
		} else {
			field = fl.Parent().FieldByName(param)
		}
	}
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		_, _, nullable := fl.ExtractType(field)
		if nullable && field.Interface() != nil {
			return true
		}
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

func isEq(field reflect.Value, value string) bool {
	switch field.Kind() {

	case reflect.String:
		return field.String() == value
	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(value)
		return int64(field.Len()) == p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(value)
		return field.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(value)
		return field.Uint() == p

	case reflect.Float32, reflect.Float64:
		p := asFloat(value)
		return field.Float() == p
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func asInt(param string) int64 {
	i, err := strconv.ParseInt(param, 0, 64)
	panicIf(err)

	return i
}

func asUint(param string) uint64 {
	i, err := strconv.ParseUint(param, 0, 64)
	panicIf(err)

	return i
}

func asFloat(param string) float64 {
	i, err := strconv.ParseFloat(param, 64)
	panicIf(err)

	return i
}

func panicIf(err error) {
	if err != nil {
		panic(err.Error())
	}
}
