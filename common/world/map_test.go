// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"os"
	"testing"
	"time"
)

func TestMapInit(t *testing.T) {
	var m Map
	m.Init()

	if m.NodeHas(1) {
		t.Fatal()
	}
	if m.getNextId() != 1 {
		t.Fatal()
	}

	loc, err := m.NodeCreate()
	if err != nil {
		t.Fatal()
	}

	if loc != 2 {
		t.Fatal()
	}
	if m.getNextId() != 3 {
		t.Fatal()
	}
	if !m.NodeHas(loc) {
		t.Fatal()
	}
}

func TestMapEinval(t *testing.T) {
	var m Map
	m.Init()

	// Test that identical, zero or inexistant locations return an error
	for _, src := range []uint64{0, 1, 2} {
		for _, dst := range []uint64{0, 1, 2} {
			if err := m.VertexCreate(src, dst, true); err == nil {
				t.Fatal()
			}
			if err := m.VertexDelete(src, dst, true); err == nil {
				t.Fatal()
			}
		}
	}
}

func TestMapMultiConnect(t *testing.T) {
	var err error
	var m Map
	m.Init()
	l0, _ := m.NodeCreate()
	l1, _ := m.NodeCreate()
	if err = m.VertexCreate(l0, l1, true); err != nil {
		t.Fatal()
	}
	if err = m.VertexCreate(l1, l0, true); err != nil {
		t.Fatal()
	}
	for i := 0; i < 5; i++ {
		if err = m.VertexCreate(l0, l1, true); err == nil {
			t.Logf("Cells %v", m.Cells)
			t.Logf("Roads %v", m.Roads)
			t.Fatal()
		}
		if err = m.VertexCreate(l1, l0, true); err == nil {
			t.Logf("Cells %v", m.Cells)
			t.Logf("Roads %v", m.Roads)
			t.Fatal()
		}
	}
}

func TestMapPathOneWay(t *testing.T) {
	var m Map
	m.Init()

	l0, _ := m.NodeCreate()
	l1, _ := m.NodeCreate()
	l2, _ := m.NodeCreate()
	l3, _ := m.NodeCreate()
	m.VertexCreate(l0, l1, true)
	m.VertexCreate(l1, l2, true)
	m.VertexCreate(l2, l3, true)

	m.Rehash()

	if step, err := m.NodeGetStep(l0, l3); err != nil {
		t.Fatal()
	} else if step != l1 {
		t.Fatal()
	}

	if step, err := m.NodeGetStep(l1, l3); err != nil {
		t.Fatal()
	} else if step != l2 {
		t.Fatal()
	}

	if step, err := m.NodeGetStep(l2, l3); err != nil {
		t.Fatal()
	} else if step != l3 {
		t.Fatal()
	}

	if step, err := m.NodeGetStep(l1, l0); err == nil {
		t.Fatal()
	} else if step != 0 {
		t.Fatal()
	}
}

func TestMapPathTwoWay(t *testing.T) {
	var m Map
	m.Init()

	l0, _ := m.NodeCreate()
	l1, _ := m.NodeCreate()
	l2, _ := m.NodeCreate()
	l3, _ := m.NodeCreate()
	biconnect := func(l0, l1 uint64) {
		m.VertexCreate(l0, l1, false)
		m.VertexCreate(l1, l0, false)
	}

	biconnect(l0, l1)
	biconnect(l1, l2)
	biconnect(l2, l3)

	m.Rehash()

	if step, err := m.NodeGetStep(l3, l0); err != nil {
		t.Fatal()
	} else if step != l2 {
		t.Fatal()
	}

	if step, err := m.NodeGetStep(l1, l3); err != nil {
		t.Fatal()
	} else if step != l2 {
		t.Fatal()
	}

	if step, err := m.NodeGetStep(l2, l3); err != nil {
		t.Fatal()
	} else if step != l3 {
		t.Fatal()
	}
}

type grid struct {
	tab []uint64
	x   int
	y   int
}

func newGrid(x, y int) *grid {
	g := grid{x: x, y: y}
	g.tab = make([]uint64, x*y, x*y)
	return &g
}

func (g *grid) loc(i, j int) int {
	return i*g.y + j
}

func (g *grid) get(i, j int) uint64 {
	return g.tab[g.loc(i, j)]
}

func (g *grid) set(i, j int, v uint64) {
	g.tab[g.loc(i, j)] = v
}

func TestMapGrid(t *testing.T) {
	var m Map
	m.Init()

	x := 10
	y := 20
	t.Logf("Building the grid at %v", time.Now())
	grid := newGrid(x, y)
	for i := 0; i < x; i++ {
		for j := 0; j < y; j++ {
			v, err := m.NodeCreate()
			if err == nil {
				grid.set(i, j, v)
			}
		}
	}

	offsets := []int{-1, 0, 1}
	neighbourhood := func(i, j int) {
		src := grid.get(i, j)
		// Iterate on the neighbors
		for _, id := range offsets {
			for _, jd := range offsets {
				if id == 0 && jd == 0 {
					// No self route
					continue
				}
				if id != 0 && jd != 0 {
					// No diagonals
					continue
				}
				if i+id < 0 || j+jd < 0 {
					// No underflow
					continue
				}
				if i+id >= x || j+jd >= y {
					// No overflow
					continue
				}
				dst := grid.get(i+id, j+jd)
				m.VertexCreateNoCheck(src, dst)
				m.VertexCreateNoCheck(dst, src)
			}
		}
	}

	// Even rows
	for i := 0; i < x; i += 2 {
		for j := 0; j < y; j += 2 {
			neighbourhood(i, j)
		}
	}
	// Odd rows
	for i := 1; i < x; i += 2 {
		for j := 1; j < y; j += 2 {
			neighbourhood(i, j)
		}
	}

	t.Logf("Rehashing at %v", time.Now())
	m.Rehash()

	t.Logf("Testing at %v", time.Now())

	dot := m.Dot()
	f, _ := os.Create("/tmp/dot")
	f.WriteString(dot)
	f.Close()

	t.Logf("Done at %v", time.Now())
}
