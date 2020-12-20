package config

type RodConfig struct {
	PagePoolSize int
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
