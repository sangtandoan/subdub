package validator

import "github.com/go-playground/validator/v10"

type Validator interface {
	Validate(object any) error
}

type appValidator struct {
	v *validator.Validate
}

func NewAppValidator() *appValidator {
	return &appValidator{
		v: validator.New(),
	}
}

func (av *appValidator) Validate(object any) error {
	return av.v.Struct(object)
}
