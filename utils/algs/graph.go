// Package graph implements a collection of useful graph algorithms used by Futfut

package graph

import (
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/services"
	"log"
	"os"
)

var logger = log.New(os.Stdout, " ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

// Dijkstra implements https://en.wikipedia.org/wiki/Dijkstra's_algorithm
func Dijkstra(in *models.TrainList, src services.StationID, departureTime string, dst services.StationID, arrivalTime string) []*Vertex {
	logger.Println("Building Graph")

	verticesList := ToVertices(in.Trains)
	logger.Printf("Having %d vertices, from %d records", len(verticesList), len(in.Trains))

	outEdges, srcV, tgtV := ToEdgeMap(verticesList, src, dst, departureTime, arrivalTime)

	logger.Printf("Running Dijkstra's algorithm on G=(V,E) with |V| = %d and |E| = %d (where %d vertices are not terminals)",
		len(verticesList), len(outEdges))

	if srcV != nil || tgtV != nil {
		srcV.SetTimeFromSource(0)
		logger.Printf("Looking for %p, %v", tgtV, *tgtV)
		logger.Printf("From for %p, %v", srcV, *srcV)
	} else {
		return nil
	}

	q := makeQueue(outEdges)

	logger.Printf("Entering main loop with |Q| = %d", len(q))
	for len(q) > 0 {
		u := q.Min()
		edges := outEdges[u.HashKey()]

		///*

		logger.Printf("Dequeued %s (distance = %d) having %d edges (%d remaining)",
			u.String(), u.TimeFromSource(), len(edges), len(q))
		//*/

		for _, e := range edges {
			e.Validate()
			v := e.to

			if q.Contains(v) {
				alt := AddWithoutOverflow(u.TimeFromSource(), e.Price())
				if v.TimeFromSource() > alt {
					v.SetTimeFromSource(alt)
					v.SetPrev(u)
					logger.Printf("Found shorter route to %s", v.String())

				}
			}
		}
	}

	logger.Println("Done!")
	return buildResult(tgtV)
}

func buildResult(v *Vertex) []*Vertex {
	path := make([]*Vertex, 0)
	for v != nil {
		to := v
		v = v.GetPrev()
		if v != nil {
			if edge, err := buildEdge(v, to, false); err == nil {
				logger.Printf("## %s", edge.String())
			} else {
				logger.Fatal("HUH!?!")
			}
		}
	}
	return path
}

func AddWithoutOverflow(a int64, b int64) int64 {
	sum := a + b
	if sum < a {
		return a
	} else {
		return sum
	}
}

func makeQueue(edgeSet map[VertexKey][]*Edge) FiFoQueue {
	q := make([]*Vertex, 0)
	for _, v := range edgeSet {
		q = append(q, v[0].from)
	}
	return FiFoQueue(q)
}

func ToEdgeMap(vertex []*Vertex, from services.StationID, to services.StationID, departureTime string, arrivalTime string) (map[VertexKey][]*Edge, *Vertex, *Vertex) {

	var src, tgt *Vertex

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
	logger.Printf("Vertices with out edges: %d", len(result))
	return result, src, tgt
}
