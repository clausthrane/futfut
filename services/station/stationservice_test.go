package stationservice

import (
	"errors"
	"github.com/clausthrane/futfut/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSuccessfulRemoteAPIChannelContentIsPropagated(t *testing.T) {
	assert := assert.New(t)

	// Buf size 1 so everyting can be sync
	succ := make(chan *models.StationList, 1)
	fail := make(chan error, 1)

	remoteAPIMock := new(mockDSB)
	remoteAPIMock.On("GetStations").Return(succ, fail)

	service := New(remoteAPIMock, 1)

	succ <- &models.StationList{}
	out, err := service.AllStations()
	assert.NotNil(out, "expecting output")
	assert.Nil(err, "errors not expected")

	fail <- errors.New("Aww")
	out, err = service.AllStations()
	assert.NotNil(err, "expecting err")
	assert.Nil(out, "output not expected")
}

type mockDSB struct {
	mock.Mock
}

func (m mockDSB) GetStations() (chan *models.StationList, chan error) {
	args := m.Called()
	return args.Get(0).(chan *models.StationList), args.Get(1).(chan error)
}

func (m mockDSB) GetTrains() (chan *models.TrainList, chan error) {
	args := m.Called()
	return args.Get(0).(chan *models.TrainList), args.Get(1).(chan error)
}
