package config

import (
	"github.com/spf13/viper"
	"log"
)

func InitConfig(fileName string) *viper.Viper {
	config := viper.New()
	config.SetConfigName(fileName)
	config.AddConfigPath(".")
	config.AddConfigPath("$HOME")
	if err := config.ReadInConfig(); err != nil {
		log.Fatal("Error while parsing configuration file", err)
	}
	return config

}
