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

func (m MockDSB) GetTrains() (chan *models.TrainList, chan error) {
	args := m.Called()
	return args.Get(0).(chan *models.TrainList), args.Get(1).(chan error)
}
