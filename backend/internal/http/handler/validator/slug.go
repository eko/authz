package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	slugReplaceCharacters = regexp.MustCompile(`[^a-z0-9-_\./*]+`)
)

func ValidateSlug(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	compared := slugReplaceCharacters.ReplaceAllString(value, "-")

	return value == compared
}
