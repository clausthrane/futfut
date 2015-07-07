package views

import (
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/services"
	"github.com/clausthrane/futfut/services/train"
	"github.com/clausthrane/futfut/views/dto"
)

type TrainView interface {
	AllTrains([]string) (*dto.JSONTrainEventList, error)
	DeparturesBetween(services.StationID, services.StationID) (*dto.JSONTrainEventList, error)
	TrainsFromStation(services.StationID) (*dto.JSONTrainEventList, error)
	Stops(services.TrainID) (*dto.JSONTrainEventList, error)
}

func NewTrainView(service trainservice.TrainService, converter dto.TrainConverter) TrainView {
	return &trainView{service, converter}
}

type trainView struct {
	service   trainservice.TrainService
	converter dto.TrainConverter
}

func (view *trainView) AllTrains(traintypes []string) (*dto.JSONTrainEventList, error) {
	if modelList, err := view.service.AllTrains(traintypes); err == nil {
		return view.converter.ConvertTrainList(modelList), nil
	} else {
		return nil, err
	}
}

func (view *trainView) Stops(trainId services.TrainID) (*dto.JSONTrainEventList, error) {
	res, err := view.service.Stops(trainId)
	return convertTrainList(view.converter, res, err)
}

func (view *trainView) TrainsFromStation(from services.StationID) (*dto.JSONTrainEventList, error) {
	res, err := view.service.TrainsFromStation(from)
	return convertTrainList(view.converter, res, err)
}

func (view *trainView) DeparturesBetween(from services.StationID, to services.StationID) (*dto.JSONTrainEventList, error) {
	if modelList, err := view.service.DeparturesBetween(from, to); err == nil {
		return view.converter.ConvertTrainList(modelList), nil
	} else {
		return nil, err
	}
}

func convertTrainList(converter dto.TrainConverter, list *models.TrainEventList, err error) (*dto.JSONTrainEventList, error) {
	if err == nil {
		return converter.ConvertTrainList(list), nil
	} else {
		return nil, err
	}
}
