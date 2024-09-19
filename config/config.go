package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type Config struct {
	ServerPort string
	BaseURL    string
	JWT_SECRET string
}

func LoadConfig() (*Config, error) {
	env := os.Getenv("ENV")
	if env == "production" {
		// For Docker, where .env is mounted at the root
		viper.AddConfigPath("/")
		viper.SetConfigFile("/.env")
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigFile(".env")
	}
	viper.SetConfigType("env")
	viper.SetDefault("sever.port", "8080")
	viper.SetDefault("BASE_URL", "http://localhost/api")
	viper.SetDefault("JWT_SECRET", "secret")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
		return nil, err
	}

	viper.AutomaticEnv()
	config := &Config{
		ServerPort: viper.GetString("WEB.PORT"),
		BaseURL:    viper.GetString("BASE_URL"),
		JWT_SECRET: viper.GetString("JWT_SECRET"),
	}
	return config, nil
}
