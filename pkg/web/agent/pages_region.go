// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"github.com/go-macaron/session"
	region "github.com/jfsmig/hegemonie/pkg/region/proto"
	"gopkg.in/macaron.v1"
)

type RawVertex struct {
	ID   uint64 `json:"id"`
	X    uint64 `json:"x"`
	Y    uint64 `json:"y"`
	City uint64 `json:"city"`
}

type RawEdge struct {
	Src uint64 `json:"src"`
	Dst uint64 `json:"dst"`
}

type RawCity struct {
	ID   uint64 `json:"id"`
	Name string `json:"name"`
	Cell uint64 `json:"cell"`
}

type RawMap struct {
	Cells map[uint64]RawVertex `json:"cells"`
	Roads []RawEdge            `json:"roads"`
}

func serveRegionMap(f *frontService) NoFlashPage {
	return func(ctx *macaron.Context, sess session.Store) {
		id := ctx.Query("id")
		if id != "calaquyr" {
			ctx.Error(400, "Invalid region")
			return
		}

		m := RawMap{
			Cells: make(map[uint64]RawVertex),
			Roads: make([]RawEdge, 0),
		}
		cli := region.NewMapClient(f.cnxRegion)
		ctx0 := contextMacaronToGrpc(ctx, sess)

		// FIXME(jfs): iterate in case of a truncated result
		vertices, err := f.loadAllLocations(ctx0, cli)
		if err != nil {
			ctx.Error(502, err.Error())
			return
		}
		for _, v := range vertices {
			m.Cells[v.Id] = RawVertex{ID: v.Id, X: v.X, Y: v.Y, City: v.CityId}
		}

		// FIXME(jfs): iterate in case of a truncated result
		edges, err := f.loadAllRoads(ctx0, cli)
		if err != nil {
			ctx.Error(502, err.Error())
			return
		}
		for _, e := range edges {
			m.Roads = append(m.Roads, RawEdge{Src: e.Src, Dst: e.Dst})
		}

		ctx.JSON(200, m)
	}
}

func serveRegionCities(f *frontService) NoFlashPage {
	return func(ctx *macaron.Context, sess session.Store) {
		id := ctx.Query("id")
		if id != "calaquyr" {
			ctx.Error(400, "Invalid region")
			return
		}

		tab := make([]RawCity, 0)
		cli := region.NewMapClient(f.cnxRegion)

		// FIXME(jfs): iterate in case of a truncated result
		cities, err := f.loadAllCities(contextMacaronToGrpc(ctx, sess), cli)
		if err != nil {
			ctx.Error(502, err.Error())
			return
		}
		for _, v := range cities {
			tab = append(tab, RawCity{ID: v.Id, Name: v.Name, Cell: v.Location})
		}

		ctx.JSON(200, tab)
	}
}
