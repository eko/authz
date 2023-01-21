package metric

import (
	"strconv"

	"github.com/eko/authz/backend/configs"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	checkCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "authz_check_counter",
		Help: "The total number of checks processed",
	}, []string{"is_allowed", "resource_kind"})

	itemCreatedCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "authz_item_counter",
		Help: "The total number of items (resource, policy, ...) created or updated in database",
	}, []string{"item_type", "action"})
)

type Observer interface {
	ObserveCheckCounter(resourceKind string, isAllowed bool)
	ObserveItemCreatedCounter(itemType, action string)
}

type observer struct {
	checkCounter       *prometheus.CounterVec
	itemCreatedCounter *prometheus.CounterVec
}

func NewObserver(
	cfg *configs.App,
) (Observer, error) {
	if !cfg.MetricsEnabled {
		return nil, nil
	}

	observer := &observer{
		checkCounter:       checkCounter,
		itemCreatedCounter: itemCreatedCounter,
	}

	if err := observer.initialize(); err != nil {
		return nil, err
	}

	return observer, nil
}

func (r *observer) initialize() (err error) {
	if err := prometheus.Register(checkCounter); err != nil {
		return err
	}

	if err := prometheus.Register(itemCreatedCounter); err != nil {
		return err
	}

	return nil
}

func (r *observer) ObserveCheckCounter(resourceKind string, isAllowed bool) {
	r.checkCounter.WithLabelValues(
		strconv.FormatBool(isAllowed),
		resourceKind,
	).Inc()
}

func (r *observer) ObserveItemCreatedCounter(itemType string, action string) {
	r.itemCreatedCounter.WithLabelValues(
		itemType,
		action,
	).Inc()
}
