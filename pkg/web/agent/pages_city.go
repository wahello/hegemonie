// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"fmt"
	"github.com/go-macaron/session"
	region "github.com/jfsmig/hegemonie/pkg/region/proto"
	"gopkg.in/macaron.v1"
)

func expandCityView(f *FrontService, lView *region.CityView) {
	f.rw.RLock()
	defer f.rw.RUnlock()

	for _, item := range lView.Assets.Units {
		item.Type = f.units[item.IdType]
	}
	for _, item := range lView.Assets.Buildings {
		item.Type = f.buildings[item.IdType]
	}
	for _, item := range lView.Assets.Knowledges {
		item.Type = f.knowledge[item.IdType]
	}
	for _, item := range lView.Assets.Armies {
		for _, u := range item.Units {
			u.Type = f.units[u.IdType]
		}
	}
}

func expandArmyView(f *FrontService, aView *region.ArmyView) {
	f.rw.RLock()
	defer f.rw.RUnlock()

	for _, u := range aView.Units {
		u.Type = f.units[u.IdType]
	}
}

func serveGameCityPage(f *FrontService, template string) ActionPage {
	return func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		uView, cView, err := f.authenticateCharacterFromSession(ctx, sess, atou(ctx.Query("cid")))
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		// Load the chosen City
		cliReg := region.NewCityClient(f.cnxRegion)
		lView, err := cliReg.Show(contextMacaronToGrpc(ctx, sess),
			&region.CityId{Character: cView.Id, City: atou(ctx.Query("lid"))})
		if err != nil {
			flash.Warning("Region error: " + err.Error())
			ctx.Redirect("/game/character?cid=" + fmt.Sprint(cView.Id))
			return
		}

		expandCityView(f, lView)

		ctx.Data["Title"] = cView.Name + "|" + lView.Name
		ctx.Data["userid"] = utoa(uView.Id)
		ctx.Data["User"] = uView
		ctx.Data["cid"] = utoa(cView.Id)
		ctx.Data["Character"] = cView
		ctx.Data["lid"] = utoa(lView.Id)
		ctx.Data["Land"] = lView
		ctx.HTML(200, template)
	}
}

func serveGameCityOverview(f *FrontService) ActionPage {
	return serveGameCityPage(f, "land_overview")
}

func serveGameCityBuildings(f *FrontService) ActionPage {
	return serveGameCityPage(f, "land_buildings")
}

func serveGameCityKnowledges(f *FrontService) ActionPage {
	return serveGameCityPage(f, "land_knowledges")
}

func serveGameCityUnits(f *FrontService) ActionPage {
	return serveGameCityPage(f, "land_units")
}

func serveGameCityArmies(f *FrontService) ActionPage {
	return serveGameCityPage(f, "land_armies")
}

type ArmyCommandExpanded struct {
	Order       int
	CommandId   int
	Location    uint64
	CityId      uint64
	ArmyId      uint64
	CityName    string
	ArmyName    string
	CommandName string
}

func serveGameArmyDetail(f *FrontService) ActionPage {
	return func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		cid := atou(ctx.Query("cid"))
		lid := atou(ctx.Query("lid"))
		aid := atou(ctx.Query("aid"))
		url := fmt.Sprintf("/game/army?cid=%d&lid=%d&aid=%d", cid, lid, aid)

		uView, cView, err := f.authenticateCharacterFromSession(ctx, sess, cid)
		if err != nil {
			flash.Warning("Auth error: " + err.Error())
			ctx.Redirect("/game/user")
			return
		}

		ctx0 := contextMacaronToGrpc(ctx, sess)

		// Load the chosen City
		cliReg := region.NewCityClient(f.cnxRegion)
		lView, err := cliReg.Show(ctx0, &region.CityId{Character: cView.Id, City: lid})
		if err != nil {
			flash.Warning("City error: " + err.Error())
			ctx.Redirect(url)
			return
		}
		expandCityView(f, lView)

		// Load the chosen Army
		cliArmy := region.NewArmyClient(f.cnxRegion)
		aView, err := cliArmy.Show(ctx0, &region.ArmyId{Character: cView.Id, City: lView.Id, Army: aid})
		if err != nil {
			flash.Warning("Army error: " + err.Error())
			ctx.Redirect(url)
			return
		}
		expandArmyView(f, aView)

		cmdv := make([]ArmyCommandExpanded, 0)
		// Build a printable list of commands
		if len(aView.Commands) > 0 {
			// Preload the description of the map
			cliMap := region.NewMapClient(f.cnxRegion)
			cities, err := f.loadAllCities(ctx0, cliMap)
			if err != nil {
				flash.Warning("Map error: " + err.Error())
				ctx.Redirect(url)
				return
			}
			locations, err := f.loadAllLocations(ctx0, cliMap)
			if err != nil {
				flash.Warning("Map error: " + err.Error())
				ctx.Redirect(url)
				return
			}
			// Generate a list of ad-hoc structures
			for idx, c := range aView.Commands {
				loc := locations[c.Target]
				city := cities[loc.CityId]
				cmdv = append(cmdv, ArmyCommandExpanded{
					Order:     idx,
					Location:  c.Target,
					ArmyId:    aid,
					CommandId: int(c.Action),
					CityId:    city.Id,
					ArmyName:  aView.Name,
					CityName:  city.Name,
				})
			}
		}

		ctx.Data["Title"] = cView.Name + "|" + lView.Name
		ctx.Data["userid"] = utoa(uView.Id)
		ctx.Data["User"] = uView
		ctx.Data["cid"] = utoa(cView.Id)
		ctx.Data["Character"] = cView
		ctx.Data["lid"] = utoa(lView.Id)
		ctx.Data["Land"] = lView
		ctx.Data["aid"] = utoa(aView.Id)
		ctx.Data["Army"] = aView
		ctx.Data["Commands"] = cmdv

		ctx.HTML(200, "army")
	}
}

func serveGameCityBudget(f *FrontService) ActionPage {
	return serveGameCityPage(f, "land_budget")
}
