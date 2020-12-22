// Using envfile and not YAML due to https://github.com/spf13/viper/issues/1029 so can't remove using . in environment variables which DigitalOcean doesn't support.
// The config methods below don't use viper.UnmarshalKey as it doesn't merge in environment variables due to https://github.com/spf13/viper/issues/1012
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	// Setup
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	// Configure directories to look for config
	viper.AddConfigPath(".")

	// Configure environment variables
	viper.SetEnvPrefix("sc")
	viper.AutomaticEnv()

	// Load configuration file
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Could not read config file: %s", err))
	}
}

func GetSchedulerConfig() SchedulerConfig {
	const intervalKey = "SCHEDULER_INTERVAL"

	keys := []string{
		intervalKey,
	}

	checkKeysExist(keys)

	return SchedulerConfig{
		Interval: viper.GetUint64(intervalKey),
	}
}

func GetTwilioConfig() TwilioConfig {
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

	checkKeysExist(keys)

	return TwilioConfig{
		Enabled:    viper.GetBool(enabledKey),
		AccountSid: viper.GetString(accountSidKey),
		AuthToken:  viper.GetString(authTokenKey),
		NumberTo:   viper.GetString(numberToKey),
		NumberFrom: viper.GetString(numberFromKey),
	}
}

func GetRodConfig() RodConfig {
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

	checkKeysExist(keys)

	return RodConfig{
		DevTools:     viper.GetBool(devToolsKey),
		Headless:     viper.GetBool(headlessKey),
		PagePoolSize: viper.GetInt(pagePoolSizeKey),
		SlowMotion:   viper.GetBool(slowMotionKey),
		Trace:        viper.GetBool(traceKey),
	}
}

func checkKeysExist(keys []string) {
	for _, k := range keys {
		if !viper.IsSet(k) {
			panic(fmt.Errorf("Configuration key %s not set", k))
		}
	}
}
