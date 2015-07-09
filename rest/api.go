package api

import (
	"github.com/clausthrane/futfut/datasources"
	"github.com/clausthrane/futfut/services"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

func NewAPI(requestHandler *RequestHandler) http.Handler {
	return NewAPIWithWebroot(requestHandler, "web/")
}

func NewAPIWithWebroot(r *RequestHandler, webroot string) http.Handler {
	api := mux.NewRouter()
	api.HandleFunc("/api/v1/stations", chain(errorHandler(r.HandleStationsRequest), CORS))
	api.HandleFunc("/api/v1/stations/{stationid}/details", chain(errorHandler(r.HandleStationsDetailRequest), CORS))

	api.HandleFunc("/api/v1/departures", chain(errorHandler(r.HandleAllTrainsRequest), CORS))
	api.HandleFunc("/api/v1/departures/from/{fromid}", chain(errorHandler(r.HandleDeparturesForStation), CORS))
	api.HandleFunc("/api/v1/departures/from/{fromid}/to/{toid}", chain(errorHandler(r.HandleDeparturesBetween), CORS))

	api.HandleFunc("/api/v1/trains/{trainid}", chain(errorHandler(r.HandleTrainStopInfo), CORS))

	api.HandleFunc("/status", chain(errorHandler(r.HandleStatusRequest), CORS))

	api.PathPrefix("/").Handler(http.FileServer(http.Dir(webroot)))
	return handlers.CombinedLoggingHandler(os.Stdout, api)
}

// handlerWrapper is a convinient function type for wrapping handlers
//
// https://golang.org/pkg/net/http/#HandlerFunc : func(ResponseWriter, *Request)
type handlerWrapper func(http.Handler) http.HandlerFunc

func chain(handler http.HandlerFunc, others ...handlerWrapper) http.HandlerFunc {
	for _, other := range others {
		handler = other(handler)
	}
	return handler
}

func CORS(handler http.Handler) http.HandlerFunc {
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

type FailableHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func errorHandler(handler FailableHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := handler(w, r); err != nil {
			switch err.(type) {
			case services.ServiceTimeoutError:
				http.Error(w, err.Error(), http.StatusRequestTimeout)
			case services.ServiceValidationError:
				http.Error(w, err.Error(), http.StatusBadRequest)
			case traindata.ClientError:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			case traindata.RemoteError:
				http.Error(w, err.Error(), http.StatusServiceUnavailable)
			default:
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			logger.Println(err.Error())
		}
	}
}
