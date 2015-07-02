package api

import (
	"encoding/json"
	"github.com/clausthrane/futfut/views"
	"net/http"
)

type RequestHandler struct {
	stationView views.StationView
}

func NewRequestHandler(stationView views.StationView) *RequestHandler {
	return &RequestHandler{stationView}
}

func (h *RequestHandler) HandleStationsRequest(w http.ResponseWriter, r *http.Request) {
	allStations, err := h.stationView.AllStations()
	if err == nil {
		logger.Println("Listing stations", len(allStations.Stations))
		encoder := json.NewEncoder(w)
		encoder.Encode(allStations.Stations)
	}

}
