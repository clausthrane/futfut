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

func (h *RequestHandler) HandleStationsDetailRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stationid := params["stationid"]
	station, err := h.stationView.GetStation(services.StationID(stationid))
	if err == nil {
		encoder := json.NewEncoder(w)
		encoder.Encode(station)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
	values := r.URL.Query()
	typeSelection := values["traintype"]
	allTrains, err := h.trainView.AllTrains(typeSelection)
	if err == nil {
		logger.Printf("Listing %d trains", allTrains.Count)
		encoder := json.NewEncoder(w)
		encoder.Encode(allTrains)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *RequestHandler) HandleDeparturesForStation(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	from := params["fromid"]
	logger.Println("Handling request for departure times at %s", from)
	trains, err := h.trainView.TrainsFromStation(services.StationID(from))
	if err == nil {
		encoder := json.NewEncoder(w)
		encoder.Encode(trains)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *RequestHandler) HandleTrainStopInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	trainid := params["trainid"]
	trains, err := h.trainView.Stops(services.TrainID(trainid))
	if err == nil {
		encoder := json.NewEncoder(w)
		encoder.Encode(trains)
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
		encoder := json.NewEncoder(w)
		encoder.Encode(trains)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
