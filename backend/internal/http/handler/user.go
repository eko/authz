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

type UserCreateRequest struct {
	Username string `json:"username" validate:"required,slug" example:"my-user"`
}

// Creates a new user
//
//	@security	Authentication
//	@Summary	Creates a new user
//	@Tags		User
//	@Produce	json
//	@Param		default	body		UserCreateRequest	true	"User creation request"
//	@Success	200		{object}	model.User
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/users [Post]
func UserCreate(
	validate *validator.Validate,
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := &UserCreateRequest{}

		// Parse request body
		if err := c.BodyParser(request); err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Validate body
		if err := validateStruct(validate, request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		user, err := manager.CreateUser(request.Username, "")
		if err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		return c.JSON(user)
	}
}

// Lists users.
//
//	@security	Authentication
//	@Summary	Lists users
//	@Tags		User
//	@Produce	json
//	@Param		page	query		int		false	"page number"			example(1)
//	@Param		size	query		int		false	"page size"				minimum(1)	maximum(1000)	default(100)
//	@Param		filter	query		string	false	"filter on a field"		example(name:contains:something)
//	@Param		sort	query		string	false	"sort field and order"	example(name:desc)
//	@Success	200		{object}	[]model.User
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/users [Get]
func UserList(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, size, err := paginate(c)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		// List actions
		users, total, err := manager.GetUserRepository().Find(
			database.WithPage(page),
			database.WithSize(size),
			database.WithFilter(httpFilterToORM(c)),
			database.WithSort(httpSortToORM(c)),
		)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(model.NewPaginated(users, total, page, size))
	}
}

// Retrieve a user.
//
//	@security	Authentication
//	@Summary	Retrieve a user
//	@Tags		User
//	@Produce	json
//	@Success	200	{object}	model.User
//	@Failure	404	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/users/{identifier} [Get]
func UserGet(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Retrieve user
		user, err := manager.GetUserRepository().GetByFields(map[string]database.FieldValue{
			"username": {Operator: "=", Value: identifier},
		})
		if err != nil {
			statusCode := http.StatusInternalServerError

			if errors.Is(err, gorm.ErrRecordNotFound) {
				statusCode = http.StatusNotFound
			}

			return returnError(c, statusCode,
				fmt.Errorf("cannot retrieve user: %v", err),
			)
		}

		return c.JSON(user)
	}
}

// Deletes a user.
//
//	@security	Authentication
//	@Summary	Deletes a user
//	@Tags		User
//	@Produce	json
//	@Success	200	{object}	model.User
//	@Failure	400	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/users/{identifier} [Get]
func UserDelete(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Retrieve user
		user, err := manager.GetUserRepository().Get(identifier)
		if err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot retrieve user: %v", err),
			)
		}

		// Delete user
		if err := manager.GetUserRepository().Delete(user); err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot delete user: %v", err),
			)
		}

		return c.JSON(model.SuccessResponse{Success: true})
	}
}
