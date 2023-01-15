package configs

import "time"

type App struct {
	DispatcherEventChannelSize int           `config:"dispatcher_event_channel_size"`
	StatsCleanDelay            time.Duration `config:"app_stats_clean_delay"`
	StatsCleanDaysToKeep       int           `config:"app_stats_clean_days_to_keep"`
	StatsFlushDelay            time.Duration `config:"app_stats_flush_delay"`
}

func newApp() *App {
	return &App{
		DispatcherEventChannelSize: 10000,
		StatsCleanDelay:            1 * time.Hour,
		StatsCleanDaysToKeep:       30,
		StatsFlushDelay:            3 * time.Second,
	}
}
