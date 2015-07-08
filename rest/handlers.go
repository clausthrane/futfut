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

func writeJsonResponse(w http.ResponseWriter, obj interface{}, err error) error {
	if err == nil {
		encoder := json.NewEncoder(w)
		encoder.Encode(obj)
		return nil
	}
	return err
}

func (h *RequestHandler) HandleStationsDetailRequest(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	stationid := params["stationid"]
	station, err := h.stationView.GetStation(services.StationID(stationid))
	return writeJsonResponse(w, station, err)
}

func (h *RequestHandler) HandleStationsRequest(w http.ResponseWriter, r *http.Request) error {
	allStations, err := h.stationView.AllStations()
	return writeJsonResponse(w, allStations, err)
}

func (h *RequestHandler) HandleAllTrainsRequest(w http.ResponseWriter, r *http.Request) error {
	values := r.URL.Query()
	typeSelection := values["traintype"]
	allTrains, err := h.trainView.AllTrains(typeSelection)
	return writeJsonResponse(w, allTrains, err)
}

func (h *RequestHandler) HandleDeparturesForStation(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	from := params["fromid"]
	logger.Println("Handling request for departure times at %s", from)
	trains, err := h.trainView.TrainsFromStation(services.StationID(from))
	return writeJsonResponse(w, trains, err)
}

func (h *RequestHandler) HandleTrainStopInfo(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	trainid := params["trainid"]
	trains, err := h.trainView.Stops(services.TrainID(trainid))
	return writeJsonResponse(w, trains, err)
}

func (h *RequestHandler) HandleDeparturesBetween(w http.ResponseWriter, r *http.Request) error {
	params := mux.Vars(r)
	from := params["fromid"]
	to := params["toid"]
	trains, err := h.trainView.DeparturesBetween(services.StationID(from), services.StationID(to))
	return writeJsonResponse(w, trains, err)
}
