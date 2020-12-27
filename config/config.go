// Using envfile and not YAML due to https://github.com/spf13/viper/issues/1029 so can't remove using . in environment variables which DigitalOcean doesn't support.
// The config methods below don't use viper.UnmarshalKey as it doesn't merge in environment variables due to https://github.com/spf13/viper/issues/1012
package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper
}

func New() Config {
	return newWithViper(viper.GetViper())
}

// For testing only
func newWithViper(v *viper.Viper) Config {
	// Setup
	v.SetConfigName("config")
	v.SetConfigType("env")

	// Configure directories to look for config
	v.AddConfigPath(".")

	// Configure environment variables
	v.SetEnvPrefix("sc")
	v.AutomaticEnv()

	// Load configuration file
	err := v.ReadInConfig()

	if err != nil {
		log.Fatalf("Could not read config file: %s.\n", err)
	}

	return Config{
		viper: v,
	}
}

func (c Config) GetLogConfig() LogConfig {
	const developmentKey = "LOG_DEVELOPMENT"

	keys := []string{
		developmentKey,
	}

	c.checkKeysExist(keys)

	return LogConfig{
		Development: c.viper.GetBool(developmentKey),
	}
}

func (c Config) GetSchedulerConfig() SchedulerConfig {
	const intervalKey = "SCHEDULER_INTERVAL"

	keys := []string{
		intervalKey,
	}

	c.checkKeysExist(keys)

	return SchedulerConfig{
		Interval: c.viper.GetUint64(intervalKey),
	}
}

func (c Config) GetTwilioConfig() TwilioConfig {
	const enabledKey = "TWILIO_ENABLED"
	const accountSidKey = "TWILIO_ACCOUNTSID"
	const authTokenKey = "TWILIO_AUTHTOKEN"
	const numberToKey = "TWILIO_NUMBERTO"
	const numberFromKey = "TWILIO_NUMBERFROM"

	keys := []string{
		enabledKey,
		accountSidKey,
		authTokenKey,
		numberToKey,
		numberFromKey,
	}

	c.checkKeysExist(keys)

	return TwilioConfig{
		Enabled:    c.viper.GetBool(enabledKey),
		AccountSid: c.viper.GetString(accountSidKey),
		AuthToken:  c.viper.GetString(authTokenKey),
		NumberTo:   c.viper.GetString(numberToKey),
		NumberFrom: c.viper.GetString(numberFromKey),
	}
}

func (c Config) GetRodConfig() RodConfig {
	const devToolsKey = "ROD_DEVTOOLS"
	const headlessKey = "ROD_HEADLESS"
	const pagePoolSizeKey = "ROD_PAGEPOOLSIZE"
	const slowMotionKey = "ROD_SLOWMOTION"
	const traceKey = "ROD_TRACE"

	keys := []string{
		devToolsKey,
		headlessKey,
		pagePoolSizeKey,
		slowMotionKey,
		traceKey,
	}

	c.checkKeysExist(keys)

	return RodConfig{
		DevTools:     c.viper.GetBool(devToolsKey),
		Headless:     c.viper.GetBool(headlessKey),
		PagePoolSize: c.viper.GetInt(pagePoolSizeKey),
		SlowMotion:   c.viper.GetBool(slowMotionKey),
		Trace:        c.viper.GetBool(traceKey),
	}
}

func (c Config) checkKeysExist(keys []string) {
	for _, k := range keys {
		if !c.viper.IsSet(k) {
			log.Fatalf("Configuration key %s not set.\n", k)
		}
	}
}
