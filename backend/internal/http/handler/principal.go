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

type CreatePrincipalRequest struct {
	ID         string         `json:"id" validate:"required,slug"`
	Roles      []string       `json:"roles" validate:"dive,slug"`
	Attributes map[string]any `json:"attributes"`
}

type UpdatePrincipalRequest struct {
	Roles      []string       `json:"roles" validate:"dive,slug"`
	Attributes map[string]any `json:"attributes"`
}

// Creates a new principal.
//
//	@security	Authentication
//	@Summary	Creates a new principal
//	@Tags		Authz
//	@Produce	json
//	@Param		default	body		CreatePrincipalRequest	true	"Principal creation request"
//	@Success	200		{object}	model.Principal
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/principals [Post]
func PrincipalCreate(
	validate *validator.Validate,
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := &CreatePrincipalRequest{}

		// Parse request body
		if err := c.BodyParser(request); err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Validate body
		if err := validateStruct(validate, request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		// Create principal
		principal, err := manager.CreatePrincipal(request.ID, request.Roles, request.Attributes)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(principal)
	}
}

// Lists principals.
//
//	@security	Authentication
//	@Summary	Lists principals
//	@Tags		Authz
//	@Produce	json
//	@Param		page	query		int		false	"page number"			example(1)
//	@Param		size	query		int		false	"page size"				minimum(1)	maximum(1000)	default(100)
//	@Param		filter	query		string	false	"filter on a field"		example(name:contains:something)
//	@Param		sort	query		string	false	"sort field and order"	example(name:desc)
//	@Success	200		{object}	[]model.Principal
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/principals [Get]
func PrincipalList(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, size, err := paginate(c)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		// List principals
		principal, total, err := manager.GetPrincipalRepository().Find(
			database.WithPage(page),
			database.WithSize(size),
			database.WithFilter(httpFilterToORM(c)),
			database.WithSort(httpSortToORM(c)),
		)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(model.NewPaginated(principal, total, page, size))
	}
}

// Retrieve a principal.
//
//	@security	Authentication
//	@Summary	Retrieve a principal
//	@Tags		Authz
//	@Produce	json
//	@Success	200	{object}	model.Principal
//	@Failure	404	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/principals/{identifier} [Get]
func PrincipalGet(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Retrieve principal
		principal, err := manager.GetPrincipalRepository().Get(identifier)
		if err != nil {
			statusCode := http.StatusInternalServerError

			if errors.Is(err, gorm.ErrRecordNotFound) {
				statusCode = http.StatusNotFound
			}

			return returnError(c, statusCode,
				fmt.Errorf("cannot retrieve principal: %v", err),
			)
		}

		return c.JSON(principal)
	}
}

// Updates a principal.
//
//	@security	Authentication
//	@Summary	Updates a principal
//	@Tags		Authz
//	@Produce	json
//	@Param		default	body		UpdatePrincipalRequest	true	"Principal update request"
//	@Success	200		{object}	model.Principal
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/principals/{identifier} [Put]
func PrincipalUpdate(
	validate *validator.Validate,
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Update request
		request := &UpdatePrincipalRequest{}

		// Parse request body
		if err := c.BodyParser(request); err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Validate body
		if err := validateStruct(validate, request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		// Retrieve principal
		principal, err := manager.UpdatePrincipal(identifier, request.Roles, request.Attributes)
		if err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot update principal: %v", err),
			)
		}

		return c.JSON(principal)
	}
}

// Deletes a principal.
//
//	@security	Authentication
//	@Summary	Deletes a principal
//	@Tags		Authz
//	@Produce	json
//	@Success	200	{object}	model.Principal
//	@Failure	400	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/principals/{identifier} [Get]
func PrincipalDelete(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Retrieve principal
		principal, err := manager.GetPrincipalRepository().Get(identifier)
		if err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot retrieve principal: %v", err),
			)
		}

		// Delete principal
		if err := manager.GetPrincipalRepository().Delete(principal); err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot delete principal: %v", err),
			)
		}

		return c.JSON(model.SuccessResponse{Success: true})
	}
}
