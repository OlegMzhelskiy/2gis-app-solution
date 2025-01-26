package config

import (
	"fmt"

	"applicationDesignTest/pkg/log"

	"github.com/spf13/viper"
)

type Server struct {
	Port string `mapstructure:"port"`
}

type Config struct {
	Server `mapstructure:"server"`
}

func LoadConfig(configPath string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	viper.AddConfigPath("../../" + configPath)

	viper.AutomaticEnv()

	viper.SetEnvPrefix("")

	// SERVER_PORT
	if err := viper.BindEnv("server.port"); err != nil {
		return nil, fmt.Errorf("failed to bind env: %w", err)
	}

	viper.SetDefault("server.port", "8080")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Error("config file not found, falling back to environment variables", err)
		} else {
			log.Error("error reading config file", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return &config, nil
}
