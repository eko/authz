package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/http/handler/model"
	"github.com/eko/authz/backend/internal/manager"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CreatePolicyRequest struct {
	ID             string   `json:"id" validate:"required,slug"`
	Resources      []string `json:"resources" validate:"required,dive,slug"`
	Actions        []string `json:"actions" validate:"required,dive,slug"`
	AttributeRules []string `json:"attribute_rules"`
}

type UpdatePolicyRequest struct {
	Resources      []string `json:"resources" validate:"required,dive,slug"`
	Actions        []string `json:"actions" validate:"required,dive,slug"`
	AttributeRules []string `json:"attribute_rules"`
}

// Creates a new policy.
//
//	@security	Authentication
//	@Summary	Creates a new policy
//	@Tags		Policy
//	@Produce	json
//	@Param		default	body		CreatePolicyRequest	true	"Policy creation request"
//	@Success	200		{object}	model.Policy
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/policies [Post]
func PolicyCreate(
	validate *validator.Validate,
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := &CreatePolicyRequest{}

		// Parse request body
		if err := c.BodyParser(request); err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Validate body
		if err := validateStruct(validate, request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		// Create policy
		policy, err := manager.CreatePolicy(
			request.ID,
			request.Resources,
			request.Actions,
			request.AttributeRules,
		)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(policy)
	}
}

// Lists policies.
//
//	@security	Authentication
//	@Summary	Lists policies
//	@Tags		Policy
//	@Produce	json
//	@Param		page	query		int		false	"page number"			example(1)
//	@Param		size	query		int		false	"page size"				minimum(1)	maximum(1000)	default(100)
//	@Param		filter	query		string	false	"filter on a field"		example(kind:contains:something)
//	@Param		sort	query		string	false	"sort field and order"	example(kind:desc)
//	@Success	200		{object}	[]model.Policy
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/policies [Get]
func PolicyList(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, size, err := paginate(c)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		// List policies
		policy, total, err := manager.GetPolicyRepository().Find(
			database.WithPreloads("Resources", "Actions"),
			database.WithPage(page),
			database.WithSize(size),
			database.WithFilter(httpFilterToORM(c)),
			database.WithSort(httpSortToORM(c)),
		)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(model.NewPaginated(policy, total, page, size))
	}
}

// Retrieve a policy.
//
//	@security	Authentication
//	@Summary	Retrieve a policy
//	@Tags		Policy
//	@Produce	json
//	@Success	200	{object}	model.Policy
//	@Failure	404	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/policies/{identifier} [Get]
func PolicyGet(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Retrieve policy
		policy, err := manager.GetPolicyRepository().Get(
			identifier,
			database.WithPreloads("Resources", "Actions"),
		)
		if err != nil {
			statusCode := http.StatusInternalServerError

			if errors.Is(err, gorm.ErrRecordNotFound) {
				statusCode = http.StatusNotFound
			}

			return returnError(c, statusCode,
				fmt.Errorf("cannot retrieve policy: %v", err),
			)
		}

		return c.JSON(policy)
	}
}

// Updates a policy.
//
//	@security	Authentication
//	@Summary	Updates a policy
//	@Tags		Policy
//	@Produce	json
//	@Param		default	body		UpdatePolicyRequest	true	"Policy update request"
//	@Success	200		{object}	model.Policy
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/policies/{identifier} [Put]
func PolicyUpdate(
	validate *validator.Validate,
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Update request
		request := &UpdatePolicyRequest{}

		// Parse request body
		if err := c.BodyParser(request); err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Validate body
		if err := validateStruct(validate, request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		// Retrieve policy
		policy, err := manager.UpdatePolicy(
			identifier,
			request.Resources,
			request.Actions,
			request.AttributeRules,
		)
		if err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot update policy: %v", err),
			)
		}

		return c.JSON(policy)
	}
}

// Deletes a policy.
//
//	@security	Authentication
//	@Summary	Deletes a policy
//	@Tags		Policy
//	@Produce	json
//	@Success	200	{object}	model.Policy
//	@Failure	400	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/policies/{identifier} [Delete]
func PolicyDelete(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Retrieve policy
		policy, err := manager.GetPolicyRepository().Get(identifier)
		if err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot retrieve policy: %v", err),
			)
		}

		// Delete policy
		if err := manager.GetPolicyRepository().Delete(policy); err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot delete policy: %v", err),
			)
		}

		return c.JSON(model.SuccessResponse{Success: true})
	}
}
