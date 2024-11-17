package utils

import (
	"github.com/spf13/viper"
	"log"
)

func InitConfig() {
	viper.AddConfigPath("..")
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
		panic(err)
	}
	viper.AutomaticEnv()
}
