package config

import "github.com/spf13/viper"

// ReadConfig reads in config from config.yml
func ReadConfig() (*viper.Viper, error) {
	viperConfig := viper.New()

	viperConfig.AutomaticEnv()
	viperConfig.SetConfigType("yaml")
	viperConfig.SetConfigName("config")
	viperConfig.AddConfigPath("./config")
	viperConfig.AddConfigPath(".")

	viperConfig.SetDefault("MAX_WORKERS", 1)
	viperConfig.SetDefault("AWS_PROFILE", "default")

	err := viperConfig.ReadInConfig()
	return viperConfig, err
}
