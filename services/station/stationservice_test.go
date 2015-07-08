package stationservice

import (
	"errors"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/tests/mockdsbfacade"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRemoteAPIChannelContentIsPropagated(t *testing.T) {
	assert := assert.New(t)

	// Buf size 1 so everyting can run in 1 go routine
	succ := make(chan *models.StationList, 1)
	fail := make(chan error, 1)

	remoteAPIMock := new(mockfacade.MockDSB)
	remoteAPIMock.On("GetStations").Return(succ, fail)

	service := New(remoteAPIMock)

	succ <- &models.StationList{}
	out, err := service.AllStations()
	assert.NotNil(out, "expecting output")
	assert.Nil(err, "errors not expected")

	fail <- errors.New("Aww")
	out, err = service.AllStations()
	assert.NotNil(err, "expecting err")
	assert.Nil(out, "output not expected")
}
