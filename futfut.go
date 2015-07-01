package main

import (
	"github.com/clausthrane/futfut/api"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal("No configuration file loaded - using defaults")
	}

	port := viper.GetString("host.port")
	logger.Printf("Starting server on: %s", port)
	http.ListenAndServe(viper.GetString("host.port"), api.NewHandler())
}
