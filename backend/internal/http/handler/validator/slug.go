package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	slugReplaceCharacters = regexp.MustCompile(`[^a-z0-9-_\./*]+`)
)

func ValidateSlug(fl validator.FieldLevel) bool {
	return ValidateSlugFromString(fl.Field().String())
}

func ValidateSlugFromString(value string) bool {
	compared := slugReplaceCharacters.ReplaceAllString(value, "-")

	return value == compared
}
