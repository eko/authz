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

type CreateSubjectRequest struct {
	Value string `json:"value" validate:"required"`
}

type UpdateSubjectRequest struct {
	Value string `json:"value" validate:"required"`
}

// Creates a new subject.
//
//	@security	Authentication
//	@Summary	Creates a new subject
//	@Tags		Authz
//	@Produce	json
//	@Param		default	body		CreateSubjectRequest	true	"Subject creation request"
//	@Success	200		{object}	model.Subject
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/subjects [Post]
func SubjectCreate(
	validate *validator.Validate,
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request := &CreateSubjectRequest{}

		// Parse request body
		if err := c.BodyParser(request); err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Validate body
		if err := validateStruct(validate, request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		// Create subject
		subject, err := manager.CreateSubject(request.Value)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(subject)
	}
}

// Lists subjects.
//
//	@security	Authentication
//	@Summary	Lists subjects
//	@Tags		Authz
//	@Produce	json
//	@Param		page	query		int		false	"page number"			example(1)
//	@Param		size	query		int		false	"page size"				minimum(1)	maximum(1000)	default(100)
//	@Param		filter	query		string	false	"filter on a field"		example(name:contains:something)
//	@Param		sort	query		string	false	"sort field and order"	example(name:desc)
//	@Success	200		{object}	[]model.Subject
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/subjects [Get]
func SubjectList(
	manager manager.Manager,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, size, err := paginate(c)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		// List subjects
		subject, total, err := manager.GetSubjectRepository().Find(
			database.WithPage(page),
			database.WithSize(size),
			database.WithFilter(httpFilterToORM(c)),
			database.WithSort(httpSortToORM(c)),
		)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(model.NewPaginated(subject, total, page, size))
	}
}

// Retrieve a subject.
//
//	@security	Authentication
//	@Summary	Retrieve a subject
//	@Tags		Authz
//	@Produce	json
//	@Success	200	{object}	model.Subject
//	@Failure	404	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/subjects/{identifier} [Get]
func SubjectGet(
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

		// Retrieve subject
		subject, err := manager.GetSubjectRepository().Get(identifier)
		if err != nil {
			statusCode := http.StatusInternalServerError

			if errors.Is(err, gorm.ErrRecordNotFound) {
				statusCode = http.StatusNotFound
			}

			return returnError(c, statusCode,
				fmt.Errorf("cannot retrieve subject: %v", err),
			)
		}

		return c.JSON(subject)
	}
}

// Updates a subject.
//
//	@security	Authentication
//	@Summary	Updates a subject
//	@Tags		Authz
//	@Produce	json
//	@Param		default	body		UpdateSubjectRequest	true	"Subject update request"
//	@Success	200		{object}	model.Subject
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/subjects/{identifier} [Put]
func SubjectUpdate(
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
		request := &UpdateSubjectRequest{}

		// Parse request body
		if err := c.BodyParser(request); err != nil {
			return returnError(c, http.StatusBadRequest, err)
		}

		// Validate body
		if err := validateStruct(validate, request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(err)
		}

		// Retrieve subject
		subject, err := manager.GetSubjectRepository().Get(identifier)
		if err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot retrieve subject: %v", err),
			)
		}

		subject.Value = request.Value

		// Update subject
		if err := manager.GetSubjectRepository().Update(subject); err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot update subject: %v", err),
			)
		}

		return c.JSON(subject)
	}
}

// Deletes a subject.
//
//	@security	Authentication
//	@Summary	Deletes a subject
//	@Tags		Authz
//	@Produce	json
//	@Success	200	{object}	model.Subject
//	@Failure	400	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/subjects/{identifier} [Get]
func SubjectDelete(
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

		// Retrieve subject
		subject, err := manager.GetSubjectRepository().Get(identifier)
		if err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot retrieve subject: %v", err),
			)
		}

		// Delete subject
		if err := manager.GetSubjectRepository().Delete(subject); err != nil {
			return returnError(c, http.StatusInternalServerError,
				fmt.Errorf("cannot delete subject: %v", err),
			)
		}

		return c.JSON(model.SuccessResponse{Success: true})
	}
}
