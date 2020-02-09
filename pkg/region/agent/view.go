// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_region_agent

import (
	"github.com/jfsmig/hegemonie/pkg/region/model"
	proto "github.com/jfsmig/hegemonie/pkg/region/proto"
)

func ShowEvolution(w *region.World, c *region.City) *proto.CityEvolution {
	cv := &proto.CityEvolution{}

	for _, kt := range c.KnowledgeFrontier(w) {
		cv.KFrontier = append(cv.KFrontier, &proto.KnowledgeTypeView{
			Id: kt.Id, Name: kt.Name,
		})
	}
	for _, bt := range c.BuildingFrontier(w) {
		cv.BFrontier = append(cv.BFrontier, &proto.BuildingTypeView{
			Id: bt.Id, Name: bt.Name,
		})
	}
	for _, ut := range c.UnitFrontier(w) {
		cv.UFrontier = append(cv.UFrontier, &proto.UnitTypeView{
			Id: ut.Id, Name: ut.Name,
		})
	}

	return cv
}

// M2P -> Model to Proto
func resMultM2P(r region.ResourcesMultiplier) *proto.ResourcesMult {
	rm := proto.ResourcesMult{}
	// Fuck, protobuf has no array of fixed size
	rm.R0 = r[0]
	rm.R1 = r[1]
	rm.R2 = r[2]
	rm.R3 = r[3]
	rm.R4 = r[4]
	rm.R5 = r[5]
	return &rm
}

// M2P -> Model to Proto
func resPlusM2P(r region.ResourcesIncrement) *proto.ResourcesPlus {
	rm := proto.ResourcesPlus{}
	// Fuck, protobuf has no array of fixed size
	rm.R0 = r[0]
	rm.R1 = r[1]
	rm.R2 = r[2]
	rm.R3 = r[3]
	rm.R4 = r[4]
	rm.R5 = r[5]
	return &rm
}

// M2P -> Model to Proto
func resAbsM2P(r region.Resources) *proto.ResourcesAbs {
	rm := proto.ResourcesAbs{}
	// Fuck, protobuf has no array of fixed size
	rm.R0 = r[0]
	rm.R1 = r[1]
	rm.R2 = r[2]
	rm.R3 = r[3]
	rm.R4 = r[4]
	rm.R5 = r[5]
	return &rm
}

func resAbsP2M(rm *proto.ResourcesAbs) region.Resources {
	r := region.Resources{}
	r[0] = rm.R0
	r[1] = rm.R1
	r[2] = rm.R2
	r[3] = rm.R3
	r[4] = rm.R4
	r[5] = rm.R5
	return r
}

// M2P -> Model to Proto
func resModM2P(r region.ResourceModifiers) *proto.ResourcesMod {
	rm := proto.ResourcesMod{}
	rm.Mult = resMultM2P(r.Mult)
	rm.Plus = resPlusM2P(r.Plus)
	return &rm
}

func ShowProduction(w *region.World, c *region.City) *proto.ProductionView {
	v := &proto.ProductionView{}
	prod := c.GetProduction(w)
	v.Base = resAbsM2P(prod.Base)
	v.Buildings = resModM2P(prod.Buildings)
	v.Knowledge = resModM2P(prod.Knowledge)
	v.Troops = resModM2P(prod.Troops)
	v.Actual = resAbsM2P(prod.Actual)
	return v
}

func ShowStock(w *region.World, c *region.City) *proto.StockView {
	v := &proto.StockView{}
	stock := c.GetStock(w)
	v.Base = resAbsM2P(stock.Base)
	v.Buildings = resModM2P(stock.Buildings)
	v.Knowledge = resModM2P(stock.Knowledge)
	v.Troops = resModM2P(stock.Troops)
	v.Actual = resAbsM2P(stock.Actual)
	v.Usage = resAbsM2P(stock.Usage)
	return v
}

func ShowAssets(w *region.World, c *region.City) *proto.CityAssets {
	v := &proto.CityAssets{}

	for _, k := range c.Knowledges {
		v.Knowledges = append(v.Knowledges, &proto.KnowledgeView{
			Id: k.Id, IdType: k.Type, Ticks: uint32(k.Ticks),
		})
	}
	for _, b := range c.Buildings {
		v.Buildings = append(v.Buildings, &proto.BuildingView{
			Id: b.Id, IdType: b.Type, Ticks: uint32(b.Ticks),
		})
	}
	for _, u := range c.Units {
		v.Units = append(v.Units, &proto.UnitView{
			Id: u.Id, IdType: u.Type, Ticks: uint32(u.Ticks),
		})
	}

	for _, a := range c.Armies() {
		v.Armies = append(v.Armies, &proto.ArmyView{
			Id: a.Id, Name: a.Name})
	}

	return v
}

func ShowCity(w *region.World, c *region.City) *proto.CityView {
	cv := &proto.CityView{
		Id:     c.Id,
		Name:   c.Name,
		Owner:  c.Owner,
		Deputy: c.Deputy,

		Cult:        c.Cult,
		Chaotic:     c.Chaotic,
		Alignment:   c.Alignment,
		EthnicGroup: c.EthnicGroup,

		TickMassacres: c.TicksMassacres,
		Auto:          c.Auto,
		Deleted:       c.Deleted,

		Politics: &proto.CityPolitics{
			Overlord: c.Overlord,
			Lieges:   []uint64{},
		},
	}

	for _, c := range c.Lieges() {
		cv.Politics.Lieges = append(cv.Politics.Lieges, c.Id)
	}

	cv.Evol = ShowEvolution(w, c)
	cv.Production = ShowProduction(w, c)
	cv.Stock = ShowStock(w, c)
	cv.Assets = ShowAssets(w, c)
	return cv
}

func ShowArmy(w *region.World, a *region.Army) *proto.ArmyView {
	view := &proto.ArmyView{}
	view.Id = a.Id
	view.Name = a.Name
	view.Location = a.Cell
	view.Stock = resAbsM2P(a.Stock)
	for _, u := range a.Units {
		view.Units = append(view.Units, ShowUnit(w, u))
	}
	return view
}

func ShowUnit(w *region.World, u *region.Unit) *proto.UnitView {
	view := &proto.UnitView{}
	view.Id = u.Id
	view.Name = ""
	view.Ticks = u.Ticks
	view.Health = u.Health
	return view
}
