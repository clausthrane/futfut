package api

import (
	"fmt"
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

	api.HandleFunc("/api/trains/from/{fromid}/to/{toid}/", chainHandlers(requestHandler.HandleDeparturesBetween, allowCORS))
	api.HandleFunc("/api/trains", chainHandlers(requestHandler.HandleAllTrainsRequest, allowCORS))

	api.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))
	return api
}

func chainHandlers(handler http.HandlerFunc, others ...handlerWrapper) http.HandlerFunc {
	for _, other := range others {
		handler = other(handler)
	}
	return handler
}

func allowCORS(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Println(fmt.Sprintf("Got request from: %s", r.RemoteAddr))
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
