package api

import (
	"encoding/json"
	"github.com/clausthrane/futfut/views"
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

func (h *RequestHandler) HandleTrainsRequest(w http.ResponseWriter, r *http.Request) {
	allTrains, err := h.trainView.AllTrains()
	if err == nil {
		logger.Printf("Listing %d trains", allTrains.Count)
		encoder := json.NewEncoder(w)
		encoder.Encode(allTrains)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
