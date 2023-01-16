package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/entity/repository"
	"github.com/eko/authz/backend/internal/http/handler/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ClientCreateRequest struct {
	Name string `json:"name" validate:"required,slug" example:"my-client"`
}

// Creates a new client
//
// @security  Authentication
// @Summary   Creates a new client
// @Tags      Client
// @Produce   json
// @Param     default  body            ClientCreateRequest  true  "Client creation request"
// @Success   200            {object}  model.Client
// @Failure   400            {object}  model.ErrorResponse
// @Failure   500            {object}  model.ErrorResponse
// @Router    /v1/clients [Post]
func ClientCreate(
	validate *validator.Validate,
	clientManager manager.Client,
	authCfg *configs.Auth,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := &ClientCreateRequest{}

		// Parse request body
		if err := c.BodyParser(request); err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Validate body
		if err := validateStruct(validate, request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		client, err := clientManager.Create(request.Name, authCfg.Domain)
		if err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		return c.JSON(client)
	}
}

// Lists clients.
//
// @security  Authentication
// @Summary   Lists clients
// @Tags      Client
// @Produce   json
// @Param     page    query            int            false                   "page number"      example(1)
// @Param     size    query            int            false                   "page size"          minimum(1)  maximum(1000)  default(100)
// @Param     filter  query            string  false  "filter on a field"                    example(name:contains:something)
// @Param     sort    query            string  false  "sort field and order"  example(name:desc)
// @Success   200            {object}  []model.Client
// @Failure   400            {object}  model.ErrorResponse
// @Failure   500            {object}  model.ErrorResponse
// @Router    /v1/clients [Get]
func ClientList(
	clientManager manager.Client,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, size, err := paginate(c)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		// List actions
		clients, total, err := clientManager.GetRepository().Find(
			repository.WithPage(page),
			repository.WithSize(size),
			repository.WithFilter(httpFilterToORM(c)),
			repository.WithSort(httpSortToORM(c)),
		)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(model.NewPaginated(clients, total, page, size))
	}
}

// Retrieve a client.
//
// @security  Authentication
// @Summary   Retrieve a client
// @Tags      Client
// @Produce   json
// @Success   200  {object}  model.Client
// @Failure   404  {object}  model.ErrorResponse
// @Failure   500  {object}  model.ErrorResponse
// @Router    /v1/clients/{identifier} [Get]
func ClientGet(
	clientManager manager.Client,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Retrieve client
		client, err := clientManager.GetRepository().Get(identifier)
		if err != nil {
			statusCode := http.StatusInternalServerError

			if errors.Is(err, gorm.ErrRecordNotFound) {
				statusCode = http.StatusNotFound
			}

			return returnError(c, statusCode,
				fmt.Errorf("cannot retrieve client: %v", err),
			)
		}

		return c.JSON(client)
	}
}

// Deletes a client.
//
// @security  Authentication
// @Summary   Deletes a client
// @Tags      Client
// @Produce   json
// @Success   200  {object}  model.Client
// @Failure   400  {object}  model.ErrorResponse
// @Failure   500  {object}  model.ErrorResponse
// @Router    /v1/clients/{identifier} [Delete]
func ClientDelete(
	clientManager manager.Client,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		if err := clientManager.Delete(identifier); err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(model.SuccessResponse{Success: true})
	}
}
