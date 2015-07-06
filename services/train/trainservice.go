package trainservice

import (
	"fmt"
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
	AllTrains() (*models.TrainList, error)
	TrainsByKeyValue(string, string) (*models.TrainList, error)
	TrainsFromStation(services.StationID) (*models.TrainList, error)
	Stops(services.TrainID) (*models.TrainList, error)
	DeparturesBetween(services.StationID, services.StationID) (*models.TrainList, error)
}

type trainService struct {
	remoteAPI dsb.DSBFacade
}

func New(remoteAPI dsb.DSBFacade) TrainService {
	return &trainService{remoteAPI}
}

func (s *trainService) AllTrains() (result *models.TrainList, err error) {
	// {"IC", "RE", "S-tog"}
	// TODO: dont restrict to S-tog but rely on caching
	return s.TrainsByKeyValue(TRAIN_TYPE, "S-tog")
}

func (s *trainService) TrainsFromStation(stationID services.StationID) (result *models.TrainList, err error) {
	return s.TrainsByKeyValue(STATION_UIC, string(stationID))
}

func (s *trainService) Stops(trainID services.TrainID) (result *models.TrainList, err error) {
	return s.TrainsByKeyValue(TRAIN_NUMBER, string(trainID))
}

func (s *trainService) DeparturesBetween(from services.StationID, to services.StationID) (result *models.TrainList, err error) {

	world, _ := s.TrainsByKeyValue("", "")
	logger.Println(fmt.Sprintf("worlds has: %d entries", len(world.Trains)))

	logger.Println(graph.Dijkstra(world, from, to))

	return s.TrainsByKeyValue(STATION_UIC, string(from))
}

func (s *trainService) TrainsByKeyValue(key string, value string) (result *models.TrainList, err error) {
	successChan, errChan := s.remoteAPI.GetTrains(key, value)
	select {
	case result := <-successChan:
		return result, nil
	case err = <-errChan:
		return nil, err
	case <-time.After(time.Second * 30):
		return nil, services.NewServiceTimeoutError("Failed to get all trains in time")
	}
}
