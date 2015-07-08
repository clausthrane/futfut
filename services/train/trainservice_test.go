package trainservice

import (
	"errors"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/services"
	"github.com/clausthrane/futfut/tests/mockdsbfacade"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

var mockStationProvider StationProvider

func init() {
	mockStationProvider = func(id services.StationID) (*models.Station, error) {
		return &models.Station{
			"short-name",
			string(id) + "-name",
			"NA",
			"NA",
			"NA",
		}, nil
	}
}

func TestRemoteAPIChannelContentIsPropagated(t *testing.T) {
	assert := assert.New(t)

	// Buf size 1 so everyting can run in 1 go routine
	fail := make(chan error, 1)
	succ := make(chan *models.TrainEventList, 1)

	remoteAPIMock := new(mockfacade.MockDSB)
	remoteAPIMock.On("GetTrains", mock.Anything, mock.Anything).Return(succ, fail)

	service := New(remoteAPIMock, mockStationProvider)

	succ <- &models.TrainEventList{}
	out, err := service.TrainsByKeyValue("a", "b")
	assert.NotNil(out)
	assert.Nil(err)

	fail <- errors.New("Aww")
	out, err = service.TrainsByKeyValue("a", "b")
	assert.NotNil(err)
	assert.Nil(out)
}

func TestJoinStationNameOnAllTrains(t *testing.T) {
	assert := assert.New(t)

	// Buf size 1 so everyting can run in 1 go routine
	fail := make(chan error, 1)
	succ := make(chan *models.TrainEventList, 1)

	remoteAPIMock := new(mockfacade.MockDSB)
	remoteAPIMock.On("GetAllTrains").Return(succ, fail)

	in := models.TrainEvent{
		"",
		"id",
		"toBeReplaced",
		"",
		"",
		0,
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		"",
		false,
		"",
		"",
		"",
	}

	succ <- &models.TrainEventList{[]models.TrainEvent{in}}

	service := New(remoteAPIMock, mockStationProvider)

	out, _ := service.AllTrains([]string{})

	assert.Equal(1, len(out.Events))
	assert.Equal("id-name", out.Events[0].StationName)
}
