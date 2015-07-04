package views

import (
	"github.com/clausthrane/futfut/services/train"
	"github.com/clausthrane/futfut/views/dto"
)

type TrainView interface {
	AllTrains() (*dto.JSONTrainList, error)
}

func NewTrainView(service trainservice.TrainService, converter dto.TrainConverter) TrainView {
	return &trainView{service, converter}
}

type trainView struct {
	service   trainservice.TrainService
	converter dto.TrainConverter
}

func (view *trainView) AllTrains() (*dto.JSONTrainList, error) {
	if modelList, err := view.service.AllTrains(); err == nil {
		return view.converter.ConvertTrainList(modelList), nil
	} else {
		return nil, err
	}
}
