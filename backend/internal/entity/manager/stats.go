package manager

import (
	"errors"
	"fmt"
	"time"

	"github.com/eko/authz/backend/internal/entity/model"
	"github.com/eko/authz/backend/internal/entity/repository"
	"gorm.io/gorm"
)

type StatsRepository repository.Base[model.Stats]

type Stats interface {
	BatchAddCheck(timestamp int64, allowed int64, denied int64) error
	GetRepository() StatsRepository
}

type statsManager struct {
	repository StatsRepository
}

// NewStats initializes a new stats manager.
func NewStats(
	repository StatsRepository,
) Stats {
	return &statsManager{
		repository: repository,
	}
}

func (m *statsManager) GetRepository() StatsRepository {
	return m.repository
}

func (m *statsManager) BatchAddCheck(timestamp int64, allowed int64, denied int64) error {
	date := time.Unix(timestamp, 0)
	formattedDate := date.Format("20060102")

	stats, err := m.repository.Get(formattedDate)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("unable to get stats for date %s: %v", formattedDate, err)
	}

	found := stats != nil

	if !found {
		stats = &model.Stats{
			ID:   formattedDate,
			Date: date.Format(time.RFC3339),
		}
	}

	stats.ChecksAllowedNumber = stats.ChecksAllowedNumber + allowed
	stats.ChecksDeniedNumber = stats.ChecksDeniedNumber + denied

	if found {
		return m.repository.Update(stats)
	}

	return m.repository.Create(stats)
}
