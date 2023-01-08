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

type AttributeKeyValue struct {
	Key   string `json:"key" validate:"required"`
	Value any    `json:"value" validate:"required"`
}

type RequestAttributes struct {
	Attributes []AttributeKeyValue `json:"attributes"`
}

func (r RequestAttributes) AttributesMap() map[string]any {
	var result = map[string]any{}

	for _, attribute := range r.Attributes {
		result[attribute.Key] = attribute.Value
	}

	return result
}

type CreateResourceRequest struct {
	RequestAttributes
	ID    string `json:"id" validate:"required,slug"`
	Kind  string `json:"kind" validate:"required,slug"`
	Value string `json:"value"`
}

type UpdateResourceRequest struct {
	RequestAttributes
	Kind  string `json:"kind" validate:"required,slug"`
	Value string `json:"value"`
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
	resourceManager manager.Resource,
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
		resource, err := resourceManager.Create(request.ID, request.Kind, request.Value, request.AttributesMap())
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
	resourceManager manager.Resource,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, size, err := paginate(c)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		// List resources
		resource, total, err := resourceManager.GetRepository().Find(
			repository.WithPage(page),
			repository.WithSize(size),
			repository.WithFilter(httpFilterToORM(c)),
			repository.WithSort(httpSortToORM(c)),
			repository.WithPreloads("Attributes"),
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
	resourceManager manager.Resource,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Retrieve resource
		resource, err := resourceManager.GetRepository().Get(
			identifier,
			repository.WithPreloads("Attributes"),
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

// Updates a resource.
//
//	@security	Authentication
//	@Summary	Updates a resource
//	@Tags		Resource
//	@Produce	json
//	@Param		default	body		UpdateResourceRequest	true	"Resource update request"
//	@Success	200		{object}	model.Resource
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/resources/{identifier} [Put]
func ResourceUpdate(
	validate *validator.Validate,
	resourceManager manager.Resource,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Update request
		request := &UpdateResourceRequest{}

		// Parse request body
		if err := c.BodyParser(request); err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Validate body
		if err := validateStruct(validate, request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		// Retrieve resource
		resource, err := resourceManager.Update(identifier, request.Kind, request.Value, request.AttributesMap())
		if err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot update resource: %v", err),
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
//	@Router		/v1/resources/{identifier} [Delete]
func ResourceDelete(
	resourceManager manager.Resource,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Retrieve resource
		resource, err := resourceManager.GetRepository().Get(identifier)
		if err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot retrieve resource: %v", err),
			)
		}

		// Delete resource
		if err := resourceManager.GetRepository().Delete(resource); err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot delete resource: %v", err),
			)
		}

		return c.JSON(model.SuccessResponse{Success: true})
	}
}
