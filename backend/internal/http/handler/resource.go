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

type CreateResourceRequest struct {
	ID         string         `json:"id" validate:"required,slug"`
	Kind       string         `json:"kind" validate:"required,slug"`
	Value      string         `json:"value"`
	Attributes map[string]any `json:"attributes"`
}

// Creates a new resource.
//
//	@security	Authentication
//	@Summary	Creates a new resource
//	@Tags		Resource
//	@Produce	json
//	@Param		default	body		CreateResourceRequest	true	"Resource creation request"
//	@Success	200		{object}	model.Resource
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/resources [Post]
func ResourceCreate(
	validate *validator.Validate,
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := &CreateResourceRequest{}

		// Parse request body
		if err := c.BodyParser(request); err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Validate body
		if err := validateStruct(validate, request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)

		}

		// Create resource
		resource, err := manager.CreateResource(request.ID, request.Kind, request.Value, request.Attributes)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(resource)
	}
}

// Lists resources.
//
//	@security	Authentication
//	@Summary	Lists resources
//	@Tags		Resource
//	@Produce	json
//	@Param		page	query		int		false	"page number"			example(1)
//	@Param		size	query		int		false	"page size"				minimum(1)	maximum(1000)	default(100)
//	@Param		filter	query		string	false	"filter on a field"		example(kind:contains:something)
//	@Param		sort	query		string	false	"sort field and order"	example(kind:desc)
//	@Success	200		{object}	[]model.Resource
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/resources [Get]
func ResourceList(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, size, err := paginate(c)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		// List resources
		resource, total, err := manager.GetResourceRepository().Find(
			database.WithPage(page),
			database.WithSize(size),
			database.WithFilter(httpFilterToORM(c)),
			database.WithSort(httpSortToORM(c)),
			database.WithPreloads("Attributes"),
		)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(model.NewPaginated(resource, total, page, size))
	}
}

// Retrieve a resource.
//
//	@security	Authentication
//	@Summary	Retrieve a resource
//	@Tags		Resource
//	@Produce	json
//	@Success	200	{object}	model.Resource
//	@Failure	404	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/resources/{identifier} [Get]
func ResourceGet(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Retrieve resource
		resource, err := manager.GetResourceRepository().Get(
			identifier,
			database.WithPreloads("Attributes"),
		)
		if err != nil {
			statusCode := http.StatusInternalServerError

			if errors.Is(err, gorm.ErrRecordNotFound) {
				statusCode = http.StatusNotFound
			}

			return returnError(c, statusCode,
				fmt.Errorf("cannot retrieve resource: %v", err),
			)
		}

		return c.JSON(resource)
	}
}

// Deletes a resource.
//
//	@security	Authentication
//	@Summary	Deletes a resource
//	@Tags		Resource
//	@Produce	json
//	@Success	200	{object}	model.Resource
//	@Failure	400	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/resources/{identifier} [Get]
func ResourceDelete(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Retrieve resource
		resource, err := manager.GetResourceRepository().Get(identifier)
		if err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot retrieve resource: %v", err),
			)
		}

		// Delete resource
		if err := manager.GetResourceRepository().Delete(resource); err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot delete resource: %v", err),
			)
		}

		return c.JSON(model.SuccessResponse{Success: true})
	}
}
