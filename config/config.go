package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	// Setup
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

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

// The config methods below don't use viper.UnmarshalKey as it doesn't merge in environment variables due to
// https://github.com/spf13/viper/issues/1012

func GetSchedulerConfig() SchedulerConfig {
	const intervalKey = "scheduler.interval"

	keys := []string{
		intervalKey,
	}

	checkKeysExist(keys)

	return SchedulerConfig{
		Interval: viper.GetUint64(intervalKey),
	}
}

func GetTwilioConfig() TwilioConfig {
	const enabledKey = "twilio.enabled"
	const accountSidKey = "twilio.accountSid"
	const authTokenKey = "twilio.authToken"
	const numberToKey = "twilio.numberTo"
	const numberFromKey = "twilio.numberFrom"

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
	const devToolsKey = "rod.devTools"
	const headlessKey = "rod.headless"
	const pagePoolSizeKey = "rod.pagePoolSize"
	const slowMotionKey = "rod.slowMotion"
	const traceKey = "rod.trace"

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
