package trainservice

import (
	"github.com/clausthrane/futfut/datasources/dsb"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/services"
	"log"
	"os"
	"time"
)

const (
	TRAIN_NUMBER = "TrainNumber"
	TRAIN_TYPE   = "TrainType"
	STATION_UIC  = "StationUic"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

type TrainService interface {
	AllTrains([]string) (*models.TrainEventList, error)
	TrainsByKeyValue(string, string) (*models.TrainEventList, error)
	TrainsFromStation(services.StationID) (*models.TrainEventList, error)
	Stops(services.TrainID) (*models.TrainEventList, error)
	DeparturesBetween(services.StationID, services.StationID) (*models.TrainEventList, error)
}

type StationProvider func(stationID services.StationID) (*models.Station, error)

type trainService struct {
	remoteAPI       dsb.DSBFacade
	stationProvider StationProvider
}

func New(remoteAPI dsb.DSBFacade, stationProvider StationProvider) TrainService {
	return &trainService{remoteAPI, stationProvider}
}

func (s *trainService) AllTrains(traintypes []string) (result *models.TrainEventList, err error) {
	if len(traintypes) > 0 {
		return s.TrainsByKeyValue(TRAIN_TYPE, traintypes[0])
	} else {
		return s.getAllTrains()
	}
}

func (s *trainService) getAllTrains() (result *models.TrainEventList, err error) {
	successChan, errChan := s.remoteAPI.GetAllTrains()
	select {
	case result := <-successChan:
		return s.joinStationName(result), nil
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Second * 30):
		return nil, services.NewServiceTimeoutError("Failed to get all trains in time")
	}
}

func (s *trainService) TrainsFromStation(stationID services.StationID) (result *models.TrainEventList, err error) {
	return s.TrainsByKeyValue(STATION_UIC, string(stationID))
}

func (s *trainService) Stops(trainID services.TrainID) (result *models.TrainEventList, err error) {
	return s.TrainsByKeyValue(TRAIN_NUMBER, string(trainID))
}

func (s *trainService) DeparturesBetween(from services.StationID, to services.StationID) (result *models.TrainEventList, err error) {
	// TODO
	return s.TrainsByKeyValue(STATION_UIC, string(from))
}

func (s *trainService) TrainsByKeyValue(key string, value string) (result *models.TrainEventList, err error) {
	successChan, errChan := s.remoteAPI.GetTrains(key, value)
	select {
	case result := <-successChan:
		return s.joinStationName(result), nil
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Second * 30):
		return nil, services.NewServiceTimeoutError("Failed to get all trains in time")
	}
}

func (s *trainService) joinStationName(list *models.TrainEventList) *models.TrainEventList {
	for idx, _ := range list.Events {
		if station, err := s.stationProvider(services.StationID(list.Events[idx].StationUic)); err == nil {
			logger.Println("getting name ")
			list.Events[idx].StationName = station.Name
		} else {
			logger.Println("not getting name ")
		}
	}
	return list
}
