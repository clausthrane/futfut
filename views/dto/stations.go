//Pagage dto implements external representations and converters of internal models
package dto

import (
	"github.com/clausthrane/futfut/models"
)

// JSONStation is the external representation of models.Station
type JSONStation struct {
	Id          string
	Name        string
	CountryCode string
}

// JSONStationList is the external representation of models.StationList
type JSONStationList struct {
	Count    int
	Stations []JSONStation
}

type stationConverter struct{}

type StationConverter interface {
	ConvertStationList(list *models.StationList) *JSONStationList
	ConvertStation(station *models.Station) JSONStation
}

func NewStationConverter() StationConverter {
	return &stationConverter{}
}

//ConvertStationList returns a external representation of list
func (c stationConverter) ConvertStationList(list *models.StationList) *JSONStationList {
	dtos := []JSONStation{}
	for _, s := range list.Stations {
		dtos = append(dtos, c.ConvertStation(&s))
	}

	return &JSONStationList{len(list.Stations), dtos}
}

//ConvertStation returns a external representation of station
func (stationConverter) ConvertStation(station *models.Station) JSONStation {
	return JSONStation{
		station.UIC,
		station.Name,
		station.CountryName,
	}
}
