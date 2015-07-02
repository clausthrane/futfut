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
	remoteAPI dsb.DSBFacade
}

// New returns a new StationService
func New(remoteAPI dsb.DSBFacade) interface{} {
	return &stationService{remoteAPI}
}

// NewDSBStationService returns a new StationService backed by a default DSBFacade
func NewDSBStationService() interface{} {
	return New(dsb.NewDSBFacade())
}

// AllStations returns all available stations provided by the underlying DSBFacade
// A timeout of 10 seconds is applied to the request, after which an empty list is
// returned
func (s *stationService) AllStations() (res *models.StationList, err error) {
	replyChan := s.remoteAPI.GetStations()
	res = &models.StationList{}
	select {
	case res = <-replyChan:
	case <-time.After(time.Second * 10):
		logger.Println("Timeout!")
		err = errors.New("Timeout loading stations")
	}
	return res, err
}

// GetStations returns stations based on the supplied parameters
func (s *stationService) GetStations(countryCode string, countryName string, page int, pageSize int) (*models.StationList, error) {
	return (*s).AllStations()
}
