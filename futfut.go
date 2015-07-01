package main

import (
	traindata "github.com/clausthrane/futfut/traindata"
	"log"
	"os"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

func main() {
	all, err := traindata.GetStations()
	if err == nil {
		logger.Println(all.Stations)
	}
}
