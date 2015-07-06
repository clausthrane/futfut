package graph

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewEdgeInvariant(t *testing.T) {
	assert := assert.New(t)

	from := NewVertex(e1.ScheduledArrivalDate(), &e1)
	to := NewVertex(e2.ScheduledArrivalDate(), &e2)
	edge, err := buildEdge(from, to, true)

	assert.Nil(edge)
	assert.NotNil(err)
}

func TestPartitionByStationId(t *testing.T) {
	assert := assert.New(t)
	in := []*Vertex{
		NewVertex(e1.ScheduledArrivalDate(), &e1),
		NewVertex(e2.ScheduledArrivalDate(), &e2),
		NewVertex(e3.ScheduledArrivalDate(), &e3),
	}

	out := PartitionByStationID(in)

	assert.Equal(3, len(out))

	par0 := out["station_id1"]
	assert.Equal(1, len(par0))
	par1 := out["station_id2"]
	assert.Equal(1, len(par1))
	par2 := out["station_id3"]
	assert.Equal(1, len(par2))
	par3 := out["station_id4"]
	assert.Equal(0, len(par3))
}

func TestPartitionByTrainBumber(t *testing.T) {
	assert := assert.New(t)
	in := []*Vertex{
		NewVertex(e1.ScheduledArrivalDate(), &e1),
		NewVertex(e2.ScheduledArrivalDate(), &e2),
		NewVertex(e3.ScheduledArrivalDate(), &e3),
	}

	out := PartitionByTrainNumber(in)

	assert.Equal(2, len(out))
	assert.Equal(2, len(out["Train2"]))
	assert.Equal("Train2", out["Train2"][0].event.TrainNumber)
	assert.True(len(out["Train2"][0].event.ScheduledArrival) > 0, "events still have their timestamps")
}

func TestEdgesOnRoute(t *testing.T) {
	assert := assert.New(t)

	in := []*Vertex{
		NewVertex(e2.ScheduledArrivalDate(), &e2),
		NewVertex(e3.ScheduledArrivalDate(), &e3),
	}

	edges := EdgesOnRoute(in)
	assert.Equal(1, len(edges))
}

func TestEdgesAtStation(t *testing.T) {
	assert := assert.New(t)

	e4 := e1
	e4.TrainNumber = "Train4"

	assert.NotEqual(e1.TrainNumber, e4.TrainNumber)

	in := []*Vertex{
		NewVertex(e1.ScheduledArrivalDate(), &e1),
		NewVertex(e1.ScheduledArrivalDate().Add(1*time.Minute), &e4),
		NewVertex(e1.ScheduledArrivalDate().Add(2*time.Minute), &e4),
		NewVertex(e1.ScheduledArrivalDate().Add(3*time.Minute), &e4),
		NewVertex(e1.ScheduledArrivalDate().Add(6*time.Minute), &e4),
	}

	out := WaitEdges(in)
	assert.Equal(10, len(out))

	//for _, edge := range out {
	//	logger.Println(edge.String())
	//}
}

func TestWaitEdges(t *testing.T) {
	t.Fail()
}

func TestTravelEdges(t *testing.T) {
	t.Fail()
}

func TestPartitionAndConvert(t *testing.T) {
	t.Fail()
}
