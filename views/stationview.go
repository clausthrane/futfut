package views

import (
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/services"
	"github.com/clausthrane/futfut/services/station"
	"github.com/clausthrane/futfut/views/dto"
)

type StationView interface {
	AllStations() (*dto.JSONStationList, error)
	GetStation(services.StationID) (dto.JSONStation, error)
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

func (view *stationView) GetStation(stationID services.StationID) (dto.JSONStation, error) {
	if station, err := view.service.Station(stationID); err == nil {
		return view.converter.ConvertStation(station), nil
	} else {
		return dto.JSONStation{}, err
	}
}
