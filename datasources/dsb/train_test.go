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
	assert.Equal(3, len(trainList.Events), "expecting 3 trains")

	logger.Println(trainList.Events[0])
	logger.Println(trainList.Events[1])
	logger.Println(trainList.Events[2])

	train := trainList.Events[2]
	assert.Equal("40dc3d87-9496-4152-802a-a3ecb117e160", train.ID)
	assert.Equal("7400028", train.StationUic)
	assert.Equal("Ib", train.TrainType)
	assert.Equal("København h", train.DestinationName)
	assert.Equal(8600626, train.DestinationID)
	assert.Equal("", train.Track)
	assert.Equal("/Date(1435928957000)/", train.Generated)
	assert.Equal("1476", train.TrainNumber)
	assert.Equal("86", train.DestinationCountryCode)
	assert.Equal("", train.ScheduledArrival)
	assert.Equal("/Date(1435954140000)/", train.ScheduledDeparture)
	assert.Equal("0", train.ArrivalDelay)
	assert.Equal("0", train.DepartureDelay)
	assert.Equal(false, train.Cancelled)
	assert.Equal("", train.Line)
	assert.Equal("", train.Direction)
	assert.Equal("", train.MinutesToDeparture)
}
