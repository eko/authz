package configs

import "time"

type App struct {
	StatsFlushDelay time.Duration `config:"app_stats_flush_delay"`
}

func newApp() *App {
	return &App{
		StatsFlushDelay: 3 * time.Second,
	}
}
