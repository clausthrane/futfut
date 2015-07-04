package main

import (
	"github.com/clausthrane/futfut/config"
	"github.com/clausthrane/futfut/datasources/dsb"
	"github.com/clausthrane/futfut/rest"
	"github.com/clausthrane/futfut/services/station"
	"github.com/clausthrane/futfut/services/train"
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

	dsbFacade := dsb.NewDSBFacade()

	stationService := stationservice.New(dsbFacade)
	stationConverter := dto.NewStationConverter()
	stationView := views.NewStationsView(stationService, stationConverter)

	trainService := trainservice.New(dsbFacade)
	trainConverter := dto.NewTrainConverter()
	trainView := views.NewTrainView(trainService, trainConverter)

	requestHandler := api.NewRequestHandler(stationView, trainView)
	return api.NewAPI(requestHandler)
}
