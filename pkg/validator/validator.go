package validator

import (
	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

func (cv *CustomValidator) ValidateStruct(s interface{}) error {
	if err := cv.validator.Struct(s); err != nil {
		return err
	}
	return nil
}
