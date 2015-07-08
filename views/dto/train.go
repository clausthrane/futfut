package dto

import (
	"github.com/clausthrane/futfut/models"
	"log"
	"os"
	"strconv"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

// JSONTrain is the external representation of models.TrainEvent
type JSONTrainEvent struct {
	TrainNumber     int
	StationId       int
	StationName     string
	DestinationName string
	ArrivalTime     string
	ArrivalDate     string
	DepartureTime   string
	DepartureDate   string
}

// JSONTrainList is the external representation of models.TrainList
type JSONTrainEventList struct {
	Count  int
	Trains []JSONTrainEvent
}

type trainConverter struct {
}

// TrainConverter maps models.TrainList to the external representation
type TrainConverter interface {
	ConvertTrainList(*models.TrainEventList) *JSONTrainEventList
	ConvertTrain(*models.TrainEvent) (*JSONTrainEvent, error)
}

func NewTrainConverter() TrainConverter {
	return &trainConverter{}
}

func (c trainConverter) ConvertTrainList(list *models.TrainEventList) *JSONTrainEventList {
	dtos := []JSONTrainEvent{}
	for _, t := range list.Events {
		if dto, err := c.ConvertTrain(&t); err == nil {
			dtos = append(dtos, *dto)
		}
	}
	return &JSONTrainEventList{len(dtos), dtos}
}

func (c trainConverter) ConvertTrain(t *models.TrainEvent) (*JSONTrainEvent, error) {
	trainNumner, err := strconv.Atoi(t.TrainNumber)
	currentStationId, err := strconv.Atoi(t.StationUic)
	if err != nil {
		return nil, err
	}

	arrivalDate := t.ScheduledArrivalDate().Format("2006-01-02")
	arrivalTime := t.ScheduledArrivalDate().Format("15:04")

	departureDate := t.ScheduledDepartureDate().Format("2006-01-02")
	departureTime := t.ScheduledDepartureDate().Format("15:04")

	return &JSONTrainEvent{
		trainNumner,
		currentStationId,
		"NA",
		t.DestinationName,
		arrivalTime,
		arrivalDate,
		departureTime,
		departureDate,
	}, nil
}
