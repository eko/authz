package handler

import (
	"net/http"

	"github.com/eko/authz/backend/internal/manager"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CheckQuery struct {
	Principal     string `json:"principal" validate:"required,slug"`
	ResourceKind  string `json:"resource_kind" validate:"required,slug"`
	ResourceValue string `json:"resource_value" validate:"required,slug"`
	Action        string `json:"action" validate:"required,slug"`
	IsAllowed     bool   `json:"is_allowed"`
}

type CheckRequest struct {
	Checks []*CheckQuery `json:"checks" validate:"required,dive"`
}

type CheckResponse struct {
	Checks []*CheckQuery `json:"checks"`
}

// Check if a principal has access to do action on resource.
//
//	@security	Authentication
//	@Summary	Check if a principal has access to do action on resource
//	@Tags		Check
//	@Produce	json
//	@Param		default	body		CheckRequest	true	"Check request"
//	@Success	200		{object}	CheckResponse
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/policies [Post]
func Check(
	validate *validator.Validate,
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := &CheckRequest{}

		// Parse request body
		if err := c.BodyParser(request); err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Validate body
		if err := validateStruct(validate, request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		// Create policy
		for _, check := range request.Checks {
			isAllowed, err := manager.IsAllowed(check.Principal, check.ResourceKind, check.ResourceValue, check.Action)
			if err != nil {
				return returnError(c, http.StatusInternalServerError, err)
			}

			check.IsAllowed = isAllowed
		}

		return c.JSON(&CheckResponse{
			Checks: request.Checks,
		})
	}
}
