package validator

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator is a wrapper for the validator.Validate
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator creates a new validator
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

// Validate validates a struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
