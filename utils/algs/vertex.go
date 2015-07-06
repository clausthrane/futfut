package graph

import (
	"fmt"
	"github.com/clausthrane/futfut/models"
	"math"
	"time"
)

const MAX_DISTANCE int64 = math.MaxInt64

type Vertex struct {
	event          *models.Train
	when           time.Time
	timeFromSource int64
	prev           *Vertex
}

type VertexKey string

func HashKey(stationid string, time string) VertexKey {
	return VertexKey(fmt.Sprintf("%s:%s", stationid, time))
}

func NewVertex(time time.Time, event *models.Train) *Vertex {
	return &Vertex{event, time, MAX_DISTANCE, nil}
}

func (v Vertex) HashKey() VertexKey {
	//return VertexKey(fmt.Sprintf("%s:%s", v.event.StationUic, v.when.String()))
	return HashKey(v.event.StationUic, v.when.String())
}

func (v *Vertex) String() string {
	return fmt.Sprintf("Vertex(at %s @ %s)", v.event.StationUic, v.when.String())
}

func (v *Vertex) HappensBefore(other *Vertex) bool {
	return v.when.Before(other.when)
}

func (v *Vertex) TimeFromSource() int64 {
	return v.timeFromSource
}

func (v *Vertex) SetTimeFromSource(minutes int64) {
	v.timeFromSource = minutes
}

func (v *Vertex) SetPrev(u *Vertex) {
	v.prev = u
}

func (v *Vertex) GetPrev() *Vertex {
	return v.prev
}

// Makes all vertices in the graph
//
//
func ToVertices(events []models.Train) []*Vertex {
	list := make([]*Vertex, 0, len(events))
	for i, e := range events {
		if isUsable(e) {
			currentEvent := &events[i]

			arrival := e.ScheduledArrivalDate()
			arrivalVertex := NewVertex(arrival, currentEvent)
			departure := e.ScheduledDepartureDate()
			departureVertex := NewVertex(departure, currentEvent)

			list = append(list, arrivalVertex)
			list = append(list, departureVertex)
		}
	}
	return list
}

// Skipping arrivals without time
// S-tog don't have timing data, so we ignore them
func isUsable(t models.Train) bool {
	return len(t.ScheduledArrival) > 0
}
