package main

import (
	"github.com/clausthrane/futfut/config"
	"github.com/clausthrane/futfut/rest"
	"github.com/clausthrane/futfut/services/station"
	"github.com/clausthrane/futfut/views"
	"github.com/clausthrane/futfut/views/dto"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	port := config.GetString("host.port")
	logger.Printf("Starting server on: %s", port)
	http.ListenAndServe(port, webApp())
}

func webApp() http.Handler {
	stationService := stationservice.NewDSBStationService()
	stationConverter := dto.NewStationConverter()
	view := views.NewStationsView(stationService, stationConverter)
	requestHandler := api.NewRequestHandler(view)
	return api.NewAPI(requestHandler)
}
