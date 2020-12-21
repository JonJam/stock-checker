package config

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
	Enabled    bool
	AccountSid string
	AuthToken  string
	NumberTo   string
	NumberFrom string
}
