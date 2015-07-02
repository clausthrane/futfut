package services

import (
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/services/station"
)

type StationsService interface {
	AllStations() (*models.StationList, error)
	GetStations(countryCode string, countryName string, page int, pageSize int) (*models.StationList, error)
}

func NewStationsService() StationsService {
	return stationservice.NewDSBStationService().(StationsService)
}
