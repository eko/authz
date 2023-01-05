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

type CreateRoleRequest struct {
	ID       string   `json:"id" validate:"required"`
	Policies []string `json:"policies" validate:"required"`
}

type UpdateRoleRequest struct {
	Policies []string `json:"policies" validate:"required"`
}

// Creates a new role.
//
//	@security	Authentication
//	@Summary	Creates a new role
//	@Tags		Role
//	@Produce	json
//	@Param		default	body		CreateRoleRequest	true	"Role creation request"
//	@Success	200		{object}	model.Role
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/roles [Post]
func RoleCreate(
	validate *validator.Validate,
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := &CreateRoleRequest{}

		// Parse request body
		if err := c.BodyParser(request); err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Validate body
		if err := validateStruct(validate, request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)

		}

		// Create role
		role, err := manager.CreateRole(request.ID, request.Policies)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(role)
	}
}

// Lists roles.
//
//	@security	Authentication
//	@Summary	Lists roles
//	@Tags		Role
//	@Produce	json
//	@Param		page	query		int		false	"page number"			example(1)
//	@Param		size	query		int		false	"page size"				minimum(1)	maximum(1000)	default(100)
//	@Param		filter	query		string	false	"filter on a field"		example(kind:contains:something)
//	@Param		sort	query		string	false	"sort field and order"	example(kind:desc)
//	@Success	200		{object}	[]model.Role
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/roles [Get]
func RoleList(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, size, err := paginate(c)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		// List roles
		role, total, err := manager.GetRoleRepository().Find(
			database.WithPreloads("Policies"),
			database.WithPage(page),
			database.WithSize(size),
			database.WithFilter(httpFilterToORM(c)),
			database.WithSort(httpSortToORM(c)),
		)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(model.NewPaginated(role, total, page, size))
	}
}

// Retrieve a role.
//
//	@security	Authentication
//	@Summary	Retrieve a role
//	@Tags		Role
//	@Produce	json
//	@Success	200	{object}	model.Role
//	@Failure	404	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/roles/{identifier} [Get]
func RoleGet(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Retrieve role
		role, err := manager.GetRoleRepository().Get(
			identifier,
			database.WithPreloads("Policies"),
		)
		if err != nil {
			statusCode := http.StatusInternalServerError

			if errors.Is(err, gorm.ErrRecordNotFound) {
				statusCode = http.StatusNotFound
			}

			return returnError(c, statusCode,
				fmt.Errorf("cannot retrieve role: %v", err),
			)
		}

		return c.JSON(role)
	}
}

// Updates a role.
//
//	@security	Authentication
//	@Summary	Updates a role
//	@Tags		Role
//	@Produce	json
//	@Param		default	body		UpdateRoleRequest	true	"Role update request"
//	@Success	200		{object}	model.Role
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/roles/{identifier} [Put]
func RoleUpdate(
	validate *validator.Validate,
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Update request
		request := &UpdateRoleRequest{}

		// Parse request body
		if err := c.BodyParser(request); err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Validate body
		if err := validateStruct(validate, request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		// Retrieve role
		role, err := manager.UpdateRole(identifier, request.Policies)
		if err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot update role: %v", err),
			)
		}

		return c.JSON(role)
	}
}

// Deletes a role.
//
//	@security	Authentication
//	@Summary	Deletes a role
//	@Tags		Role
//	@Produce	json
//	@Success	200	{object}	model.Role
//	@Failure	400	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/roles/{identifier} [Get]
func RoleDelete(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Retrieve role
		role, err := manager.GetRoleRepository().Get(identifier)
		if err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot retrieve role: %v", err),
			)
		}

		// Delete role
		if err := manager.GetRoleRepository().Delete(role); err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot delete role: %v", err),
			)
		}

		return c.JSON(model.SuccessResponse{Success: true})
	}
}
