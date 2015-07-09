package graph

import (
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

	//fromid := "8600020"
	toid := "8600029"

	var from *models.TrainEvent
	var to *models.TrainEvent
	for idx, _ := range data.Events {
		e := &data.Events[idx]
		if e.ID == "b24ba020-7545-43b0-8135-48fffe742345" {
			Dijkstra(data, e, services.StationID(toid))

			if from == nil {
				logger.Printf("Trying src %s %v", data.Events[idx].ID, data.Events[idx])
				from = e
			}
		}
	}

	logger.Println(from)
	logger.Println(to)

	out, err := Dijkstra(data, from, services.StationID(toid))
	assert.Nil(err)
	assert.NotNil(out)
	for _, evnt := range out.Events {
		logger.Printf("##> %v", evnt)
	}

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
