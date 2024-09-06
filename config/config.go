package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	ServerPort string
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
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
		DBUser:     viper.GetString("DB.USERNAME"),
		DBPassword: viper.GetString("DB.PASSWORD"),
		DBName:     viper.GetString("DB.NAME"),
		DBHost:     viper.GetString("DB.HOST"),
		DBPort:     viper.GetString("DB.PORT"),
	}
	return config, nil
}
