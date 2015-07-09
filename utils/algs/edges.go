package graph

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

type Edge struct {
	// The source vertex of the edge
	from *Vertex

	// The target vertex of the edge
	to *Vertex

	// enRoute flags wether or not this edge represents a "physical" move
	enRoute bool
}

// NewTravelEdge returns a end modelling a "physical" move
func NewTravelEdge(from *Vertex, to *Vertex) (*Edge, error) {
	return buildEdge(from, to, true)
}

// NewTravelEdge returns a end modelling a "temoral" move
func NewWaitEdge(from *Vertex, to *Vertex) (*Edge, error) {
	return buildEdge(from, to, false)
}

// buildEdge constructs an Edge object and validates the in
func buildEdge(from *Vertex, to *Vertex, enRoute bool) (*Edge, error) {
	if from.event.TrainNumber != to.event.TrainNumber && from.event.StationUic != to.event.StationUic {
		logger.Printf("Not supposed to connect %s, %s", from.String(), to.String())
		return nil, errors.New("Cannot be connected!")
	}
	return &Edge{from, to, enRoute}, nil
}

// String pretty prints the edge
func (e Edge) String() string {

	fh, fm, _ := e.from.when.Clock()
	th, tm, _ := e.to.when.Clock()
	now := fmt.Sprintf("%d:%d", fh, fm)
	then := fmt.Sprintf("%d:%d", th, tm)

	action := "Continue on"
	if !e.enRoute {
		action = "Wait for"
	}

	return fmt.Sprintf("==> (%s @ %s) %s train(%s) to station(%s) (until %s)",
		now, e.from.event.StationUic, action, e.from.event.TrainNumber, e.to.event.StationUic, then)
}

func (e *Edge) Validate() {
	buildEdge(e.from, e.to, e.enRoute)
}

// Price calculates the distance in minutes from 'from' to 'to'
func (e Edge) Price() int64 {
	dur := e.to.when.Sub(e.from.when)
	return int64(dur) / int64(time.Minute)
}

// Partitioner represents a partitioning of the input so that a pair
// of vertices are in a single partition if they shoud be connected
type Partitioner func([]*Vertex) map[string][]*Vertex

// Edger represents a strategy for connecting the given vertex list
type Edger func([]*Vertex) []*Edge

// Edge when waiting to get on a train or switching
func WaitEdges(vertices []*Vertex) []*Edge {
	return PartitionAndConvert(vertices, PartitionByStationID, EdgesAtStation)
}

// Edge when en route
func TravelEdges(vertices []*Vertex) []*Edge {
	return PartitionAndConvert(vertices, PartitionByTrainNumber, EdgesOnRoute)
}

// PartitionAndConvert is a generic helper to drive the partitioner and edger
//
// Applying the Partitioner to the all parameter results in a collection of partitions
// each of which we apply the Edger to. The result of all edgers are appended in the output
func PartitionAndConvert(all []*Vertex, toPartitions Partitioner, edger Edger) []*Edge {
	result := make([]*Edge, 0)
	for _, partition := range toPartitions(all) {
		result = append(result, edger(partition)...)
	}
	return result
}

// EdgesAtStation generates all edges obtainable by waiting
// Output is size (cf. Triangular number series) N * (N+1)/2
func EdgesAtStation(vertices []*Vertex) []*Edge {
	sortedVerticies := TimeSortableEvents(vertices)
	sort.Sort(sortedVerticies)
	result := make([]*Edge, 0)
	for i := 0; i < len(sortedVerticies); i++ {
		for j := i + 1; j < len(sortedVerticies); j++ {
			edge, err := NewWaitEdge(sortedVerticies[i], sortedVerticies[j])
			if err == nil {
				result = append(result, edge)
			}
		}
	}
	return result
}

// EdgesOnRoute generate all edges obtained by following a route
// Output size N-1
func EdgesOnRoute(sequence []*Vertex) []*Edge {
	sortedSequence := TimeSortableEvents(sequence)
	sort.Sort(sortedSequence)
	edges := make([]*Edge, 0)
	if len(sortedSequence) > 0 {
		from := sortedSequence[0]
		for _, to := range sortedSequence[1:] {
			newEdge, err := NewTravelEdge(from, to)
			if err == nil {
				edges = append(edges, newEdge)
				from = to
			} else {
				logger.Fatalf("Unable to add edge %s => %s", from.String(), to.String())
			}
		}
	}
	return edges
}

func PartitionByStationID(events []*Vertex) map[string][]*Vertex {
	return IndexEvents(events, func(e *Vertex) string {
		return e.event.StationUic
	})
}

func PartitionByTrainNumber(events []*Vertex) map[string][]*Vertex {
	return IndexEvents(events, func(e *Vertex) string {
		return e.event.TrainNumber
	})
}

func IndexEvents(vertices []*Vertex, keyGenerator func(*Vertex) string) map[string][]*Vertex {
	m := make(map[string][]*Vertex)
	for _, vertex := range vertices {
		key := keyGenerator(vertex)
		list := m[key]
		m[key] = append(list, vertex)
	}
	return m
}
