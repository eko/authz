package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/entity/repository"
	"github.com/eko/authz/backend/internal/http/handler/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CreatePrincipalRequest struct {
	RequestAttributes
	ID    string   `json:"id" validate:"required,slug"`
	Roles []string `json:"roles" validate:"dive,slug"`
}

type UpdatePrincipalRequest struct {
	RequestAttributes
	Roles []string `json:"roles" validate:"dive,slug"`
}

// Creates a new principal.
//
//	@security	Authentication
//	@Summary	Creates a new principal
//	@Tags		Principal
//	@Produce	json
//	@Param		default	body		CreatePrincipalRequest	true	"Principal creation request"
//	@Success	200		{object}	model.Principal
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/principals [Post]
func PrincipalCreate(
	validate *validator.Validate,
	principalManager manager.Principal,
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
		principal, err := principalManager.Create(request.ID, request.Roles, request.AttributesMap())
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
//	@Tags		Principal
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
	principalManager manager.Principal,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, size, err := paginate(c)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		// List principals
		principal, total, err := principalManager.GetRepository().Find(
			repository.WithPage(page),
			repository.WithSize(size),
			repository.WithFilter(httpFilterToORM(c)),
			repository.WithSort(httpSortToORM(c)),
			repository.WithPreloads("Attributes", "Roles"),
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
//	@Tags		Principal
//	@Produce	json
//	@Success	200	{object}	model.Principal
//	@Failure	404	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/principals/{identifier} [Get]
func PrincipalGet(
	principalManager manager.Principal,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Retrieve principal
		principal, err := principalManager.GetRepository().Get(
			identifier,
			repository.WithPreloads("Attributes", "Roles"),
		)
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
//	@Tags		Principal
//	@Produce	json
//	@Param		default	body		UpdatePrincipalRequest	true	"Principal update request"
//	@Success	200		{object}	model.Principal
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/principals/{identifier} [Put]
func PrincipalUpdate(
	validate *validator.Validate,
	principalManager manager.Principal,
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
		principal, err := principalManager.Update(identifier, request.Roles, request.AttributesMap())
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
//	@Tags		Principal
//	@Produce	json
//	@Success	200	{object}	model.Principal
//	@Failure	400	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/principals/{identifier} [Delete]
func PrincipalDelete(
	principalManager manager.Principal,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Retrieve principal
		principal, err := principalManager.GetRepository().Get(identifier)
		if err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot retrieve principal: %v", err),
			)
		}

		// Delete principal
		if err := principalManager.GetRepository().Delete(principal); err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot delete principal: %v", err),
			)
		}

		return c.JSON(model.SuccessResponse{Success: true})
	}
}
