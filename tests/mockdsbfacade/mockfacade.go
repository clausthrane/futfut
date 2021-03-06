package mockfacade

import (
	"github.com/clausthrane/futfut/models"
	"github.com/stretchr/testify/mock"
)

type MockDSB struct {
	mock.Mock
}

func (m MockDSB) GetStations() (chan *models.StationList, chan error) {
	args := m.Called()
	return args.Get(0).(chan *models.StationList), args.Get(1).(chan error)
}

func (m MockDSB) GetStation(stationId string) (chan *models.StationList, chan error) {
	args := m.Called(stationId)
	return args.Get(0).(chan *models.StationList), args.Get(1).(chan error)
}

func (m MockDSB) GetAllTrains() (chan *models.TrainEventList, chan error) {
	args := m.Called()
	return args.Get(0).(chan *models.TrainEventList), args.Get(1).(chan error)
}

func (m MockDSB) GetTrains(key string, value string) (chan *models.TrainEventList, chan error) {
	args := m.Called(key, value)
	return args.Get(0).(chan *models.TrainEventList), args.Get(1).(chan error)
}
