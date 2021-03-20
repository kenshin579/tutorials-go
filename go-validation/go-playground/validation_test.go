package go_playground

import (
	"fmt"
	"log"
	"testing"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required"`
}

type CustomUser struct {
	Email    string `json:"email" validate:"required,email"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"passwd"` // <-- a custom validation rule
}

func TestValidation_Email(t *testing.T) {
	v := validator.New()
	a := User{
		Email: "a",
	}
	err := v.Struct(a)

	for _, e := range err.(validator.ValidationErrors) {
		fmt.Println(e)
	}
}

func TestValidation_Custom_Validator(t *testing.T) {
	v := validator.New()
	_ = v.RegisterValidation("passwd", func(fl validator.FieldLevel) bool {
		return len(fl.Field().String()) > 6
	})

	a := CustomUser{
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

	a := CustomUser{
		Email:    "a",
		Password: "1234",
	}
	err := v.Struct(a)

	for _, e := range err.(validator.ValidationErrors) {
		fmt.Println(e.Translate(trans))
	}
}
