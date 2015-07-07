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
	succ := make(chan *models.TrainEventList, 1)

	remoteAPIMock := new(mockfacade.MockDSB)
	remoteAPIMock.On("GetTrains", mock.Anything, mock.Anything).Return(succ, fail)

	service := New(remoteAPIMock)

	succ <- &models.TrainEventList{}
	out, err := service.TrainsByKeyValue("a", "b")
	assert.NotNil(out)
	assert.Nil(err)

	fail <- errors.New("Aww")
	out, err = service.TrainsByKeyValue("a", "b")
	assert.NotNil(err)
	assert.Nil(out)
}
