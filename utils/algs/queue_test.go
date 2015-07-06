package graph

import (
	"github.com/clausthrane/futfut/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFifoQueue(t *testing.T) {
	assert := assert.New(t)

	e1 := &models.Train{"event_id1", "station_id1", "type", "end_dest", 123, "track1",
		"", "Train1", "86", "/Date(1436068203501)/", "", "", "", false, "", "", "",
	}

	e2 := &models.Train{"event_id2", "station_id2", "type", "end_dest", 123, "track1",
		"", "Train2", "86", "/Date(1436068203502)/",
		"", "", "", false, "", "", "",
	}

	e3 := &models.Train{"event_id3", "station_id3", "type", "end_dest", 123, "track1",
		"", "Train2", "86", "/Date(1436068203503)/",
		"", "", "", false, "", "", "",
	}

	verticies := make([]*Vertex, 0)
	verticies = append(verticies, NewVertex(e1.ScheduledArrivalDate(), e1))
	verticies = append(verticies, NewVertex(e2.ScheduledArrivalDate(), e2))

	q := FiFoQueue(verticies)

	q.Push(NewVertex(e3.ScheduledArrivalDate(), e3))

	a := q.Pop().event.StationUic
	state1 := len(q)
	b := q.Pop().event.StationUic
	c := q.Pop().event.StationUic
	state2 := len(q)

	assert.Equal("station_id1", a)
	assert.Equal("station_id2", b)
	assert.Equal("station_id3", c)
	assert.Equal(2, state1)
	assert.Equal(0, state2)
}

func TestContains(t *testing.T) {
	assert := assert.New(t)

	e1 := &models.Train{"event_id1", "station_id1", "type", "end_dest", 123, "track1",
		"", "Train1", "86", "/Date(1436068203501)/", "", "", "", false, "", "", "",
	}

	e2 := &models.Train{"event_id2", "station_id2", "type", "end_dest", 123, "track1",
		"", "Train2", "86", "/Date(1436068203502)/",
		"", "", "", false, "", "", "",
	}

	inSider := NewVertex(e1.ScheduledArrivalDate(), e1)
	outSider := NewVertex(e2.ScheduledArrivalDate(), e2)

	verticies := make([]*Vertex, 0)
	verticies = append(verticies, inSider)

	q := FiFoQueue(verticies)

	assert.True(q.Contains(inSider))
	assert.False(q.Contains(outSider))
}
