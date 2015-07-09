package graph

type FiFoQueue []*Vertex

func (q *FiFoQueue) Push(v *Vertex) {
	*q = append(*q, v)
}

func (q *FiFoQueue) Pop() *Vertex {
	v := (*q)[0]
	*q = (*q)[1:]
	return v
}

func (q *FiFoQueue) Min() *Vertex {
	idx := 0
	res := (*q)[idx]
	for i := 1; i < len(*q); i++ {
		v := (*q)[i]
		if res.TimeFromSource() > v.TimeFromSource() {
			res = v
			idx = i
		}
	}

	// remove the item!!!
	x := (*q)[:idx]
	y := (*q)[idx+1:]
	z := append(x, y...)
	*q = z
	return res
}

func (q *FiFoQueue) Contains(v *Vertex) bool {
	for i := 0; i < len(*q); i++ {
		in := (*q)[i]
		if in == v {
			return true
		}
	}
	return false
}
