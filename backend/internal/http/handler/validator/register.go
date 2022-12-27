package validator

import (
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

func Register(
	logger *slog.Logger,
	validate *validator.Validate,
) {
	if err := validate.RegisterValidation("slug", ValidateSlug); err != nil {
		logger.Error("Unable to register slug validation rule", err)
	}
}
