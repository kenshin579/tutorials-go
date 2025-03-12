package router

import (
	"github.com/kenshin579/tutorials-go/go-validation/article/model"
	"github.com/kenshin579/tutorials-go/go-validation/article/utils"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *Validator {
	v := validator.New()

	v.RegisterValidation(`postStatus`, ValidatePostStatus)

	return &Validator{
		validator: v,
	}
}

type Validator struct {
	validator *validator.Validate
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}

func ValidatePostStatus(fl validator.FieldLevel) bool {
	return utils.Include(model.AllPostStatus, fl.Field().Interface())
}
