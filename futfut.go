package main

import (
	"github.com/spf13/viper"
	"log"
	"os"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal("No configuration file loaded - using defaults")
	}
}
