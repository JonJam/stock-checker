// Using envfile and not YAML due to https://github.com/spf13/viper/issues/1029 so can't remove using . in environment variables which DigitalOcean doesn't support.
// The config methods below don't use viper.UnmarshalKey as it doesn't merge in environment variables due to https://github.com/spf13/viper/issues/1012
package config

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	viper *viper.Viper
}

func NewAppConfig() AppConfig {
	v := viper.GetViper()

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

	return AppConfig{
		viper: v,
	}
}

func (c AppConfig) GetLogConfig() LogConfig {
	const developmentKey = "LOG_DEVELOPMENT"

	keys := []string{
		developmentKey,
	}

	c.checkKeysExist(keys)

	return LogConfig{
		Development: c.viper.GetBool(developmentKey),
	}
}

func (c AppConfig) GetNotifierConfig() NotifierConfig {
	const enabledKey = "NOTIFIER_ENABLED"

	keys := []string{
		enabledKey,
	}

	c.checkKeysExist(keys)

	return NotifierConfig{
		Enabled: c.viper.GetBool(enabledKey),
	}
}

func (c AppConfig) GetSchedulerConfig() SchedulerConfig {
	const intervalKey = "SCHEDULER_INTERVAL"

	keys := []string{
		intervalKey,
	}

	c.checkKeysExist(keys)

	return SchedulerConfig{
		Interval: c.viper.GetUint64(intervalKey),
	}
}

func (c AppConfig) GetTwilioConfig() TwilioConfig {
	const accountSidKey = "TWILIO_ACCOUNTSID"
	const authTokenKey = "TWILIO_AUTHTOKEN"
	const numberToKey = "TWILIO_NUMBERTO"
	const numberFromKey = "TWILIO_NUMBERFROM"

	keys := []string{
		accountSidKey,
		authTokenKey,
		numberToKey,
		numberFromKey,
	}

	c.checkKeysExist(keys)

	return TwilioConfig{
		AccountSid: c.viper.GetString(accountSidKey),
		AuthToken:  c.viper.GetString(authTokenKey),
		NumberTo:   c.viper.GetString(numberToKey),
		NumberFrom: c.viper.GetString(numberFromKey),
	}
}

func (c AppConfig) GetRodConfig() RodConfig {
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

func (c AppConfig) checkKeysExist(keys []string) {
	for _, k := range keys {
		if !c.viper.IsSet(k) {
			log.Fatalf("Configuration key %s not set.\n", k)
		}
	}
}
