package mockserver

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

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
