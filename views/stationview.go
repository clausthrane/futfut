package views

import (
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/services"
)

type StationView interface {
	AllStations() (*models.StationList, error)
	GetStations(countryCode string, countryName string, page int, pageSize int) (*models.StationList, error)
}

func NewStationsView(service services.StationsService) StationView {
	return &stationView{service}
}

type stationView struct {
	service services.StationsService
}

func (view *stationView) AllStations() (*models.StationList, error) {
	return view.service.AllStations()
}

func (view *stationView) GetStations(countryCode string, countryName string, page int, pageSize int) (*models.StationList, error) {
	return view.service.AllStations()
}
