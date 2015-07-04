package dsb

import (
	"encoding/json"
	"github.com/clausthrane/futfut/tests/mockserver"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var stationTestData = `{
    "d": [
        {
            "__metadata": {
                "uri": "http://traindata.dsb.dk/StationDeparture/opendataprotocol.svc/Station('7400002')",
                "type": "ITogLogic.Model.Station"
            },
            "Abbreviation": "%G",
            "Name": "Göteborg",
            "UIC": "7400002",
            "CountryCode": "74",
            "CountryName": "S"
        },
        {
            "__metadata": {
                "uri": "http://traindata.dsb.dk/StationDeparture/opendataprotocol.svc/Station('7400003')",
                "type": "ITogLogic.Model.Station"
            },
            "Abbreviation": "%M",
            "Name": "Malmö Central",
            "UIC": "7400003",
            "CountryCode": "74",
            "CountryName": "S"
        },
        {
            "__metadata": {
                "uri": "http://traindata.dsb.dk/StationDeparture/opendataprotocol.svc/Station('7400006')",
                "type": "ITogLogic.Model.Station"
            },
            "Abbreviation": "%HM",
            "Name": "Hässleholm",
            "UIC": "7400006",
            "CountryCode": "74",
            "CountryName": "S"
        },{
            "__metadata": {
                "uri": "http://traindata.dsb.dk/StationDeparture/opendataprotocol.svc/Station('8600909')",
                "type": "ITogLogic.Model.Station"
            },
            "Abbreviation": "HØV",
            "Name": "Høvelte",
            "UIC": "8600909",
            "CountryCode": "86",
            "CountryName": "DK"
        }
    ]
}`

var corruptStationTestData = `{
    "d": [
        {
            "__metadata": {
                "uri": "http://traindata.dsb.dk/StationDeparture/opendataprotocol.svc/Station('7400002')",
                "type": "ITogLogic.Model.Station"
            },
            "Abbreviation": "%G",
            "Name": "Göteborg",
            "UIC": "7400002",
            "CountryCode": "74",
            "CountryName": "S"
        },{
            "__metadata": {
                "uri": "http://traindata.dsb.dk/StationDeparture/opendataprotocol.svc/Station('8600909')",
                "type": "ITogLogic.Model.Station"
            },
            "Abbreviation": "HØV",
            "Name": "Høvelte",
            "UIC": "8600909",
            "CountryCode": 86,
            "CountryName": "DK"
        }
    ]
}`

func TestBuildRequest(t *testing.T) {
	assert := assert.New(t)

	apiFacade := NewDSBFacade()
	apiFacade.setEndpoint("http://example.com")

	out, noerr := apiFacade.buildRequest(httpGET, "/hello")
	headerValue := out.Header.Get("Accept")
	url := *out.URL

	assert.Nil(noerr, "no errors expected here")
	assert.Equal("Application/JSON", headerValue, "accept type should be JSON")
	assert.Equal("http://example.com/hello", url.String(), "expecting url to have been assembled")
}

func TestJsonUnmarshaller(t *testing.T) {
	assert := assert.New(t)

	var container map[string][]json.RawMessage
	err := json.Unmarshal([]byte(`en lille nisse rejste`), &container)

	assert.NotNil(err, "error expected")
}

func TestUnmarshalStations(t *testing.T) {
	assert := assert.New(t)

	var container map[string][]json.RawMessage
	err := json.Unmarshal([]byte(stationTestData), &container)
	assert.Nil(err, "no errors expected")

	rawList := container["d"]
	out := convertStationJSONList(rawList)

	assert.Equal(4, len(out.Stations), "all stations should have been unmarshalled")
	assert.Equal("Göteborg", out.Stations[0].Name)
	assert.Equal("7400003", out.Stations[1].UIC)
	assert.Equal("74", out.Stations[2].CountryCode)
	assert.Equal("DK", out.Stations[3].CountryName)
}

func TestUnmarshalStationsFails(t *testing.T) {
	assert := assert.New(t)

	var container map[string][]json.RawMessage
	err := json.Unmarshal([]byte(corruptStationTestData), &container)
	assert.Nil(err, "no errors expected")

	rawList := container["d"]
	out := convertStationJSONList(rawList)

	assert.Nil(err, "error not expected")
	assert.NotNil(out, "result expected")
	assert.Equal(1, len(out.Stations))
}

func TestGetStations(t *testing.T) {
	assert := assert.New(t)

	mockserver.HttpServerWithStatusCode(7771, 500)

	client := NewDSBFacade()
	client.setEndpoint("http://localhost:7771")

	var expectedError error
	succ, fail := client.GetStations()
	select {
	case <-succ:
		t.Error("Request supposed to fail")
	case expectedError = <-fail:
	case <-time.After(time.Second * 1):
		t.Error("Not supposed to timeout for local mock server")
	}
	assert.NotNil(expectedError, "fdsf")
}
