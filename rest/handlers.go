package api

import (
	"encoding/json"
	"github.com/clausthrane/futfut/services"
	"github.com/clausthrane/futfut/views"
	"github.com/gorilla/mux"
	"net/http"
)

type RequestHandler struct {
	stationView views.StationView
	trainView   views.TrainView
}

func NewRequestHandler(stationView views.StationView, trainView views.TrainView) *RequestHandler {
	return &RequestHandler{stationView, trainView}
}

func (h *RequestHandler) HandleStationsRequest(w http.ResponseWriter, r *http.Request) {
	allStations, err := h.stationView.AllStations()
	if err == nil {
		logger.Printf("Listing %d stations", allStations.Count)
		encoder := json.NewEncoder(w)
		encoder.Encode(allStations)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *RequestHandler) HandleAllTrainsRequest(w http.ResponseWriter, r *http.Request) {
	allTrains, err := h.trainView.AllTrains()
	if err == nil {
		logger.Printf("Listing %d trains", allTrains.Count)
		encoder := json.NewEncoder(w)
		encoder.Encode(allTrains)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *RequestHandler) HandleDeparturesBetween(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	from := params["fromid"]
	to := params["toid"]

	trains, err := h.trainView.DeparturesBetween(services.StationID(from), services.StationID(to))
	if err == nil {
		logger.Printf("Listing %d trains", trains.Count)
		encoder := json.NewEncoder(w)
		encoder.Encode(trains)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
