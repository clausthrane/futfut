package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"html"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal("No configuration file loaded - using defaults")
	}

	router := mux.NewRouter()
	router.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})
	router.HandleFunc("/baz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Go away, %q", html.EscapeString(r.URL.Path))
	})

	port := viper.GetString("host.port")
	logger.Printf("Starting server on: %s", port)
	http.ListenAndServe(viper.GetString("host.port"), router)
}
