// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_region_proto

import (
	"github.com/jfsmig/hegemonie/pkg/region/model"
)

func ShowEvolution(w *region.World, c *region.City) *CityEvolution {
	cv := &CityEvolution{}

	for _, kt := range c.KnowledgeFrontier(w) {
		cv.KFrontier = append(cv.KFrontier, &KnowledgeTypeView{
			Id: kt.Id, Name: kt.Name,
		})
	}
	for _, bt := range c.BuildingFrontier(w) {
		cv.BFrontier = append(cv.BFrontier, &BuildingTypeView{
			Id: bt.Id, Name: bt.Name,
		})
	}
	for _, ut := range c.UnitFrontier(w) {
		cv.UFrontier = append(cv.UFrontier, &UnitTypeView{
			Id: ut.Id, Name: ut.Name,
		})
	}

	return cv
}

func resMult(r region.ResourcesMultiplier) *ResourcesMult {
	rm := ResourcesMult{}
	// Fuck, protobuf has no array of fixed size
	rm.R0 = r[0]
	rm.R1 = r[1]
	rm.R2 = r[2]
	rm.R3 = r[3]
	rm.R4 = r[4]
	rm.R5 = r[5]
	return &rm
}

func resPlus(r region.ResourcesIncrement) *ResourcesPlus {
	rm := ResourcesPlus{}
	// Fuck, protobuf has no array of fixed size
	rm.R0 = r[0]
	rm.R1 = r[1]
	rm.R2 = r[2]
	rm.R3 = r[3]
	rm.R4 = r[4]
	rm.R5 = r[5]
	return &rm
}

func resAbs(r region.Resources) *ResourcesAbs {
	rm := ResourcesAbs{}
	// Fuck, protobuf has no array of fixed size
	rm.R0 = r[0]
	rm.R1 = r[1]
	rm.R2 = r[2]
	rm.R3 = r[3]
	rm.R4 = r[4]
	rm.R5 = r[5]
	return &rm
}

func resMod(r region.ResourceModifiers) *ResourcesMod {
	rm := ResourcesMod{}
	rm.Mult = resMult(r.Mult)
	rm.Plus = resPlus(r.Plus)
	return &rm
}

func ShowProduction(w *region.World, c *region.City) *ProductionView {
	v := &ProductionView{}
	prod := c.GetProduction(w)
	v.Base = resAbs(prod.Base)
	v.Buildings = resMod(prod.Buildings)
	v.Knowledge = resMod(prod.Knowledge)
	v.Troops = resMod(prod.Troops)
	v.Actual = resAbs(prod.Actual)
	return v
}

func ShowStock(w *region.World, c *region.City) *StockView {
	v := &StockView{}
	stock := c.GetStock(w)
	v.Base = resAbs(stock.Base)
	v.Buildings = resMod(stock.Buildings)
	v.Knowledge = resMod(stock.Knowledge)
	v.Troops = resMod(stock.Troops)
	v.Actual = resAbs(stock.Actual)
	v.Usage = resAbs(stock.Usage)
	return v
}

func ShowAssets(w *region.World, c *region.City) *CityAssets {
	v := &CityAssets{}

	for _, k := range c.Knowledges {
		v.Knowledges = append(v.Knowledges, &KnowledgeView{
			Id: k.Id, IdType: k.Type, Ticks: uint32(k.Ticks),
		})
	}
	for _, b := range c.Buildings {
		v.Buildings = append(v.Buildings, &BuildingView{
			Id: b.Id, IdType: b.Type, Ticks: uint32(b.Ticks),
		})
	}
	for _, u := range c.Units {
		v.Units = append(v.Units, &UnitView{
			Id: u.Id, IdType: u.Type, Ticks: uint32(u.Ticks),
		})
	}

	for _, a := range c.Armies() {
		v.Armies = append(v.Armies, &ArmyView{
			Id: a.Id, Name: a.Name})
	}

	return v
}

func Show(w *region.World, c *region.City) *CityView {
	cv := &CityView{
		Id: c.Id, Name: c.Name, Owner: c.Owner, Deputy: c.Deputy,
	}
	cv.Evol = ShowEvolution(w, c)
	cv.Production = ShowProduction(w, c)
	cv.Stock = ShowStock(w, c)
	cv.Assets = ShowAssets(w, c)
	return cv
}
