// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapclient

import (
	"fmt"
	"sort"
)

type SiteSeed struct {
	ID   string `json:"id"`
	X    uint64 `json:"x"`
	Y    uint64 `json:"y"`
	City bool   `json:"city"`
}

type RoadSeed struct {
	Src string `json:"src"`
	Dst string `json:"dst"`
}

// Handy representation of a map where nearly each node of the graph carries a city.
// - The sites are identified by a textual name, typically decided by a content creator
// that will name the cities in a meaningful way.
// - Each site name MUST be unique
// - The reads bidirectional.
// - Duplicated roads are treated as is, no simplification is done.
type MapSeed struct {
	ID    string     `json:"id"`
	Sites []SiteSeed `json:"sites"`
	Roads []RoadSeed `json:"roads"`
}

func (s SiteSeed) cityName() string {
	if s.City {
		return s.ID
	}
	return ""
}

// extractRawMap to a extractRawMap mapMem. extractRawMap maps are very similar to mapMem Seeds, they mostly differ
// by the kind/type of index used for the Sites (vertices): mapMem Seeds use strings, extractRawMap Maps
// use sequential indexes.
func (ms *MapSeed) extractRawMap() (MapRaw, error) {
	var err error
	byName := make(map[string]uint64)

	rawMap := makeRawMap()
	rawMap.ID = ms.ID

	// Sorted Sites / Vertices
	sort.Slice(ms.Sites, func(i, j int) bool { return ms.Sites[i].ID < ms.Sites[j].ID })
	for idx, s := range ms.Sites {
		// We need a non-zero unique ID that is monotonically increasing
		id := uint64(idx) + 1
		sr := SiteRaw{ID: id, X: s.X, Y: s.Y, City: s.cityName()}
		rawMap.Sites = append(rawMap.Sites, sr)
		byName[s.ID] = id
	}

	// Sorted Roads / Edges
	for i, r := range ms.Roads {
		src, dst := byName[r.Src], byName[r.Dst]
		if src <= 0 {
			err = fmt.Errorf("No such site [%s] [%v:%v->%v]", r.Src, i, r.Src, r.Dst)
			break
		} else if dst <= 0 {
			err = fmt.Errorf("No such site [%s] [%v:%v->%v]", r.Dst, i, r.Src, r.Dst)
			break
		} else {
			// map seeds are graphs, raw maps are digraphs ... we need to expand in directions
			rawMap.Roads = append(rawMap.Roads, RoadRaw{src, dst})
			rawMap.Roads = append(rawMap.Roads, RoadRaw{dst, src})
		}
	}
	sort.Slice(rawMap.Roads, func(i, j int) bool {
		ri, rj := rawMap.Roads[i], rawMap.Roads[j]
		return ri.Src < rj.Src || (ri.Src == rj.Src && ri.Dst < rj.Dst)
	})

	return rawMap, err
}
