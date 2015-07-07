package models

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
