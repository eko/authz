package handler

import (
	"net/http"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/entity/repository"
	"github.com/eko/authz/backend/internal/http/handler/model"
	"github.com/gofiber/fiber/v2"
)

// Retrieve audits for last days
//
//	@security	Authentication
//	@Summary	Retrieve audits for last days
//	@Tags		Check
//	@Produce	json
//	@Param		page	query		int		false	"page number"			example(1)
//	@Param		size	query		int		false	"page size"				minimum(1)	maximum(1000)	default(100)
//	@Param		filter	query		string	false	"filter on a field"		example(kind:contains:something)
//	@Param		sort	query		string	false	"sort field and order"	example(kind:desc)
//	@Success	200		{object}	[]model.Audit
//	@Failure	400		{object}	model.ErrorResponse
//	@Failure	500		{object}	model.ErrorResponse
//	@Router		/v1/audits [Get]
func AuditGet(
	auditManager manager.Audit,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		page, size, err := paginate(c)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		audits, total, err := auditManager.GetRepository().Find(
			repository.WithPage(page),
			repository.WithSize(size),
			repository.WithFilter(httpFilterToORM(c)),
			repository.WithSort(httpSortToORM(c)),
		)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(model.NewPaginated(audits, total, page, size))
	}
}
