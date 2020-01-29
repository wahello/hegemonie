// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"errors"
	"sort"
	"strconv"
	"strings"
	"sync/atomic"
)

func (r SetOfVertices) Len() int           { return len(r) }
func (r SetOfVertices) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r SetOfVertices) Less(i, j int) bool { return r[i].Id < r[j].Id }

func (s SetOfVertices) Get(id uint64) *MapVertex {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].Id >= id
	})
	if i < len(s) && s[i].Id == id {
		return s[i]
	}
	return nil
}

func (s *SetOfVertices) Add(v *MapVertex) {
	*s = append(*s, v)
	if nb := len(*s); nb > 2 && !sort.IsSorted((*s)[nb-2:]) {
		sort.Sort(*s)
	}
}

func (r SetOfEdges) Len() int      { return len(r) }
func (r SetOfEdges) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r SetOfEdges) Less(i, j int) bool {
	s, d := r[i], r[j]
	return s.S < d.S || (s.S == d.S && s.D < d.D)
}

func (r SetOfEdges) First(at uint64) int {
	i := sort.Search(len(r), func(i int) bool { return r[i].S >= at })
	return i
}

func (s *SetOfEdges) Add(e *MapEdge) {
	*s = append(*s, e)
	if nb := len(*s); nb > 2 && !sort.IsSorted((*s)[nb-2:]) {
		sort.Sort(*s)
	}
}

func (s SetOfEdges) Get(src, dst uint64) *MapEdge {
	i := sort.Search(len(s), func(i int) bool {
		return s[i].S >= src || (s[i].S == src && s[i].D == dst)
	})
	if i < len(s) && s[i].S == src && s[i].D == dst {
		return s[i]
	}
	return nil

}

func (m *Map) Init() {
	m.Cells = make([]*MapVertex, 0)
	m.Roads = make([]*MapEdge, 0)
}

func (m *Map) getNextId() uint64 {
	return atomic.AddUint64(&m.NextId, 1)
}

func (m *Map) CellGet(id uint64) *MapVertex {
	return m.Cells.Get(id)
}

func (m *Map) CellHas(id uint64) bool {
	return m.CellGet(id) != nil
}

func (m *Map) CellCreate() *MapVertex {
	c := &MapVertex{Id: m.getNextId()}
	m.Cells.Add(c)
	return c
}

// Raw creation of an edge, with no check the Source and Destination exist
// The set of roads isn't sorted afterwards
func (m *Map) RoadCreateRaw(src, dst uint64) *MapEdge {
	if src == dst || src == 0 || dst == 0 {
		panic("Invalid Edge parameters")
	}

	e := &MapEdge{src, dst, false}
	m.Roads = append(m.Roads, e)
	return e
}

func (m *Map) RoadCreate(src, dst uint64, check bool) error {
	if src == dst || src == 0 || dst == 0 {
		return errors.New("EINVAL")
	}

	if check && !m.CellHas(src) {
		return errors.New("Source not found")
	}
	if check && !m.CellHas(dst) {
		return errors.New("Destination not found")
	}

	if r := m.Roads.Get(src, dst); r != nil {
		if r.Deleted {
			return errors.New("MapEdge exists")
		} else {
			r.Deleted = false
			return nil
		}
	} else {
		m.Roads.Add(&MapEdge{src, dst, false})
		return nil
	}
}

func (m *Map) RoadDelete(src, dst uint64, check bool) error {
	if src == dst || src == 0 || dst == 0 {
		return errors.New("EINVAL")
	}

	if check && !m.CellHas(src) {
		return errors.New("Source not found")
	}
	if check && !m.CellHas(dst) {
		return errors.New("Destination not found")
	}

	if r := m.Roads.Get(src, dst); r != nil {
		if !r.Deleted {
			r.Deleted = true
			return nil
		} else {
			return errors.New("MapEdge closed")
		}
	} else {
		return errors.New("MapEdge not found")
	}
}

func (m *Map) PathNextStep(src, dst uint64) (uint64, error) {
	if src == dst || src == 0 || dst == 0 {
		return 0, errors.New("EINVAL")
	}

	next, ok := m.steps[vector{src, dst}]
	if ok {
		return next, nil
	} else {
		return 0, errors.New("No route")
	}
}

func (m *Map) CellAdjacency(id uint64) []uint64 {
	adj := make([]uint64, 0)

	for i := m.Roads.First(id); i < len(m.Roads); i++ {
		r := m.Roads[i]
		if r.S != id {
			break
		}
		adj = append(adj, r.D)
	}

	return adj
}

func (m *Map) Check(w *World) error {
	return nil
}

func (m *Map) Dot() string {
	var sb strings.Builder
	sb.WriteString("digraph g {")
	for _, c := range m.Cells {
		sb.WriteString("n" + strconv.FormatUint(c.Id, 10))
		sb.WriteRune(';')
		sb.WriteRune('\n')
	}
	for _, r := range m.Roads {
		sb.WriteRune(' ')
		sb.WriteString("n" + strconv.FormatUint(r.S, 10))
		sb.WriteString(" -> ")
		sb.WriteString("n" + strconv.FormatUint(r.D, 10))
		sb.WriteRune(';')
		sb.WriteRune('\n')
	}
	sb.WriteString("}")
	return sb.String()
}

func (m *Map) Rehash() {
	next := make(map[vector]uint64)

	// Ensure the locations are sorted
	sort.Sort(&m.Cells)
	sort.Sort(&m.Roads)

	// Fill with the immediate neighbors
	for _, r := range m.Roads {
		next[vector{r.S, r.D}] = r.D
	}

	add := func(src, dst, step uint64) {
		_, found := next[vector{src, dst}]
		if !found {
			next[vector{src, dst}] = step
		}
	}

	// Call one DFS per node and shortcut when possible
	for _, cell := range m.Cells {
		already := make(map[uint64]bool)
		q := newQueue()

		// Bootstrap the DFS with adjacent nodes
		for _, next := range m.CellAdjacency(cell.Id) {
			q.push(next, next)
			already[next] = true
			// No need to add this in the known routes, we already did it
			// with an iteration on the roads (much faster)
		}

		for !q.empty() {
			current, first := q.pop()
			neighbors := m.CellAdjacency(current)
			// TODO(jfs): shuffle the neighbors
			for _, next := range neighbors {
				if !already[next] {
					// Avoid passing again in the neighbor
					already[next] = true
					// Tell to contine at that neighbor
					q.push(next, first)
					// We already learned the shortest path to that neighbor
					add(cell.Id, next, first)
				}
			}
		}
	}

	m.steps = next
}

type vector struct {
	src uint64
	dst uint64
}

type dfsTrack struct {
	current uint64
	first   uint64
}

type queue struct {
	tab   []dfsTrack
	start int
}

func newQueue() queue {
	var q queue
	q.tab = make([]dfsTrack, 0)
	return q
}

func (q *queue) push(node, first uint64) {
	q.tab = append(q.tab, dfsTrack{node, first})
}

func (q *queue) pop() (uint64, uint64) {
	v := q.tab[q.start]
	q.start++
	return v.current, v.first
}

func (q *queue) empty() bool {
	return q.start == len(q.tab)
}
