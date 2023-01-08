package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/eko/authz/backend/internal/entity/repository"
	"github.com/eko/authz/backend/internal/http/handler/model"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

const (
	// filterHTTPDelimiter character allows to delimit the field
	// to filter (example: name), the operator (example: contains) and the value (example: hello).
	filterHTTPDelimiter = ":"

	// sortHTTPDelimiter character allows to delimit the field
	// to sort (example: name) and the order (example: asc).
	sortHTTPDelimiter = ":"
)

func paginate(c *fiber.Ctx) (int64, int64, error) {
	pageStr := c.Query("page")
	sizeStr := c.Query("size")

	page, detailErr := convertStringToInt64(pageStr)
	if detailErr != nil {
		return 0, 0, detailErr
	}

	// Decrement page to avoid having page 0. Starts at 1.
	if page > 0 {
		page--
	}

	size, detailErr := convertStringToInt64(sizeStr)
	if detailErr != nil {
		return 0, 0, detailErr
	}

	// Set a default size value in case no size is specified in HTTP request.
	if size == 0 || size > 1000 {
		size = 100
	}

	return page, size, nil
}

func convertStringToInt64(value string) (int64, error) {
	if value == "" {
		return 0, nil
	}

	intValue, err := strconv.ParseInt(value, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("unable to convert string to uint64: %v", err)
	}

	return intValue, nil
}

func httpFilterToORM(c *fiber.Ctx) map[string]repository.FieldValue {
	var result = map[string]repository.FieldValue{}

	if c == nil {
		return result
	}

	filter := c.Query("filter")
	values := strings.Split(filter, filterHTTPDelimiter)

	// We should have 3 values:
	// - field
	// - operator
	// - value
	if len(values) != 3 {
		return result
	}

	field, operator, value := values[0], values[1], values[2]

	switch operator {
	case "contains":
		result[field] = repository.FieldValue{Operator: "LIKE", Value: "%" + value + "%"}
	case "is":
		result[field] = repository.FieldValue{Operator: "=", Value: value}
	}

	return result
}

func httpSortToORM(c *fiber.Ctx) string {
	if c == nil {
		return ""
	}

	sort := c.Query("sort")
	values := strings.Split(sort, sortHTTPDelimiter)

	// We should have 2 values:
	// - field
	// - sort order
	if len(values) != 2 {
		return ""
	}

	field, order := values[0], values[1]

	if order != "asc" && order != "desc" {
		return ""
	}

	return fmt.Sprintf("%s %s", field, order)
}

func returnError(c *fiber.Ctx, statusCode int, err error) error {
	c.Response().SetStatusCode(statusCode)

	return c.JSON(model.ErrorResponse{
		Error:   true,
		Message: err.Error(),
	})
}

func returnHTTPError(w http.ResponseWriter, statusCode int, err error) {
	w.WriteHeader(statusCode)

	responseBytes, err := json.Marshal(model.ErrorResponse{
		Error:   true,
		Message: err.Error(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(responseBytes)
}

func validateStruct(validate *validator.Validate, request any) *model.ErrorResponse {
	var errors []*model.ValidateErrorResponse

	if err := validate.Struct(request); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, &model.ValidateErrorResponse{
				FailedField: err.Field(),
				Tag:         err.Tag(),
				Value:       err.Param(),
			})
		}
	}

	if len(errors) == 0 {
		return nil
	}

	return &model.ErrorResponse{
		Error:            true,
		Message:          "Unable to validate request.",
		ValidationErrors: errors,
	}
}
