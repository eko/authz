package handler

import (
	"net/http"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/entity/repository"
	"github.com/eko/authz/backend/internal/http/handler/model"
	"github.com/gofiber/fiber/v2"
)

// Retrieve compiled policies
//
//	@security	Authentication
//	@Summary	Retrieve compiled policies
//	@Tags		Policy
//	@Produce	json
//	@Success	200	{object}	[]model.CompiledPolicy
//	@Failure	404	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/policies/{identifier}/matches [Get]
func CompiledList(
	compiledManager manager.CompiledPolicy,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, size, err := paginate(c)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		// List policies
		compiledPolicies, total, err := compiledManager.GetRepository().Find(
			repository.WithPage(page),
			repository.WithSize(size),
			repository.WithFilter(httpFilterToORM(c)),
			repository.WithSort(httpSortToORM(c)),
		)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(model.NewPaginated(compiledPolicies, total, page, size))
	}
}
