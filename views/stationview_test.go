package views

import (
	"errors"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/views/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestNewStationsView(t *testing.T) {
	assert := assert.New(t)

	out := NewStationsView(nil, nil)

	assert.Panics(nil, func() {
		out.AllStations()
	}, "Calling AllStations should panic due to nil")

}

func TestStationView(t *testing.T) {
	assert := assert.New(t)

	var emptyList = &models.StationList{nil}
	mockService := new(mockStationService)
	mockService.On("AllStations").Return(emptyList, nil)
	mockConverter := new(mockStationConverter)
	mockConverter.On("ConvertStationList", emptyList).Return(&dto.JSONStationList{1, []dto.JSONStation{dto.JSONStation{"id", "name", "country"}}})

	view := NewStationsView(mockService, mockConverter)

	out, err := view.AllStations()
	assert.Nil(err)
	assert.Equal(1, out.Count)
	assert.Equal("id", out.Stations[0].Id)
}

func TestStationViewServiceErrorsArePropagated(t *testing.T) {
	assert := assert.New(t)

	service := new(mockStationService)
	service.On("AllStations").Return(nil, errors.New("fake error"))

	converter := new(mockStationConverter)

	view := NewStationsView(service, converter)
	out, err := view.AllStations()

	assert.NotNil(err)
	assert.Nil(out)
	assert.Equal("fake error", err.Error())
}

func TestStationViewConverterIsNotCalledOnServiceError(t *testing.T) {
	assert := assert.New(t)

	service := new(mockStationService)
	service.On("AllStations").Return(nil, errors.New("fake error"))

	converter := new(mockStationConverter)

	view := NewStationsView(service, converter)
	out, err := view.AllStations()

	assert.NotNil(err)
	assert.Nil(out)

	converter.AssertNotCalled(t, "ConvertStations", mock.Anything)
}

type mockStationService struct {
	mock.Mock
}

func (s *mockStationService) AllStations() (res *models.StationList, err error) {
	args := s.Called()
	e := args.Get(1)
	if e != nil {
		return nil, e.(error)
	} else {
		return args.Get(0).(*models.StationList), nil
	}
}

func (s *mockStationService) GetStations(countryCode string, countryName string, page int, pageSize int) (*models.StationList, error) {
	args := s.Called(countryCode, countryName, page, pageSize)
	return args.Get(0).(*models.StationList), args.Error(1)
}

type mockStationConverter struct {
	mock.Mock
}

func (c mockStationConverter) ConvertStationList(list *models.StationList) *dto.JSONStationList {
	args := c.Called(list)
	return args.Get(0).(*dto.JSONStationList)
}

func (c mockStationConverter) ConvertStation(station *models.Station) dto.JSONStation {
	args := c.Called(station)
	return args.Get(0).(dto.JSONStation)
}
