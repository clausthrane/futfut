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
		out.AllTrains([]string{"anytype"})
	}, "without service and convereter we should panic")
}

func TestTrainView(t *testing.T) {
	assert := assert.New(t)

	emptyList := &models.TrainEventList{nil}
	service := new(mockTrainService)
	service.On("AllTrains", mock.Anything).Return(emptyList, nil)

	dtos := &dto.JSONTrainEventList{1, []dto.JSONTrainEvent{{42, 10, "dest", "", ""}}}
	converter := new(mockTrainConverter)
	converter.On("ConvertTrainList", emptyList).Return(dtos)

	view := NewTrainView(service, converter)
	out, err := view.AllTrains([]string{"anytype"})

	assert.Nil(err)
	assert.Equal(1, out.Count)
	assert.Equal(42, out.Trains[0].TrainNumber)
}

func TestTrainViewConvertionNotCalledOnServiceError(t *testing.T) {
	assert := assert.New(t)

	service := new(mockTrainService)
	service.On("AllTrains", mock.Anything).Return(nil, errors.New("fake error"))

	converter := new(mockTrainConverter)

	view := NewTrainView(service, converter)
	out, err := view.AllTrains([]string{"anytype"})

	assert.NotNil(err)
	assert.Nil(out)

	converter.AssertNotCalled(t, "ConvertTrains", mock.Anything)
}

func TestTrainViewServiceErrorsArePropagated(t *testing.T) {
	assert := assert.New(t)

	service := new(mockTrainService)
	service.On("AllTrains", mock.Anything).Return(nil, errors.New("fake error"))

	converter := new(mockTrainConverter)

	view := NewTrainView(service, converter)
	out, err := view.AllTrains([]string{})

	assert.NotNil(err)
	assert.Nil(out)
	assert.Equal("fake error", err.Error())
}

type mockTrainService struct {
	mock.Mock
}

func (m *mockTrainService) AllTrains(types []string) (result *models.TrainEventList, err error) {
	args := m.Called(types)
	e := args.Get(1)
	if e != nil {
		return nil, e.(error)
	} else {
		return args.Get(0).(*models.TrainEventList), nil
	}
}

func (m *mockTrainService) TrainsByKeyValue(key string, value string) (result *models.TrainEventList, err error) {
	args := m.Called(key, value)
	e := args.Get(1)
	if e != nil {
		return nil, e.(error)
	} else {
		return args.Get(0).(*models.TrainEventList), nil
	}
}

func (m *mockTrainService) TrainsFromStation(stationID services.StationID) (result *models.TrainEventList, err error) {
	args := m.Called(stationID)
	e := args.Get(1)
	if e != nil {
		return nil, e.(error)
	} else {
		return args.Get(0).(*models.TrainEventList), nil
	}
}

func (m *mockTrainService) Stops(trainID services.TrainID) (result *models.TrainEventList, err error) {
	args := m.Called(trainID)
	e := args.Get(1)
	if e != nil {
		return nil, e.(error)
	} else {
		return args.Get(0).(*models.TrainEventList), nil
	}
}

func (m *mockTrainService) DeparturesBetween(from services.StationID, to services.StationID) (result *models.TrainEventList, err error) {
	args := m.Called(from, to)
	e := args.Get(1)
	if e != nil {
		return nil, e.(error)
	} else {
		return args.Get(0).(*models.TrainEventList), nil
	}
}

type mockTrainConverter struct {
	mock.Mock
}

func (m mockTrainConverter) ConvertTrainList(list *models.TrainEventList) *dto.JSONTrainEventList {
	args := m.Called(list)
	return args.Get(0).(*dto.JSONTrainEventList)
}

func (m mockTrainConverter) ConvertTrain(t *models.TrainEvent) (*dto.JSONTrainEvent, error) {
	args := m.Called(t)
	e := args.Get(1)
	if e != nil {
		return nil, e.(error)
	} else {
		return args.Get(0).(*dto.JSONTrainEvent), nil
	}
}
