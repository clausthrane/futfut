package stationservice

import (
	"errors"
	"github.com/clausthrane/futfut/datasources/dsb"
	"github.com/clausthrane/futfut/models"
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
	GetStations(countryCode string, countryName string, page int, pageSize int) (*models.StationList, error)
}

// New returns a new StationService
func New(remoteAPI dsb.DSBFacade, timeout int) StationsService {
	return &stationService{remoteAPI, timeout}
}

// NewDSBStationService returns a new StationService backed by a default DSBFacade
func NewDSBStationService() StationsService {
	return New(dsb.NewDSBFacade(), 30)
}

// AllStations returns all available stations provided by the underlying DSBFacade
// A timeout of 10 seconds is applied to the request, after which an empty list is
// returned
func (s *stationService) AllStations() (res *models.StationList, err error) {
	successChan, errChan := s.remoteAPI.GetStations()
	select {
	case res = <-successChan:
		return res, nil
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Second * 30):
		logger.Println("Timeout!")
		return nil, errors.New("Timeout loading stations")
	}
}

// GetStations returns stations based on the supplied parameters
func (s *stationService) GetStations(countryCode string, countryName string, page int, pageSize int) (*models.StationList, error) {
	return (*s).AllStations()
}

func (s *stationService) setTimeoutInSeconds(seconds int) {
	s.timeoutInSeconds = seconds
}
