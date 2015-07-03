package api

import (
	"errors"
	"github.com/clausthrane/futfut/views/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"testing"
)

func TestHandleStationsRequestWillMarshallViewsResult(t *testing.T) {
	assert := assert.New(t)

	view := new(mockView)
	view.On("AllStations").Return(&dto.JSONStationList{1, []dto.JSONStation{{"id", "name", "country"}}}, nil)

	responseWriter := new(mockResponseWriter)
	responseWriter.On("Write", mock.Anything).Return().Run(func(args mock.Arguments) {
		a := args.Get(0).([]byte)
		assert.NotNil(a, "expected come thing!")
		assert.Equal([]byte(`{"Count":1,"Stations":[{"Id":"id","Name":"name","CountryCode":"country"}]}`), a[:len(a)-1], "is marshalled")
	})

	handler := NewRequestHandler(view)
	handler.HandleStationsRequest(responseWriter, &http.Request{})
}

func TestHttpErrorCodes(t *testing.T) {
	assert := assert.New(t)

	view := new(mockView)
	view.On("AllStations").Return(nil, errors.New("some message"))

	responseWriter := new(mockResponseWriter)
	responseWriter.On("Header").Return(make(http.Header))

	responseWriter.On("Write", mock.Anything).Return().Run(func(args mock.Arguments) {
		a := args.Get(0).([]byte)
		assert.Equal("some message\n", string(a))
	})
	responseWriter.On("WriteHeader", 500).Return()

	handler := NewRequestHandler(view)
	handler.HandleStationsRequest(responseWriter, &http.Request{})
}

type mockView struct {
	mock.Mock
}

func (m mockView) AllStations() (list *dto.JSONStationList, err error) {
	args := m.Called()
	if e := args.Get(1); e == nil {
		return args.Get(0).(*dto.JSONStationList), nil
	} else {
		return nil, e.(error)
	}
}

func (m mockView) GetStations(countryCode string, countryName string, page int, pageSize int) (*dto.JSONStationList, error) {
	args := m.Called(countryCode, countryName, page, pageSize)
	return args.Get(0).(*dto.JSONStationList), args.Get(1).(error)
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
