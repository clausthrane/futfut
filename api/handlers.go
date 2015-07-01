package api

import (
	"encoding/json"
	traindata "github.com/clausthrane/futfut/traindata"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

func NewHandler() http.Handler {
	router := mux.NewRouter()
	router.HandleFunc("/api/stations", stationsHandler())
	return router

}

func stationsHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		allStations, err := traindata.GetStations()
		if err == nil {
			logger.Println("Listing stations", len(allStations.Stations))
			encoder := json.NewEncoder(w)
			encoder.Encode(allStations.Stations)
		}
	}
}
