package core

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	*validator.Validate
}

func NewValidator() *Validator {
	return &Validator{
		validator.New(),
	}
}
