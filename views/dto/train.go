package dto

import (
	"github.com/clausthrane/futfut/models"
	"log"
	"os"
	"strconv"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

// JSONTrain is the external representation of models.Train
type JSONTrain struct {
	TrainNumber      int
	CurrentStationId int
	DestinationName  string
}

// JSONTrainList is the external representation of models.TrainList
type JSONTrainList struct {
	Count  int
	Trains []JSONTrain
}

type trainConverter struct {
}

// TrainConverter maps models.TrainList to the external representation
type TrainConverter interface {
	ConvertTrainList(*models.TrainList) *JSONTrainList
	ConvertTrain(*models.Train) (*JSONTrain, error)
}

func NewTrainConverter() TrainConverter {
	return &trainConverter{}
}

func (c trainConverter) ConvertTrainList(list *models.TrainList) *JSONTrainList {
	dtos := []JSONTrain{}
	for _, t := range list.Trains {
		if dto, err := c.ConvertTrain(&t); err == nil {
			dtos = append(dtos, *dto)
		}
	}
	return &JSONTrainList{len(dtos), dtos}
}

func (c trainConverter) ConvertTrain(t *models.Train) (*JSONTrain, error) {
	trainNumner, err := strconv.Atoi(t.TrainNumber)
	currentStationId, err := strconv.Atoi(t.StationUic)
	if err != nil {
		logger.Printf("Failed to covert %s, due to %s", t, err.Error())
		return nil, err
	} else {
		return &JSONTrain{
			trainNumner,
			currentStationId,
			t.DestinationName,
		}, nil
	}
}
