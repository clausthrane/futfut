package views

import (
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

func TestServiceErrorsArePropagated(t *testing.T) {
	t.Fail()
}

func TestConverterIsCalled(t *testing.T) {
	assert := assert.New(t)

	var emptyList = &models.StationList{nil}
	mockService := new(mockSrv)
	mockService.On("AllStations").Return(emptyList, nil)
	mockConverter := new(mockConv)
	mockConverter.On("ConvertStationList", emptyList).Return(&dto.JSONStationList{1, []dto.JSONStation{dto.JSONStation{"id", "name", "country"}}})

	view := NewStationsView(mockService, mockConverter)

	out, err := view.AllStations()
	assert.Nil(err, "errors not expected")
	assert.Equal(1, out.Count)
	assert.Equal("id", out.Stations[0].Id, "first element expected to match our mock")
}

func TestConverterIsNotCalledOnServiceError(t *testing.T) {
	t.Fail()
}

type mockSrv struct {
	mock.Mock
}

func (s *mockSrv) AllStations() (res *models.StationList, err error) {
	args := s.Called()
	return args.Get(0).(*models.StationList), args.Error(1)
}

func (s *mockSrv) GetStations(countryCode string, countryName string, page int, pageSize int) (*models.StationList, error) {
	args := s.Called(countryCode, countryName, page, pageSize)
	return args.Get(0).(*models.StationList), args.Error(1)
}

type mockConv struct {
	mock.Mock
}

func (c mockConv) ConvertStationList(list *models.StationList) *dto.JSONStationList {
	args := c.Called(list)
	return args.Get(0).(*dto.JSONStationList)
}

func (c mockConv) ConvertStation(station *models.Station) dto.JSONStation {
	args := c.Called(station)
	return args.Get(0).(dto.JSONStation)
}
