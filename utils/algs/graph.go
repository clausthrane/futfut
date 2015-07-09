// Package graph implements a collection of useful graph algorithms used by Futfut

package graph

import (
	"errors"
	"fmt"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/services"
	"github.com/clausthrane/futfut/utils"
	"log"
	"os"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

// Dijkstra implements https://en.wikipedia.org/wiki/Dijkstra's_algorithm
func Dijkstra(in *models.TrainEventList, src *models.TrainEvent, dst services.StationID) (*models.TrainEventList, error) {
	logger.Println("Building Graph")
	verticesList := ToVertices(in.Events)
	outEdges := ToEdgeMap(verticesList)
	logger.Printf("Exec Dijkstra's with |V| = %d and |E| = %d", len(verticesList), len(outEdges))

	targets := []*Vertex{}

	if _, err := seedVertex(src, outEdges); err != nil {
		logger.Println("Failed to seed an apropriate vertex for %v", src)
		return nil, err
	}

	q := FiFoQueue(verticesList)
	logger.Printf("|Q| = %d", len(q))
	for len(q) > 0 {
		u := q.Min()
		edges := outEdges[u.HashKey()]
		for _, e := range edges {
			e.Validate()
			v := e.to

			if q.Contains(v) {
				alt := utils.AddWithoutOverflow(u.TimeFromSource(), e.Price())
				if v.TimeFromSource() > alt {
					v.SetTimeFromSource(alt)
					v.SetPrev(u)
					v.inEdge = e
				}
			}
		}

		if dst == services.StationID(u.event.StationUic) && u.timeFromSource < MAX_DISTANCE {
			targets = append(targets, u)
		}
	}

	logger.Printf("Done! - %d targets found", len(targets))
	if result, err := route(targets); err == nil {
		return result, nil
	} else {
		msg := fmt.Sprintf("Could not produce route from %s at %s to %s", src.StationUic, src.ScheduledDeparture, dst)
		return nil, errors.New(msg)
	}
}

func seedVertex(src *models.TrainEvent, outEdges map[VertexKey][]*Edge) (*Vertex, error) {
	key := NewVertex(src.ScheduledDepartureDate(), src).HashKey()
	edges := outEdges[key]
	if len(edges) > 0 {
		srcV := edges[0].from
		srcV.SetTimeFromSource(0)
		logger.Printf("Looking for %p, %v", src, *src)
		return srcV, nil
	}
	return nil, errors.New("src has no outgoing edges!")
}

func ToEdgeMap(vertex []*Vertex) map[VertexKey][]*Edge {
	travelEdges := TravelEdges(vertex)
	logger.Printf("Generated %d travel edges", len(travelEdges))
	waitEdges := WaitEdges(vertex)
	logger.Printf("Generated %d wait edges", len(waitEdges))
	allEdges := append(travelEdges, waitEdges...)
	logger.Println("Indexing edges")

	result := make(map[VertexKey][]*Edge, 403)
	for _, edge := range allEdges {
		fromVertex := edge.from
		outEdges := result[fromVertex.HashKey()]
		result[fromVertex.HashKey()] = append(outEdges, edge)
	}
	return result
}

func route(options []*Vertex) (*models.TrainEventList, error) {
	if len(options) > 0 {
		best := earliestTarget(options)
		return buildResult(best), nil
	}
	return nil, errors.New("Cant build route with no options for a target")

}

func earliestTarget(targets []*Vertex) *Vertex {
	quickest := targets[0]
	for _, t := range targets {
		if t.timeFromSource < quickest.timeFromSource {
			quickest = t
		}
	}
	return quickest
}

func buildResult(v *Vertex) *models.TrainEventList {
	path := make([]models.TrainEvent, 0)
	for v != nil {
		path = append(path, *v.event)
		v = v.GetPrev()
	}
	return &models.TrainEventList{path}
}
