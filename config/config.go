package config

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal("No configuration file loaded - exiting")
	}
}

func GetString(key string) string {
	return viper.GetString(key)
}
