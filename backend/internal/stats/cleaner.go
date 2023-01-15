package stats

import (
	"context"
	lib_time "time"

	"github.com/eko/authz/backend/configs"
	"github.com/eko/authz/backend/internal/entity/manager"
	"github.com/eko/authz/backend/internal/entity/repository"
	"github.com/eko/authz/backend/internal/helper/time"
	"go.uber.org/fx"
	"golang.org/x/exp/slog"
)

type cleaner struct {
	logger       *slog.Logger
	clock        time.Clock
	statsManager manager.Stats
	cleanDelay   lib_time.Duration
	daysToKeep   int
}

func NewCleaner(
	cfg *configs.App,
	logger *slog.Logger,
	clock time.Clock,
	statsManager manager.Stats,
) *cleaner {
	return &cleaner{
		logger:       logger,
		clock:        clock,
		statsManager: statsManager,
		cleanDelay:   cfg.StatsCleanDelay,
		daysToKeep:   cfg.StatsCleanDaysToKeep,
	}
}

func RunCleaner(lc fx.Lifecycle, cleaner *cleaner) {
	ticker := lib_time.NewTicker(cleaner.cleanDelay)

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			go func() {
				for range ticker.C {
					cleaner.logger.Info("Stats: cleaning stats older than 30 days")

					if err := cleaner.statsManager.GetRepository().DeleteByFields(map[string]repository.FieldValue{
						"date": {
							Operator: "<=",
							Value:    cleaner.clock.Now().AddDate(0, 0, -cleaner.daysToKeep),
						},
					}); err != nil {
						cleaner.logger.Error("Stats: unable to clean stats", err)
					}
				}
			}()

			cleaner.logger.Info("Stats: cleaner started")

			return nil
		},
		OnStop: func(_ context.Context) error {
			ticker.Stop()

			cleaner.logger.Info("Stats: cleaner stopped")

			return nil
		},
	})
}
