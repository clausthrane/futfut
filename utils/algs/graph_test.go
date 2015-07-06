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

var e1 models.Train
var e2 models.Train
var e3 models.Train

func TestMain(m *testing.M) {
	logger.Println("Setting default test values in TestSuite")
	e1 = models.Train{"event_id1", "station_id1", "type", "end_dest", 123, "track1",
		"", "Train1", "86", "/Date(1436068203501)/", "/Date(1436068203503)/",
		"", "", false, "", "", "",
	}

	e2 = models.Train{"event_id2", "station_id2", "type", "end_dest", 123, "track1",
		"", "Train2", "86", "/Date(1436068203505)/", "/Date(1436068203507)/",
		"", "", false, "", "", "",
	}

	e3 = models.Train{"event_id3", "station_id3", "type", "end_dest", 123, "track1",
		"", "Train2", "86", "/Date(1436068203508)/", "/Date(1436068203510)/",
		"", "", false, "", "", "",
	}
	os.Exit(m.Run())
}

func TestDijkstra(t *testing.T) {
	data := loadData(t)

	// Vertex(at 8600192 @ 2015-07-05 16:45:00 +0800 CST)
	from := services.StationID("8600646")
	to := services.StationID("8600783")

	Dijkstra(data, from, "", to, "")
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

	out, _, _ := ToEdgeMap(in, services.StationID(in[0].event.StationUic),
		services.StationID(in[3].event.StationUic), "", "")

	assert.Equal(4, len(out), "Comming in on v1 only v1,v2,v3,v4 has extension")

	//	for k, v := range out {
	//		logger.Printf("\n\nk(%s) = (%d) : %s", k, len(v), v)
	//	}

	assert.NotNil(out)

}

// Not a test
func loadData(t *testing.T) *models.TrainList {

	// quickest way
	mockserver.HttpServerDSBTestApi(t, 44444)
	remoteAPI := dsb.NewDSBFacadeWithEndpoint("http://localhost:44444")
	successC, errC := remoteAPI.GetTrains("", "")

	select {
	case response := <-successC:
		return response
	case <-errC:
		return &models.TrainList{nil}
	}
}
