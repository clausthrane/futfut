package main

import (
	"fmt"
	"html"

	"github.com/spf13/viper"
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

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.ListenAndServe(":8080", nil)
}
