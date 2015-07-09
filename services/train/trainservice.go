package trainservice

import (
	"github.com/clausthrane/futfut/datasources/dsb"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/services"
	"github.com/clausthrane/futfut/utils/algs"
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
	logger.Printf("Computing route from %v to %v", from, to)

	if world, err := s.getAllTrains(); err == nil {
		if departures, departuresErr := s.TrainsFromStation(from); departuresErr == nil {
			for departureCandidateIdx, _ := range departures.Events {
				if path, err := graph.Dijkstra(world, &departures.Events[departureCandidateIdx], to); err == nil {
					return path, nil
				}
			}
			return nil, services.NewServiceValidationError("Cannot find any route")
		} else {
			return nil, departuresErr
		}
	} else {
		return nil, err
	}
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
			list.Events[idx].StationName = station.Name
		}
	}
	return list
}
