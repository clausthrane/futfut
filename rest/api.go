package api

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

type handlerWrapper func(http.Handler) http.HandlerFunc

func NewAPI(requestHandler *RequestHandler) http.Handler {
	api := mux.NewRouter()
	api.HandleFunc("/api/stations", chainHandlers(requestHandler.HandleStationsRequest, allowCORS))
	api.HandleFunc("/api/stations/{stationid}/details", chainHandlers(requestHandler.HandleStationsDetailRequest, allowCORS))

	api.HandleFunc("/api/departures", chainHandlers(requestHandler.HandleAllTrainsRequest, allowCORS))
	api.HandleFunc("/api/departures/from/{fromid}", chainHandlers(requestHandler.HandleDeparturesBetween, allowCORS))
	api.HandleFunc("/api/departures/from/{fromid}/to/{toid}", chainHandlers(requestHandler.HandleDeparturesBetween, allowCORS))

	api.HandleFunc("/api/trains/{trainid}", chainHandlers(requestHandler.HandleTrainStopInfo, allowCORS))

	api.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))
	return handlers.CombinedLoggingHandler(os.Stdout, api)
}

func chainHandlers(handler http.HandlerFunc, others ...handlerWrapper) http.HandlerFunc {
	for _, other := range others {
		handler = other(handler)
	}
	return handler
}

func allowCORS(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
			w.Header().Set("Access-Control-Max-Age", "1000")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
		} else {
			handler.ServeHTTP(w, r)
		}
	}
}
