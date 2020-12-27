package config

type Config interface {
	GetLogConfig() LogConfig
	GetNotifierConfig() NotifierConfig
	GetRodConfig() RodConfig
	GetSchedulerConfig() SchedulerConfig
	GetTwilioConfig() TwilioConfig
}

type LogConfig struct {
	Development bool
}

type NotifierConfig struct {
	Enabled bool
}

type RodConfig struct {
	DevTools     bool
	Headless     bool
	PagePoolSize int
	SlowMotion   bool
	Trace        bool
}

type SchedulerConfig struct {
	Interval uint64
}

type TwilioConfig struct {
	AccountSid string
	AuthToken  string
	NumberTo   string
	NumberFrom string
}
