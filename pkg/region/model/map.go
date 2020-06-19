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

func (ev SetOfEdges) Len() int      { return len(ev) }
func (ev SetOfEdges) Swap(i, j int) { ev[i], ev[j] = ev[j], ev[i] }
func (ev SetOfEdges) Less(i, j int) bool {
	return ev.edgeLess(i, *ev[j])
}

func (ev SetOfEdges) edgeLess(i int, d MapEdge) bool {
	s := ev[i]
	return s.S < d.S || (s.S == d.S && s.D < d.D)
}

func (ev SetOfEdges) First(at uint64) int {
	return sort.Search(len(ev), func(i int) bool { return ev[i].S >= at })
}

func (ev *SetOfEdges) Add(e *MapEdge) {
	*ev = append(*ev, e)
	if nb := len(*ev); nb > 2 && !sort.IsSorted((*ev)[nb-2:]) {
		sort.Sort(*ev)
	}
}

func (ev SetOfEdges) Get(src, dst uint64) *MapEdge {
	i := sort.Search(len(ev), func(i int) bool {
		return ev[i].S >= src || (ev[i].S == src && ev[i].D >= dst)
	})
	if i < len(ev) && ev[i].S == src && ev[i].D == dst {
		return ev[i]
	}
	return nil

}

func (ev SetOfEdges) Slice(markerSrc, markerDst uint64, max uint32) []MapEdge {
	tab := make([]MapEdge, 0)

	iMax := ev.Len()
	i := ev.First(markerSrc)
	if i < iMax && ev[i].S == markerSrc && ev[i].D == markerDst {
		i++
	}

	needle := MapEdge{S: markerSrc, D: markerDst}
	for ; i < iMax; i++ {
		if ev.edgeLess(i, needle) {
			continue
		}
		tab = append(tab, *ev[i])
		if uint32(len(tab)) >= max {
			break
		}
	}
	return tab
}

func (m *Map) Init() {
	m.Cells = make(SetOfVertices, 0)
	m.Roads = make(SetOfEdges, 0)
}

func (m *Map) getNextID() uint64 {
	return atomic.AddUint64(&m.nextID, 1)
}

func (m *Map) CellGet(id uint64) *MapVertex {
	return m.Cells.Get(id)
}

func (m *Map) CellHas(id uint64) bool {
	return m.Cells.Has(id)
}

func (m *Map) CellCreate() *MapVertex {
	id := m.getNextID()
	c := &MapVertex{ID: id}
	m.Cells.Add(c)
	return c
}

// Raw creation of an edge, with no check the Source and Destination exist
// The set of roads isn't sorted afterwards
func (m *Map) RoadCreateRaw(src, dst uint64) *MapEdge {
	if src == dst || src == 0 || dst == 0 {
		panic("Invalid Edge parameters")
	}

	e := &MapEdge{src, dst}
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
		return errors.New("MapEdge exists")
	}
	m.Roads.Add(&MapEdge{src, dst})
	return nil
}

func (m *Map) PathNextStep(src, dst uint64) (uint64, error) {
	if src == dst || src == 0 || dst == 0 {
		return 0, errors.New("EINVAL")
	}

	next, ok := m.steps[vector{src, dst}]
	if ok {
		return next, nil
	}
	return 0, errors.New("No route")
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

func (m *Map) Dot() string {
	var sb strings.Builder
	sb.WriteString("digraph g {")
	for _, c := range m.Cells {
		sb.WriteString("n" + strconv.FormatUint(c.ID, 10))
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
		for _, next := range m.CellAdjacency(cell.ID) {
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
					add(cell.ID, next, first)
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
