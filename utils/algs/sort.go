package graph

// Convinience interface for sorting by date
type TimeSortableEvents []*Vertex

func (s TimeSortableEvents) Len() int {
	return len(s)
}

func (s TimeSortableEvents) Less(i int, j int) bool {
	return s[i].HappensBefore(s[j])
}

func (s TimeSortableEvents) Swap(i int, j int) {
	s[i], s[j] = s[j], s[i]
}
