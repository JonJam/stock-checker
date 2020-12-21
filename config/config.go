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

func GetSchedulerConfig() SchedulerConfig {
	const key = "scheduler"

	c := SchedulerConfig{}

	unmarshalComplexConfig(key, &c)

	return c
}

func GetTwilioConfig() TwilioConfig {
	const key = "twilio"

	c := TwilioConfig{}

	unmarshalComplexConfig(key, &c)

	return c
}

func GetRodConfig() RodConfig {
	const key = "rod"

	c := RodConfig{}

	unmarshalComplexConfig(key, &c)

	return c
}

func unmarshalComplexConfig(key string, c interface{}) {
	if !viper.IsSet(key) {
		panic(fmt.Errorf("Configuration key %s not set", key))
	}

	if err := viper.UnmarshalKey(key, c); err != nil {
		panic(fmt.Errorf("Unable to decode configuration with key %s: %s", key, err))
	}
}
