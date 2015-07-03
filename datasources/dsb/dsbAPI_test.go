package dsb

import (
	"github.com/clausthrane/futfut/tests/mockserver"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var testdata = `{
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

var corrupttestdata = `{
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

func TestUnmarshalStations(t *testing.T) {
	assert := assert.New(t)

	out, err := unmarshalStations([]byte(testdata))

	assert.Nil(err, "no errors expected")
	assert.Equal(4, len(out.Stations), "all stations should have been unmarshalled")
	assert.Equal("Göteborg", out.Stations[0].Name)
	assert.Equal("7400003", out.Stations[1].UIC)
	assert.Equal("74", out.Stations[2].CountryCode)
	assert.Equal("DK", out.Stations[3].CountryName)
}

func TestUnmarshalStationsFails(t *testing.T) {
	assert := assert.New(t)

	out, err := unmarshalStations([]byte(`en lille nisse rejste`))

	assert.Nil(out, "only error expected")
	assert.NotNil(err, "error expected")

	out, err = unmarshalStations([]byte(corrupttestdata))

	assert.Nil(err, "error not expected")
	assert.NotNil(out, "result expected")
	assert.Equal(1, len(out.Stations))
}

func TestBuildRequestSetsJSONHeadder(t *testing.T) {
	assert := assert.New(t)

	out, err := NewDSBFacade().buildRequest()

	assert.Nil(err, "no errors expected")
	assert.NotNil(out, "output expected")

	headerValue := out.Header.Get("Accept")
	assert.Equal("Application/JSON", headerValue, "accept type should be JSON")
}

func TestGetStations(t *testing.T) {
	assert := assert.New(t)

	mockserver.HttpServerWithStatusCode(7771, 500)

	client := NewDSBFacade()
	client.SetEndpoint("http://localhost:7771")

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
