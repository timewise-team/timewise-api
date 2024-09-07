package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	ServerPort string
}

func LoadConfig() (*Config, error) {
	// Load config here
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")

	viper.SetDefault("sever.port", "3000")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
		return nil, err
	}

	viper.AutomaticEnv()
	config := &Config{
		ServerPort: viper.GetString("WEB.PORT"),
	}
	return config, nil
}
