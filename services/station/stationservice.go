package stationservice

import (
	"errors"
	"fmt"
	"github.com/clausthrane/futfut/datasources/dsb"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/services"
	"log"
	"os"
	"time"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

type stationService struct {
	remoteAPI        dsb.DSBFacade
	timeoutInSeconds int
}

type StationsService interface {
	AllStations() (*models.StationList, error)
	Station(services.StationID) (*models.Station, error)
}

// New returns a new StationService
func New(remoteAPI dsb.DSBFacade) StationsService {
	return &stationService{remoteAPI, 0}
}

// NewDSBStationService returns a new StationService backed by a default DSBFacade
func NewDSBStationService() StationsService {
	return New(dsb.NewDSBFacade())
}

// AllStations returns all available stations provided by the underlying DSBFacade
// A timeout of 29 seconds (normal HTTP timeout is 30) is applied to the request
func (s *stationService) AllStations() (res *models.StationList, err error) {
	successChan, errChan := s.remoteAPI.GetStations()
	select {
	case res = <-successChan:
		return res, nil
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Second * 29):
		logger.Println("Timeout!")
		return nil, errors.New("Timeout loading stations")
	}
}

func (s *stationService) Station(stationID services.StationID) (*models.Station, error) {
	if all, err := s.AllStations(); err == nil {
		for _, s := range all.Stations {
			if services.StationID(s.UIC) == stationID {
				return &s, nil
			}
		}
	}
	return nil, services.NewServiceValidationError(fmt.Sprintf("Unknown station id: %", string(stationID)))
}

func (s *stationService) setTimeoutInSeconds(seconds int) {
	s.timeoutInSeconds = seconds
}
