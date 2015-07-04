package dsb

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testTrainData = `
{
    "d": [
        {
            "__metadata": {
                "uri": "http://traindata.dsb.dk/StationDeparture/opendataprotocol.svc/Queue('230fffca-dfa4-4fae-bece-229a40bdf234')",
                "type": "ITogLogic.Model.Queue"
            },
            "ID": "230fffca-dfa4-4fae-bece-229a40bdf234",
            "StationUic": "7400028",
            "TrainType": "Ib",
            "DestinationName": "København h",
            "DestinationID": 8600626,
            "Track": "",
            "Generated": "/Date(1435928957000)/",
            "TrainNumber": "1468",
            "DestinationCountryCode": "86",
            "ScheduledArrival": null,
            "ScheduledDeparture": "/Date(1435939740000)/",
            "ArrivalDelay": "0",
            "DepartureDelay": "0",
            "Cancelled": false,
            "Line": null,
            "Direction": null,
            "MinutesToDeparture": null
        },
        {
            "__metadata": {
                "uri": "http://traindata.dsb.dk/StationDeparture/opendataprotocol.svc/Queue('40dc3d87-9496-4152-802a-a3ecb117e160')",
                "type": "ITogLogic.Model.Queue"
            },
            "ID": "40dc3d87-9496-4152-802a-a3ecb117e160",
            "StationUic": "7400028",
            "TrainType": "Ib",
            "DestinationName": "København h",
            "DestinationID": 8600626,
            "Track": "",
            "Generated": "/Date(1435928957000)/",
            "TrainNumber": "1476",
            "DestinationCountryCode": "86",
            "ScheduledArrival": null,
            "ScheduledDeparture": "/Date(1435954140000)/",
            "ArrivalDelay": "0",
            "DepartureDelay": "0",
            "Cancelled": false,
            "Line": null,
            "Direction": null,
            "MinutesToDeparture": null
        },{
            "__metadata": {
                "uri": "http://traindata.dsb.dk/StationDeparture/opendataprotocol.svc/Queue('fd0a665e-ba47-42d1-ac85-1927337a99e7')",
                "type": "ITogLogic.Model.Queue"
            },
            "ID": "fd0a665e-ba47-42d1-ac85-1927337a99e7",
            "StationUic": "8600761",
            "TrainType": "S-tog",
            "DestinationName": "Hundige",
            "DestinationID": 8600769,
            "Track": "2",
            "Generated": "/Date(1435936635503)/",
            "TrainNumber": "10246",
            "DestinationCountryCode": "86",
            "ScheduledArrival": null,
            "ScheduledDeparture": null,
            "ArrivalDelay": null,
            "DepartureDelay": null,
            "Cancelled": false,
            "Line": "A",
            "Direction": "Syd",
            "MinutesToDeparture": "23"
        }
    ]
}
`

func TestUnmarshalTrains(t *testing.T) {
	assert := assert.New(t)

	var container map[string][]json.RawMessage
	err := json.Unmarshal([]byte(testTrainData), &container)
	assert.Nil(err, "no errors expected")

	rawList := container["d"]
	assert.Equal(3, len(rawList), "expecting 3 trains")

	trainList := convertTrainJSONList(rawList)
	assert.Equal(3, len(trainList.Trains), "expecting 3 trains")

	train := trainList.Trains[2]
	assert.Equal("fd0a665e-ba47-42d1-ac85-1927337a99e7", train.ID, "should be the same")
	assert.Equal("8600761", train.StationUic, "should be the same")
	assert.Equal("S-tog", train.TrainType, "should be the same")
	assert.Equal("Hundige", train.DestinationName, "should be the same")
	assert.Equal(8600769, train.DestinationID, "should be the same")
	assert.Equal("2", train.Track, "should be the same")
	assert.Equal("/Date(1435936635503)/", train.Generated, "should be the same")
	assert.Equal("10246", train.TrainNumber, "should be the same")
	assert.Equal("86", train.DestinationCountryCode, "should be the same")
	assert.Equal("", train.ScheduledArrival, "should be the same")
	assert.Equal("", train.ScheduledDeparture, "should be the same")
	assert.Equal("", train.ArrivalDelay, "should be the same")
	assert.Equal("", train.DepartureDelay, "should be the same")
	assert.Equal(false, train.Cancelled, "should be the same")
	assert.Equal("A", train.Line, "should be the same")
	assert.Equal("Syd", train.Direction, "should be the same")
	assert.Equal("23", train.MinutesToDeparture, "should be the same")
}
