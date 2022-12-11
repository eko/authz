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

type PolicyResourceRequest struct {
	Kind  string `json:"kind" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type PolicyResourceRequests []PolicyResourceRequest

func (r PolicyResourceRequests) ToMap() []map[string]string {
	var resources = make([]map[string]string, 0)

	for _, resource := range r {
		resources = append(resources, map[string]string{
			"kind":  resource.Kind,
			"value": resource.Value,
		})
	}

	return resources
}

type CreatePolicyRequest struct {
	Name      string                 `json:"name" validate:"required"`
	Resources PolicyResourceRequests `json:"resources" validate:"required"`
	Actions   []string               `json:"actions" validate:"required"`
}

type UpdatePolicyRequest struct {
	Name      string                 `json:"name" validate:"required"`
	Resources PolicyResourceRequests `json:"resources" validate:"required"`
	Actions   []string               `json:"actions" validate:"required"`
}

// Creates a new policy.
//
//	@security	Authentication
//	@Summary	Creates a new policy
//	@Tags		Authz
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
		policy, err := manager.CreatePolicy(request.Name, request.Resources.ToMap(), request.Actions)
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
//	@Tags		Authz
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
//	@Tags		Authz
//	@Produce	json
//	@Success	200	{object}	model.Policy
//	@Failure	404	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/policies/{identifier} [Get]
func PolicyGet(
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
//	@Tags		Authz
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
		identifierStr := c.Params("identifier")

		identifier, err := strconv.ParseInt(identifierStr, 10, 64)
		if err != nil {
			return returnError(c, http.StatusBadRequest,
				fmt.Errorf("cannot convert identifier to int64: %v", err),
			)
		}

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
		policy, err := manager.UpdatePolicy(identifier, request.Name, request.Resources.ToMap(), request.Actions)
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
//	@Tags		Authz
//	@Produce	json
//	@Success	200	{object}	model.Policy
//	@Failure	400	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/policies/{identifier} [Get]
func PolicyDelete(
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
