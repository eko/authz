package handler

import (
	"net/http"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type CheckRequestQuery struct {
	Principal     string `json:"principal" validate:"required,slug"`
	ResourceKind  string `json:"resource_kind" validate:"required,slug"`
	ResourceValue string `json:"resource_value" validate:"required,slug"`
	Action        string `json:"action" validate:"required,slug"`
}

type CheckResponseQuery struct {
	*CheckRequestQuery
	IsAllowed bool `json:"is_allowed"`
}

type CheckRequest struct {
	Checks []*CheckRequestQuery `json:"checks" validate:"required,dive"`
}

type CheckResponse struct {
	Checks []*CheckResponseQuery `json:"checks"`
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
//	@Router		/v1/check [Post]
func Check(
	validate *validator.Validate,
	compiledManager manager.CompiledPolicy,
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
		var responseChecks = make([]*CheckResponseQuery, len(request.Checks))

		for i, check := range request.Checks {
			isAllowed, err := compiledManager.IsAllowed(check.Principal, check.ResourceKind, check.ResourceValue, check.Action)
			if err != nil {
				return returnError(c, http.StatusInternalServerError, err)
			}

			responseChecks[i] = &CheckResponseQuery{
				CheckRequestQuery: check,
				IsAllowed:         isAllowed,
			}
		}

		return c.JSON(&CheckResponse{
			Checks: responseChecks,
		})
	}
}
