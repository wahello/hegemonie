// Copyright (C) 2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"errors"
	"fmt"
)

type SiteRaw struct {
	Id   string
	X, Y float64
	City bool
}

type RoadRaw struct {
	Src, Dst string
}

type MapRaw struct {
	Sites []SiteRaw
	Roads []RoadRaw
}

func makeRawMap() MapRaw {
	return MapRaw{
		Sites: make([]SiteRaw, 0),
		Roads: make([]RoadRaw, 0),
	}
}

func (mr *MapRaw) Finalize() (Map, error) {
	var err error
	m := makeMap()

	for _, s := range mr.Sites {
		m.sites[s.Id] = &Site{
			raw:   s,
			peers: make(map[*Site]bool),
		}
	}
	for _, r := range mr.Roads {
		if src, ok := m.sites[r.Src]; !ok {
			err = errors.New(fmt.Sprintf("No such site [%s]", r.Src))
			break
		} else if dst, ok := m.sites[r.Dst]; !ok {
			err = errors.New(fmt.Sprintf("No such site [%s]", r.Dst))
			break
		} else {
			src.peers[dst] = true
			dst.peers[src] = true
		}
	}
	return m, err
}
