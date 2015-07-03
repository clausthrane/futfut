package views

import (
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/services/station"
	"github.com/clausthrane/futfut/views/dto"
)

type StationView interface {
	AllStations() (*dto.JSONStationList, error)
	GetStations(countryCode string, countryName string, page int, pageSize int) (*dto.JSONStationList, error)
}

// NewStationsView provides a new view for the stations resource
func NewStationsView(service stationservice.StationsService, converter dto.StationConverter) StationView {
	return &stationView{service, converter}
}

type stationView struct {
	service   stationservice.StationsService
	converter dto.StationConverter
}

func (view *stationView) AllStations() (list *dto.JSONStationList, err error) {
	var mlist *models.StationList
	if mlist, err = view.service.AllStations(); err == nil {
		return view.converter.ConvertStationList(mlist), nil
	}
	return nil, err
}

func (view *stationView) GetStations(countryCode string, countryName string, page int, pageSize int) (*dto.JSONStationList, error) {
	return view.AllStations()
}
