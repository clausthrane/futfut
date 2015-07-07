// Package mockserver implement a collection of test harneses to be used in automated tests

package mockserver

import (
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

// HttpServer starts a HTTP server on the provided port, where requests
// are handled as specified by the hander
func HttpServer(port int, handler http.HandlerFunc) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logger.Fatal("Unable to start mock server")
	}
	go http.Serve(listener, handler)
}

// HttpServerWithStatusCode starts a HTTP server which always replies
// with the given status code
func HttpServerWithStatusCode(port int, statusCode int) {
	HttpServer(port, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Mocking", statusCode)
	})
}

// HttpServerFromFiles stats a HTTP server pretending to be the DSP API
//
// The sever answers on /Station() with the content of the fiel stationsFilename
// and on /Queue() - ignoring queryparames - with the content of the file queueFilename
func HttpServerFromFiles(stationsFilename string, queueFilename string, port int) {

	router := mux.NewRouter()
	router.HandleFunc("/Station()", func(w http.ResponseWriter, r *http.Request) {
		logger.Println("Serving stations")

		stations, err := os.Open(stationsFilename)
		if err != nil {
			logger.Println("no stations file found")
		}
		w.Header().Set("Content-Type", "Application/JSON; charset=utf-8")
		io.Copy(w, stations)
	})
	router.HandleFunc("/Queue()", func(w http.ResponseWriter, r *http.Request) {
		logger.Println("Serving queue")
		queue, err := os.Open(queueFilename)
		if err != nil {
			logger.Println("no queue file found")
		}
		w.Header().Set("Content-Type", "Application/JSON; charset=utf-8")
		io.Copy(w, queue)
	})

	go func() { logger.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router)) }()
	logger.Printf("Started server on port %d", port)
}

func HttpServerDSBTestApi(t *testing.T, port int) {
	stations := "/Users/thrane/Dev/other/GO/src/github.com/clausthrane/futfut/tests/mockserver/stations.txt"
	if _, err := os.Open(stations); err != nil {
		t.Skip("skipping test; content for mock server not present")
	}

	queue := "/Users/thrane/Dev/other/GO/src/github.com/clausthrane/futfut/tests/mockserver/queue.txt"
	if _, err := os.Open(queue); err != nil {
		t.Skip("skipping test; content for mock server not present")
	}

	HttpServerFromFiles(stations, queue, port)
}
