// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapclient

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync/atomic"
)

type siteMem struct {
	Raw   SiteRaw
	Peers map[*siteMem]bool
}

type roadMem struct {
	Src, Dst *siteMem
}

// Human-unfriendly representation of a Map
// - The Sites are indexed by a unique number
// - The roads are bidirectional.
// - Road can be duplicated
type mapMem struct {
	ID     string
	Sites  map[uint64]*siteMem
	nextID uint64
}

func makeMemMap() mapMem {
	return mapMem{
		Sites: make(map[uint64]*siteMem),
	}
}

func makeSite(raw SiteRaw) *siteMem {
	return &siteMem{
		Raw:   raw,
		Peers: make(map[*siteMem]bool),
	}
}

func (s *siteMem) getDotName() string {
	if s.Raw.City != "" {
		return s.Raw.City
	}
	return fmt.Sprintf("x%v", s.Raw.ID)
}

func (m *mapMem) uniqueRoads() <-chan roadMem {
	out := make(chan roadMem)
	go func() {
		seen := make(map[RoadRaw]bool)
		for _, s := range m.Sites {
			for peer := range s.Peers {
				r0 := RoadRaw{Src: s.Raw.ID, Dst: peer.Raw.ID}
				if !seen[r0] {
					seen[r0] = true
					out <- roadMem{s, peer}
				}
			}
		}
		close(out)
	}()
	return out
}

func (m *mapMem) sortedSites() <-chan *siteMem {
	out := make(chan *siteMem)
	go func() {
		keys := make([]uint64, 0, len(m.Sites))
		for k := range m.Sites {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
		for _, k := range keys {
			out <- m.Sites[k]
		}
		close(out)
	}()
	return out
}

// Produces a MapRaw with sorted sites and unique roads.
func (m *mapMem) extractRawMap() MapRaw {
	rawMap := makeRawMap()
	rawMap.ID = m.ID
	for s := range m.sortedSites() {
		rawMap.Sites = append(rawMap.Sites, s.Raw)
	}
	for r := range m.uniqueRoads() {
		rawRoad := RoadRaw{Src: r.Src.Raw.ID, Dst: r.Dst.Raw.ID}
		rawMap.Roads = append(rawMap.Roads, rawRoad)
	}
	return rawMap
}

func (m *mapMem) deepCopy() mapMem {
	mFinal := makeMemMap()
	for id, site := range m.Sites {
		mFinal.Sites[id] = makeSite(site.Raw)
	}
	for _, s := range m.Sites {
		src := mFinal.Sites[s.Raw.ID]
		for d := range s.Peers {
			dst := mFinal.Sites[d.Raw.ID]
			src.Peers[dst] = true
			dst.Peers[src] = true
		}
	}
	return mFinal
}

func (m *mapMem) computeBox() (xmin, xmax, ymin, ymax uint64) {
	const Max = math.MaxUint64
	const Min = 0
	xmin, ymin, xmax, ymax = Max, Max, Min, Min
	for _, s := range m.Sites {
		x, y := s.Raw.X, s.Raw.Y
		if x < xmin {
			xmin = x
		}
		if x > xmax {
			xmax = x
		}
		if y < ymin {
			ymin = y
		}
		if y > ymax {
			ymax = y
		}
	}
	if xmin == Max {
		xmin, xmax, ymin, ymax = 0, 0, 0, 0
	}
	return
}

func (m *mapMem) shiftAt(xabs, yabs uint64) {
	xmin, _, ymin, _ := m.computeBox()
	m.shift(xabs-xmin, yabs-ymin)
}

func (m *mapMem) shift(xrel, yrel uint64) {
	for _, s := range m.Sites {
		s.Raw.X += xrel
		s.Raw.Y += yrel
	}
}

func (m *mapMem) resizeRatio(xratio, yratio float64) {
	for _, s := range m.Sites {
		s.Raw.X = uint64(math.Round(float64(s.Raw.X) * xratio))
		s.Raw.Y = uint64(math.Round(float64(s.Raw.Y) * yratio))
	}
}

func (m *mapMem) resizeStretch(x, y float64) {
	m.shiftAt(0, 0)
	_, xmax, _, ymax := m.computeBox()
	m.resizeRatio(x/float64(xmax), y/float64(ymax))
}

func (m *mapMem) resizeAndAdjust(x, y uint64) {
	m.shiftAt(0, 0)
	_, xmax, _, ymax := m.computeBox()
	xRatio := float64(x) / float64(xmax)
	yRatio := float64(y) / float64(ymax)
	ratio := math.Min(xRatio, yRatio)
	m.resizeRatio(ratio, ratio)
}

func (m *mapMem) SiftToTheCenter(xbound, ybound uint64) {
	xmin, xmax, ymin, ymax := m.computeBox()
	xdelta, ydelta := xbound-(xmax-xmin), ybound-(ymax-ymin)
	xpad, ypad := xdelta/2.0, ydelta/2.0
	m.shift(xpad-xmin, ypad-ymin)
}

func (m *mapMem) splitOneRoad(src, dst *siteMem, nbSegments uint) {
	if nbSegments < 2 {
		panic("bug")
	}

	xinc := uint64(math.Round(float64(dst.Raw.X-src.Raw.X) / float64(nbSegments)))
	yinc := uint64(math.Round(float64(dst.Raw.Y-src.Raw.Y) / float64(nbSegments)))
	segments := make([]*siteMem, 0, nbSegments+1)

	delete(src.Peers, dst)
	delete(dst.Peers, src)

	// Create segment boundaries
	segments = append(segments, src)
	for i := uint(0); i < nbSegments-1; i++ {
		last := segments[len(segments)-1]
		x := last.Raw.X + xinc
		y := last.Raw.Y + yinc
		id := atomic.AddUint64(&m.nextID, 1)
		raw := SiteRaw{ID: id, City: "", X: x, Y: y}
		middle := makeSite(raw)
		m.Sites[middle.Raw.ID] = middle
		segments = append(segments, middle)
	}
	segments = append(segments, dst)

	// Link the segment boundaries
	for i, end := range segments[1:] {
		start := segments[i]
		start.Peers[end] = true
		end.Peers[start] = true
	}
}

func (m *mapMem) splitLongRoads(max float64) mapMem {
	// Work on a deep copy to iterate on the original map while we alter the copy
	mCopy := m.deepCopy()
	for r := range m.uniqueRoads() {
		src := mCopy.Sites[r.Src.Raw.ID]
		dst := mCopy.Sites[r.Dst.Raw.ID]
		dist := distance(src, dst)
		if max < dist {
			mCopy.splitOneRoad(src, dst, uint(math.Ceil(dist/max)))
		}
	}
	return mCopy
}

func (m *mapMem) applyNoiseOnPositions(xjitter, yjitter float64) {
	for _, s := range m.Sites {
		if s.Raw.City != "" {
			continue
		}
		s.Raw.X += uint64(math.Round((0.5 - rand.Float64()) * xjitter))
		s.Raw.Y += uint64(math.Round((0.5 - rand.Float64()) * yjitter))
	}
}

func distance(src, dst *siteMem) float64 {
	dx := (dst.Raw.X - src.Raw.X)
	dy := (dst.Raw.Y - src.Raw.Y)
	return math.Sqrt(float64(dx*dx) + float64(dy*dy))
}
