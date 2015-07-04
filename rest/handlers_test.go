package api

import (
	"errors"
	"github.com/clausthrane/futfut/views/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestHandleStationsRequestWillMarshallViewsOutput(t *testing.T) {
	assert := assert.New(t)

	view := new(mockStationView)
	view.On("AllStations").Return(&dto.JSONStationList{1, []dto.JSONStation{{"id", "name", "country"}}}, nil)

	responseWriter := new(mockResponseWriter)
	responseWriter.On("Write", mock.Anything).Return().Run(func(args mock.Arguments) {
		a := args.Get(0).([]byte)
		assert.NotNil(a)
		assert.Equal([]byte(`{"Count":1,"Stations":[{"Id":"id","Name":"name","CountryCode":"country"}]}`), a[:len(a)-1])
	})

	handler := NewRequestHandler(view, nil)
	handler.HandleStationsRequest(responseWriter, &http.Request{})
}

func TestHandleTrainsRequestWillMarshallViewOutput(t *testing.T) {
	assert := assert.New(t)

	view := new(mockTrainView)

	viewResponse := &dto.JSONTrainList{1, []dto.JSONTrain{{1, 2, "hello"}}}
	view.On("AllTrains").Return(viewResponse, nil)

	responseWriterCalled := false
	responseWriter := new(mockResponseWriter)
	responseWriter.On("Write", mock.Anything).Return().Run(func(args mock.Arguments) {
		a := args.Get(0).([]byte)
		assert.NotNil(a)
		expected := `{"Count":1,"Trains":[{"TrainNumber":1,"CurrentStationId":2,"DestinationName":"hello"}]}`
		assert.Equal([]byte(expected), a[:len(a)-1])
		responseWriterCalled = true
	})

	handler := NewRequestHandler(nil, view)
	handler.HandleTrainsRequest(responseWriter, &http.Request{})
	assert.True(responseWriterCalled, "responsWriter has not been invoked")
}

func TestHttpErrorCodes(t *testing.T) {
	assert := assert.New(t)

	view := new(mockStationView)
	view.On("AllStations").Return(nil, errors.New("some message"))

	responseWriter := new(mockResponseWriter)
	responseWriter.On("Header").Return(make(http.Header))

	responseWriter.On("Write", mock.Anything).Return().Run(func(args mock.Arguments) {
		a := args.Get(0).([]byte)
		assert.Equal("some message\n", string(a))
	})
	responseWriter.On("WriteHeader", 500).Return()

	handler := NewRequestHandler(view, nil)
	handler.HandleStationsRequest(responseWriter, &http.Request{})
}

type mockTrainView struct {
	mock.Mock
}

func (m mockTrainView) AllTrains() (list *dto.JSONTrainList, err error) {
	args := m.Called()
	if e := args.Get(1); e == nil {
		return args.Get(0).(*dto.JSONTrainList), nil
	} else {
		return nil, e.(error)
	}
}

type mockStationView struct {
	mock.Mock
}

func (m mockStationView) AllStations() (list *dto.JSONStationList, err error) {
	args := m.Called()
	if e := args.Get(1); e == nil {
		return args.Get(0).(*dto.JSONStationList), nil
	} else {
		return nil, e.(error)
	}
}

func (m mockStationView) GetStations(countryCode string, countryName string, page int, pageSize int) (*dto.JSONStationList, error) {
	args := m.Called(countryCode, countryName, page, pageSize)
	if e := args.Get(1); e == nil {
		return args.Get(0).(*dto.JSONStationList), nil
	} else {
		return nil, e.(error)
	}
}

type mockResponseWriter struct {
	mock.Mock
}

func (m mockResponseWriter) Write(bytes []byte) (int, error) {
	m.Called(bytes)
	return len(bytes), nil
}

func (m mockResponseWriter) WriteHeader(i int) {
	m.Called(i)
}

func (m mockResponseWriter) Header() http.Header {
	args := m.Called()
	return args.Get(0).(http.Header)
}
