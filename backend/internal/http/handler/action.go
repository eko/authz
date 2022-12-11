package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eko/authz/backend/internal/database"
	"github.com/eko/authz/backend/internal/http/handler/model"
	"github.com/eko/authz/backend/internal/manager"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CreateActionRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateActionRequest struct {
	Name string `json:"name" validate:"required"`
}

// Creates a new action.
//
//	@security	Authentication
//	@Summary	Creates a new action
//	@Tags		Authz
//	@Produce	json
//	@Param		default	body		CreateActionRequest	true	"Action creation request"
//	@Success	200		{object}	model.Action
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/actions [Post]
func ActionCreate(
	validate *validator.Validate,
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := &CreateActionRequest{}

		// Parse request body
		if err := c.BodyParser(request); err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Validate body
		if err := validateStruct(validate, request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		// Create action
		action, err := manager.CreateAction(request.Name)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(action)
	}
}

// Lists actions.
//
//	@security	Authentication
//	@Summary	Lists actions
//	@Tags		Authz
//	@Produce	json
//	@Param		page	query		int		false	"page number"			example(1)
//	@Param		size	query		int		false	"page size"				minimum(1)	maximum(1000)	default(100)
//	@Param		filter	query		string	false	"filter on a field"		example(name:contains:something)
//	@Param		sort	query		string	false	"sort field and order"	example(name:desc)
//	@Success	200		{object}	[]model.Action
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/actions [Get]
func ActionList(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, size, err := paginate(c)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		// List actions
		action, total, err := manager.GetActionRepository().Find(
			database.WithPage(page),
			database.WithSize(size),
			database.WithFilter(httpFilterToORM(c)),
			database.WithSort(httpSortToORM(c)),
		)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(model.NewPaginated(action, total, page, size))
	}
}

// Retrieve an action.
//
//	@security	Authentication
//	@Summary	Retrieve an action
//	@Tags		Authz
//	@Produce	json
//	@Success	200	{object}	model.Action
//	@Failure	404	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/actions/{identifier} [Get]
func ActionGet(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifierStr := c.Params("identifier")

		identifier, err := strconv.ParseInt(identifierStr, 10, 64)
		if err != nil {
			return returnError(c, http.StatusBadRequest,
				fmt.Errorf("cannot convert identifier to int64: %v", err),
			)
		}

		// Retrieve action
		action, err := manager.GetActionRepository().Get(identifier)
		if err != nil {
			statusCode := http.StatusInternalServerError

			if errors.Is(err, gorm.ErrRecordNotFound) {
				statusCode = http.StatusNotFound
			}

			return returnError(c, statusCode,
				fmt.Errorf("cannot retrieve action: %v", err),
			)
		}

		return c.JSON(action)
	}
}

// Updates an action.
//
//	@security	Authentication
//	@Summary	Updates an action
//	@Tags		Authz
//	@Produce	json
//	@Param		default	body		UpdateActionRequest	true	"Action update request"
//	@Success	200		{object}	model.Action
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/actions/{identifier} [Put]
func ActionUpdate(
	validate *validator.Validate,
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifierStr := c.Params("identifier")

		identifier, err := strconv.ParseInt(identifierStr, 10, 64)
		if err != nil {
			return returnError(c, http.StatusBadRequest,
				fmt.Errorf("cannot convert identifier to int64: %v", err),
			)
		}

		// Update request
		request := &UpdateActionRequest{}

		// Parse request body
		if err := c.BodyParser(request); err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Validate body
		if err := validateStruct(validate, request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		// Retrieve action
		action, err := manager.GetActionRepository().Get(identifier)
		if err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot retrieve action: %v", err),
			)
		}

		action.Name = request.Name

		// Update action
		if err := manager.GetActionRepository().Update(action); err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot update action: %v", err),
			)
		}

		return c.JSON(action)
	}
}

// Deletes an action.
//
//	@security	Authentication
//	@Summary	Deletes an action
//	@Tags		Authz
//	@Produce	json
//	@Success	200	{object}	model.Action
//	@Failure	400	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/actions/{identifier} [Get]
func ActionDelete(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifierStr := c.Params("identifier")

		identifier, err := strconv.ParseInt(identifierStr, 10, 64)
		if err != nil {
			return returnError(c, http.StatusBadRequest,
				fmt.Errorf("cannot convert identifier to int64: %v", err),
			)
		}

		// Retrieve action
		action, err := manager.GetActionRepository().Get(identifier)
		if err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot retrieve action: %v", err),
			)
		}

		// Delete action
		if err := manager.GetActionRepository().Delete(action); err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot delete action: %v", err),
			)
		}

		return c.JSON(model.SuccessResponse{Success: true})
	}
}
