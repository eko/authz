package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/entity/repository"
	"github.com/eko/authz/backend/internal/http/handler/model"
	"github.com/eko/authz/backend/internal/http/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var (
	// ErrCannotDeleteOwnAccount is returned when a user tries to delete their own account.
	ErrCannotDeleteOwnAccount = errors.New("a user cannot delete their own account")
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
	userManager manager.User,
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

		user, err := userManager.Create(request.Username, "")
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
	userManager manager.User,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, size, err := paginate(c)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		// List actions
		users, total, err := userManager.GetRepository().Find(
			repository.WithPage(page),
			repository.WithSize(size),
			repository.WithFilter(httpFilterToORM(c)),
			repository.WithSort(httpSortToORM(c)),
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
	userManager manager.User,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")

		// Retrieve user
		user, err := userManager.GetRepository().GetByFields(map[string]repository.FieldValue{
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
//	@Router		/v1/users/{identifier} [Delete]
func UserDelete(
	userManager manager.User,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx := c.UserContext()
		userID := ctx.Value(middleware.UserIdentifierKey).(string)

		identifier := c.Params("identifier")

		if userID == identifier {
			return returnError(c, http.StatusBadRequest, ErrCannotDeleteOwnAccount)
		}

		if err := userManager.Delete(identifier); err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(model.SuccessResponse{Success: true})
	}
}
