package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

type GoogleOauthConfig struct {
	ClientID     string
	ClientSecret string
}

type Config struct {
	ServerPort   string
	BaseURL      string
	JWT_SECRET   string
	GoogleOauth  GoogleOauthConfig
	SMPHost      string
	SMTPPort     int
	SMTPEmail    string
	SMTPPassword string
	BASEURL      string
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
		GoogleOauth: GoogleOauthConfig{
			ClientID:     viper.GetString("GOOGLE_CLIENT_ID"),
			ClientSecret: viper.GetString("GOOGLE_CLIENT_SECRET"),
		},
		SMPHost:      viper.GetString("SMTP_HOST"),
		SMTPPort:     viper.GetInt("SMTP_PORT"),
		SMTPEmail:    viper.GetString("SMTP_EMAIL"),
		SMTPPassword: viper.GetString("SMTP_PASSWORD"),
		BASEURL:      viper.GetString("BASE_URL"),
	}
	return config, nil
}
