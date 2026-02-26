package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Environment string
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.SetDefault("SERVER_PORT", "8010")
	viper.SetDefault("ENVIRONMENT", "development")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Note: No .env file found, using environment variables\n")
		}
	}

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("SERVER_PORT"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetString("POSTGRES_PORT"),
			User:     viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			DBName:   viper.GetString("POSTGRES_DB"),
			SSLMode:  viper.GetString("POSTGRES_SSLMODE"),
		},
		Environment: viper.GetString("ENVIRONMENT"),
	}

	return config, nil
}
