package configs

import "time"

type App struct {
	AuditCleanDelay            time.Duration `config:"app_audit_clean_delay"`
	AuditCleanDaysToKeep       int           `config:"app_audit_clean_days_to_keep"`
	AuditFlushDelay            time.Duration `config:"app_audit_flush_delay"`
	AuditResourceKindRegex     string        `config:"app_audit_resource_kind_regex"`
	DispatcherEventChannelSize int           `config:"dispatcher_event_channel_size"`
	StatsCleanDelay            time.Duration `config:"app_stats_clean_delay"`
	StatsCleanDaysToKeep       int           `config:"app_stats_clean_days_to_keep"`
	StatsFlushDelay            time.Duration `config:"app_stats_flush_delay"`
	StatsResourceKindRegex     string        `config:"app_stats_resource_kind_regex"`
}

func newApp() *App {
	return &App{
		AuditCleanDelay:            1 * time.Hour,
		AuditCleanDaysToKeep:       7,
		AuditFlushDelay:            3 * time.Second,
		AuditResourceKindRegex:     `.*`,
		DispatcherEventChannelSize: 10000,
		StatsCleanDelay:            1 * time.Hour,
		StatsCleanDaysToKeep:       30,
		StatsFlushDelay:            3 * time.Second,
		StatsResourceKindRegex:     `.*`,
	}
}
