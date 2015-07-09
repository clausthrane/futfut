package graph

import (
	"fmt"
	"github.com/clausthrane/futfut/datasources/dsb"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/services"
	"github.com/clausthrane/futfut/tests/mockserver"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var e1 models.TrainEvent
var e2 models.TrainEvent
var e3 models.TrainEvent

func TestMain(m *testing.M) {
	logger.Println("Setting default test values in TestSuite")
	e1 = models.TrainEvent{"event_id1", "station_id1", "station_name", "type", "end_dest", 123, "track1",
		"", "Train1", "86", "/Date(1436068203501)/", "/Date(1436068203503)/",
		"", "", false, "", "", "",
	}

	e2 = models.TrainEvent{"event_id2", "station_id2", "station_name", "type", "end_dest", 123, "track1",
		"", "Train2", "86", "/Date(1436068203505)/", "/Date(1436068203507)/",
		"", "", false, "", "", "",
	}

	e3 = models.TrainEvent{"event_id3", "station_id3", "station_name", "type", "end_dest", 123, "track1",
		"", "Train2", "86", "/Date(1436068203508)/", "/Date(1436068203510)/",
		"", "", false, "", "", "",
	}
	os.Exit(m.Run())
}

func TestDijkstra(t *testing.T) {
	assert := assert.New(t)

	data := loadData(t)

	//fromid := "8600020" // Aalborg
	toid := services.StationID("8600029") // Arden

	var from *models.TrainEvent
	for idx, _ := range data.Events {
		e := &data.Events[idx]
		if e.ID == "b24ba020-7545-43b0-8135-48fffe742345" { // Aalborg departure
			if from == nil {
				logger.Printf("Trying src %s %v", data.Events[idx].ID, data.Events[idx])
				from = e
			}
		}
	}

	out, err := Dijkstra(data, from, toid)
	assert.Nil(err)
	assert.NotNil(out)
}

func TestDijkstraMultipleDepartures(t *testing.T) {
	assert := assert.New(t)

	time1 := time.Now()
	e1 := NewTestEvent("A", time1, "1")

	time2 := next(next(time1))
	e2 := NewTestEvent("A", time2, "1")

	time3 := next(next(time2))
	e3 := NewTestEvent("A", time3, "1")

	time4 := next(next(time3))
	e4 := NewTestEvent("A", time4, "2")

	time5 := next(next(time4))
	e5 := NewTestEvent("B", time5, "2")

	time6 := next(next(time5))
	e6 := NewTestEvent("B", time6, "2")

	time7 := next(next(time6))
	e7 := NewTestEvent("C", time7, "2")

	time8 := next(next(time7))
	e8 := NewTestEvent("D", time8, "3")

	time9 := next(next(time8))
	e9 := NewTestEvent("E", time9, "3")

	in := &models.TrainEventList{[]models.TrainEvent{*e1, *e2, *e3, *e4, *e5, *e6, *e7, *e8, *e9}}

	out1, err1 := Dijkstra(in, e1, services.StationID("C-id"))
	assert.Nil(err1)
	assert.NotNil(out1)
	assert.Equal(5, len(out1.Events))

	out2, err2 := Dijkstra(in, e8, services.StationID("E-id"))
	assert.Nil(err2)
	assert.NotNil(out2)
	assert.Equal(2, len(out2.Events))

	out3, err3 := Dijkstra(in, e1, services.StationID("E-id"))
	assert.NotNil(err3)
	assert.Nil(out3)

	/*
		for _, e := range out.Events {
			logger.Printf("$$$ %s", e.String())
		}
	*/
}

func next(t time.Time) time.Time {
	return t.Add(5 * time.Minute)
}

func TestToEdgeMap(t *testing.T) {
	assert := assert.New(t)

	e4 := e1
	e4.TrainNumber = "Train4"

	in := []*Vertex{
		NewVertex(e1.ScheduledArrivalDate(), &e1),
		NewVertex(e1.ScheduledArrivalDate().Add(1*time.Minute), &e4),
		NewVertex(e1.ScheduledArrivalDate().Add(2*time.Minute), &e4),
		NewVertex(e1.ScheduledArrivalDate().Add(3*time.Minute), &e4),
		NewVertex(e1.ScheduledArrivalDate().Add(6*time.Minute), &e4),
	}

	out := ToEdgeMap(in)

	assert.Equal(4, len(out), "Comming in on v1 only v1,v2,v3,v4 has extension")

	//	for k, v := range out {
	//		logger.Printf("\n\nk(%s) = (%d) : %s", k, len(v), v)
	//	}

	assert.NotNil(out)

}

// Not a test
func loadData(t *testing.T) *models.TrainEventList {

	// quickest way
	mockserver.HttpServerDSBTestApi(t, 44444)
	remoteAPI := dsb.NewDSBFacadeWithEndpoint("http://localhost:44444")
	successC, errC := remoteAPI.GetTrains("", "")

	select {
	case response := <-successC:
		return response
	case <-errC:
		return &models.TrainEventList{nil}
	}
}

func NewTestEvent(station string, arrival time.Time, train string) *models.TrainEvent {
	departure := next(arrival)
	dept := fmt.Sprintf("/Date(%d)/", departure.Unix()*1000)
	arri := fmt.Sprintf("/Date(%d)/", arrival.Unix()*1000)
	stationid := fmt.Sprintf("%s-id", station)

	return &models.TrainEvent{"", stationid, station, "", "", 0, "", "",
		train, "86", arri, dept, "", "", false, "", "", "",
	}
}
