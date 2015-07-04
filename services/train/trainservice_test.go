package trainservice

import (
	"errors"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/tests/mockdsbfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestRemoteAPIChannelContentIsPropagated(t *testing.T) {
	assert := assert.New(t)

	// Buf size 1 so everyting can run in 1 go routine
	fail := make(chan error, 1)
	succ := make(chan *models.TrainList, 1)

	remoteAPIMock := new(mockfacade.MockDSB)
	remoteAPIMock.On("GetTrains", mock.Anything, mock.Anything).Return(succ, fail)

	service := New(remoteAPIMock)

	succ <- &models.TrainList{}
	out, err := service.AllTrains()
	assert.NotNil(out, "expecting output")
	assert.Nil(err, "errors not expected")

	fail <- errors.New("Aww")
	out, err = service.AllTrains()
	assert.NotNil(err, "expecting err")
	assert.Nil(out, "output not expected")
}
