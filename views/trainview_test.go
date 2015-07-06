package views

import (
	"errors"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/services"
	"github.com/clausthrane/futfut/views/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestTrainViewPanicsWhenNotInitialized(t *testing.T) {
	assert := assert.New(t)

	out := NewTrainView(nil, nil)

	assert.Panics(nil, func() {
		out.AllTrains()
	}, "without service and convereter we should panic")
}

func TestTrainView(t *testing.T) {
	assert := assert.New(t)

	emptyList := &models.TrainList{nil}
	service := new(mockTrainService)
	service.On("AllTrains").Return(emptyList, nil)

	dtos := &dto.JSONTrainList{1, []dto.JSONTrain{{42, 10, "dest", "", ""}}}
	converter := new(mockTrainConverter)
	converter.On("ConvertTrainList", emptyList).Return(dtos)

	view := NewTrainView(service, converter)
	out, err := view.AllTrains()

	assert.Nil(err)
	assert.Equal(1, out.Count)
	assert.Equal(42, out.Trains[0].TrainNumber)
}

func TestTrainViewConvertionNotCalledOnServiceError(t *testing.T) {
	assert := assert.New(t)

	service := new(mockTrainService)
	service.On("AllTrains").Return(nil, errors.New("fake error"))

	converter := new(mockTrainConverter)

	view := NewTrainView(service, converter)
	out, err := view.AllTrains()

	assert.NotNil(err)
	assert.Nil(out)

	converter.AssertNotCalled(t, "ConvertTrains", mock.Anything)
}

func TestTrainViewServiceErrorsArePropagated(t *testing.T) {
	assert := assert.New(t)

	service := new(mockTrainService)
	service.On("AllTrains").Return(nil, errors.New("fake error"))

	converter := new(mockTrainConverter)

	view := NewTrainView(service, converter)
	out, err := view.AllTrains()

	assert.NotNil(err)
	assert.Nil(out)
	assert.Equal("fake error", err.Error())
}

type mockTrainService struct {
	mock.Mock
}

func (m *mockTrainService) AllTrains() (result *models.TrainList, err error) {
	args := m.Called()
	e := args.Get(1)
	if e != nil {
		return nil, e.(error)
	} else {
		return args.Get(0).(*models.TrainList), nil
	}
}

func (m *mockTrainService) TrainsByKeyValue(key string, value string) (result *models.TrainList, err error) {
	args := m.Called(key, value)
	e := args.Get(1)
	if e != nil {
		return nil, e.(error)
	} else {
		return args.Get(0).(*models.TrainList), nil
	}
}

func (m *mockTrainService) TrainsFromStation(stationID services.StationID) (result *models.TrainList, err error) {
	args := m.Called(stationID)
	e := args.Get(1)
	if e != nil {
		return nil, e.(error)
	} else {
		return args.Get(0).(*models.TrainList), nil
	}
}

func (m *mockTrainService) Stops(trainID services.TrainID) (result *models.TrainList, err error) {
	args := m.Called(trainID)
	e := args.Get(1)
	if e != nil {
		return nil, e.(error)
	} else {
		return args.Get(0).(*models.TrainList), nil
	}
}

func (m *mockTrainService) DeparturesBetween(from services.StationID, to services.StationID) (result *models.TrainList, err error) {
	args := m.Called(from, to)
	e := args.Get(1)
	if e != nil {
		return nil, e.(error)
	} else {
		return args.Get(0).(*models.TrainList), nil
	}
}

type mockTrainConverter struct {
	mock.Mock
}

func (m mockTrainConverter) ConvertTrainList(list *models.TrainList) *dto.JSONTrainList {
	args := m.Called(list)
	return args.Get(0).(*dto.JSONTrainList)
}

func (m mockTrainConverter) ConvertTrain(t *models.Train) (*dto.JSONTrain, error) {
	args := m.Called(t)
	e := args.Get(1)
	if e != nil {
		return nil, e.(error)
	} else {
		return args.Get(0).(*dto.JSONTrain), nil
	}
}
