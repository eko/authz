package validator

import (
	"log/slog"

	"github.com/go-playground/validator/v10"
)

func Register(
	logger *slog.Logger,
	validate *validator.Validate,
) {
	if err := validate.RegisterValidation("slug", ValidateSlug); err != nil {
		logger.Error("Unable to register slug validation rule", slog.Any("error", err))
	}
}
