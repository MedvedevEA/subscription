package validator

import (
	"github.com/go-playground/validator"
)

var global *validator.Validate

func init() {
	SetValidator(New())
}

func New() *validator.Validate {
	v := validator.New()
	return v
}
func SetValidator(v *validator.Validate) {
	global = v
}
func Validator() *validator.Validate {
	return global
}
func Validate(structure any) error {
	return Validator().Struct(structure)
}
