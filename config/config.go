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

func IsDevMode() bool {
	const key = "developmentMode"

	if !viper.IsSet(key) {
		panic(fmt.Errorf("Configuration key %s not set", key))
	}

	return viper.GetBool(key)
}

// TODO Implement scheduler config

func GetTwilioConfig() TwilioConfig {
	const key = "twilio"

	c := TwilioConfig{}

	if !viper.IsSet(key) {
		panic(fmt.Errorf("Configuration key %s not set", key))
	}

	if err := viper.UnmarshalKey(key, &c); err != nil {
		panic(fmt.Errorf("Unable to decode configuration with key %s: %s", key, err))
	}

	return c
}

func GetRodConfig() RodConfig {
	const key = "rod"

	c := RodConfig{}

	if !viper.IsSet(key) {
		panic(fmt.Errorf("Configuration key %s not set", key))
	}

	if err := viper.UnmarshalKey(key, &c); err != nil {
		panic(fmt.Errorf("Unable to decode configuration with key %s: %s", key, err))
	}

	return c
}
