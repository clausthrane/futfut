package main

import (
	"github.com/clausthrane/futfut/rest"
	"github.com/clausthrane/futfut/services"
	"github.com/clausthrane/futfut/views"
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
	http.ListenAndServe(viper.GetString("host.port"), webApp())
}

func webApp() http.Handler {
	stationsService := services.NewStationsService()
	view := views.NewStationsView(stationsService)
	requestHandler := api.NewRequestHandler(view)
	return api.NewAPI(requestHandler)
}
