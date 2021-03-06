// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package mapclient

import (
	"github.com/juju/errors"
)

// SiteRaw is the information of a node in the graph representation on a map by MapRaw.
type SiteRaw struct {
	ID uint64 `json:"id"`
	X  uint64 `json:"x"`
	Y  uint64 `json:"y"`
	// City is the name of the city that should exist at the instanciation a region based on the current map.
	// The presence of a City is achieved by a non-empty string in the City field, that will be the name of the city
	City string `json:"city"`
}

// RoadRaw is a minimal repreentation of a directional road on the map.
type RoadRaw struct {
	Src uint64 `json:"src"`
	Dst uint64 `json:"dst"`
}

// MapRaw is a human-unfriendly representation of a Map
// - The Sites are indexed by a unique number
// - The presence of a non-empty name on the site makes it an established City when the map is instantiated.
// - The roads are unidirectional (the MapRaw is a digraph).
// - A road might be duplicated, without effect
type MapRaw struct {
	ID    string    `json:"id"`
	Sites []SiteRaw `json:"sites"`
	Roads []RoadRaw `json:"roads"`
}

func makeRawMap() MapRaw {
	return MapRaw{
		Sites: make([]SiteRaw, 0),
		Roads: make([]RoadRaw, 0),
	}
}

func (mr *MapRaw) extractMemMap() (mapMem, error) {
	var err error
	memMap := makeMemMap()
	memMap.ID = mr.ID
	for _, s := range mr.Sites {
		memMap.Sites[s.ID] = &siteMem{
			Raw:   s,
			Peers: make(map[*siteMem]bool),
		}
	}
	for _, r := range mr.Roads {
		if src, ok := memMap.Sites[r.Src]; !ok {
			err = errors.NotValidf("no such source %v", r.Src)
			break
		} else if dst, ok := memMap.Sites[r.Dst]; !ok {
			err = errors.NotValidf("no such destination %v", r.Dst)
			break
		} else {
			// raw maps are digraphs, mem maps are digraphs... no need to duplicate any road
			src.Peers[dst] = true
		}
	}
	return memMap, err
}
