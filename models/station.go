package models

import (
	"bytes"
)

type Station struct {
	Abbreviation string
	Name         string
	UIC          string
	CountryCode  string
	CountryName  string
}

type StationList struct {
	Stations []Station
}

func (list *StationList) Len() int {
	return len(list.Stations)
}

func (list *StationList) Less(i int, j int) bool {
	a := list.Stations[i].Name
	b := list.Stations[j].Name
	return bytes.Compare([]byte(a), []byte(b)) < 0
}

func (list *StationList) Swap(i int, j int) {
	list.Stations[i], list.Stations[j] = list.Stations[j], list.Stations[i]
}
