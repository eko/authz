package handler

import (
	"net/http"

	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/entity/repository"
	"github.com/gofiber/fiber/v2"
)

// Retrieve statistics for last days
//
//	@security	Authentication
//	@Summary	Retrieve statistics for last days
//	@Tags		Check
//	@Produce	json
//	@Success	200	{object}	[]model.Stats
//	@Failure	400	{object}	model.ErrorResponse
//	@Failure	500	{object}	model.ErrorResponse
//	@Router		/v1/stats [Get]
func StatsGet(
	statsManager manager.Stats,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		stats, _, err := statsManager.GetRepository().Find(
			repository.WithSort("date desc"),
		)
		if err != nil {
			return returnError(c, http.StatusInternalServerError, err)
		}

		return c.JSON(stats)
	}
}
