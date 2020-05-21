// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"context"
	"fmt"
	"github.com/go-macaron/session"
	region "github.com/jfsmig/hegemonie/pkg/region/proto"
	"gopkg.in/macaron.v1"
)

func serveGameCityPage(f *FrontService, template string) ActionPage {
	return func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		uView, cView, err := f.authenticateCharacterFromSession(sess, atou(ctx.Query("cid")))
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		// Load the chosen City
		cliReg := region.NewCityClient(f.cnxRegion)
		lView, err := cliReg.Show(context.Background(),
			&region.CityId{Character: cView.Id, City: atou(ctx.Query("lid"))})
		if err != nil {
			flash.Warning("Region error: " + err.Error())
			ctx.Redirect("/game/character?cid=" + fmt.Sprint(cView.Id))
			return
		}

		// Expand the view
		f.rw.RLock()
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
		f.rw.RUnlock()

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

func serveGameArmyDetail(f *FrontService) ActionPage {
	return func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		uView, cView, err := f.authenticateCharacterFromSession(sess, atou(ctx.Query("cid")))
		if err != nil {
			flash.Warning("Auth error: " + err.Error())
			ctx.Redirect("/game/user")
			return
		}

		// Load the chosen City
		cliReg := region.NewCityClient(f.cnxRegion)
		lView, err := cliReg.Show(context.Background(),
			&region.CityId{Character: cView.Id, City: atou(ctx.Query("lid"))})
		if err != nil {
			flash.Warning("City error: " + err.Error())
			ctx.Redirect(fmt.Sprintf("/game/land/armies?cid=%d&lid=%d", cView.Id, lView.Id))
			return
		}

		// Load the chosen Army
		cliArmy := region.NewArmyClient(f.cnxRegion)
		aView, err := cliArmy.Show(context.Background(),
			&region.ArmyId{Character: cView.Id, City: lView.Id, Army: atou(ctx.Query("aid"))})
		if err != nil {
			flash.Warning("Army error: " + err.Error())
			ctx.Redirect(fmt.Sprintf("/game/land/armies?cid=%d&lid=%d", cView.Id, lView.Id))
			return
		}

		// Expand the view
		f.rw.RLock()
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
		for _, u := range aView.Units {
			u.Type = f.units[u.IdType]
		}
		f.rw.RUnlock()

		ctx.Data["Title"] = cView.Name + "|" + lView.Name
		ctx.Data["userid"] = utoa(uView.Id)
		ctx.Data["User"] = uView
		ctx.Data["cid"] = utoa(cView.Id)
		ctx.Data["Character"] = cView
		ctx.Data["lid"] = utoa(lView.Id)
		ctx.Data["Land"] = lView
		ctx.Data["aid"] = utoa(aView.Id)
		ctx.Data["Army"] = aView

		ctx.HTML(200, "army")
	}
}

func serveGameCityBudget(f *FrontService) ActionPage {
	return serveGameCityPage(f, "land_budget")
}
