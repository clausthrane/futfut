package trainservice

import (
	"github.com/clausthrane/futfut/datasources/dsb"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/services"
	"log"
	"os"
	"time"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Ltime|log.Lshortfile)

type trainService struct {
	remoteAPI dsb.DSBFacade
}

type TrainService interface {
	AllTrains() (*models.TrainList, error)
}

func New(remoteAPI dsb.DSBFacade) TrainService {
	return &trainService{remoteAPI}
}

func (s *trainService) AllTrains() (result *models.TrainList, err error) {
	successChan, errChan := s.remoteAPI.GetTrains("TrainType", "S-tog")
	select {
	case result := <-successChan:
		return result, nil
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Second * 30):
		return nil, services.NewServiceTimeoutError("Failed to get all trains in time")
	}
}
