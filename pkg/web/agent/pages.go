// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"context"
	"fmt"
	"github.com/go-macaron/session"
	auth "github.com/jfsmig/hegemonie/pkg/auth/proto"
	region "github.com/jfsmig/hegemonie/pkg/region/proto"
	"gopkg.in/macaron.v1"
)

type ActionPage func(*macaron.Context, session.Store, *session.Flash)

type NoFlashPage func(*macaron.Context, session.Store)

type StatelessPage func(*macaron.Context)

func (f *FrontService) routePages(m *macaron.Macaron) {
	m.Get("/", serveRoot)
	m.Get("/game/admin", serveGameAdmin(f))
	m.Get("/game/user", serveGameUser(f))
	m.Get("/game/character", serveGameCharacter(f))
	m.Get("/game/land/overview", serveGameCityOverview(f))
	m.Get("/game/land/budget", serveGameCityBudget(f))
	m.Get("/game/land/buildings", serveGameCityBuildings(f))
	m.Get("/game/land/armies", serveGameCityArmies(f))
	m.Get("/game/land/units", serveGameCityUnits(f))
	m.Get("/game/land/knowledges", serveGameCityKnowledges(f))
	m.Get("/game/army", serveGameArmyDetail(f))

	m.Get("/map/region", serveRegionMap(f))
	m.Get("/map/cities", serveRegionCities(f))
}

func serveRoot(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	ctx.Data["Title"] = "Hegemonie"
	ctx.Data["userid"] = sess.Get("userid")
	ctx.HTML(200, "index")
}

func serveGameAdmin(f *FrontService) ActionPage {
	return func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		uView, err := f.authenticateAdminFromSession(sess)
		if err != nil {
			flash.Error(err.Error())
			ctx.Redirect("/")
			return
		}
		ctx.Data["Title"] = uView.Name
		ctx.Data["userid"] = utoa(uView.Id)
		ctx.Data["User"] = uView
		ctx.HTML(200, "admin")
	}
}

func serveGameUser(f *FrontService) ActionPage {
	return func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		uView, err := f.authenticateUserFromSession(sess)
		if err != nil {
			flash.Error(err.Error())
			ctx.Redirect("/")
			return
		}

		cliReg := region.NewCityClient(f.cnxRegion)
		for _, c := range uView.Characters {
			l, err := cliReg.List(context.Background(), &region.ListReq{Character: c.Id})
			if err != nil {
				flash.Warning("Error with " + c.Name)
			} else {
				for _, ni := range l.Items {
					c.Cities = append(c.Cities,
						&auth.NamedItem{Id: ni.Id, Name: ni.Name})
				}
			}
		}
		ctx.Data["Title"] = uView.Name
		ctx.Data["userid"] = utoa(uView.Id)
		ctx.Data["User"] = uView
		ctx.HTML(200, "user")
	}
}

func serveGameCharacter(f *FrontService) ActionPage {
	return func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		uView, cView, err := f.authenticateCharacterFromSession(sess, atou(ctx.Query("cid")))
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		// Load the Cities managed by the current Character
		cliReg := region.NewCityClient(f.cnxRegion)
		list, err := cliReg.List(context.Background(), &region.ListReq{Character: cView.Id})
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		// Query the World server for the Character
		ctx.Data["Title"] = uView.Name + "|" + cView.Name
		ctx.Data["userid"] = utoa(uView.Id)
		ctx.Data["User"] = uView
		ctx.Data["cid"] = utoa(cView.Id)
		ctx.Data["Character"] = cView
		ctx.Data["Cities"] = list.Items
		ctx.HTML(200, "character")
	}
}

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

type RawVertex struct {
	Id   uint64 `json:"id"`
	X    uint64 `json:"x"`
	Y    uint64 `json:"y"`
	City uint64 `json:"city"`
}

type RawEdge struct {
	Src uint64 `json:"src"`
	Dst uint64 `json:"dst"`
}

type RawCity struct {
	Id uint64 `json:"id"`
	Name string `json:"name"`
	Cell uint64 `json:"cell"`
}

type RawMap struct {
	Cells map[uint64]RawVertex `json:"cells"`
	Roads []RawEdge            `json:"roads"`
}

func serveRegionMap(f *FrontService) StatelessPage {
	return func(ctx *macaron.Context) {
		id := ctx.Query("id")
		if id != "calaquyr" {
			ctx.Error(400, "Invalid region")
			return
		}

		m := RawMap{
			Cells: make(map[uint64]RawVertex),
			Roads: make([]RawEdge, 0),
		}
		cliReg := region.NewMapClient(f.cnxRegion)

		// FIXME(jfs): iterate in case of a truncated result
		vertices, err := cliReg.Vertices(context.Background(), &region.ListVerticesReq{})
		if err != nil {
			ctx.Error(502, err.Error())
			return
		}
		for _, v := range vertices.Items {
			m.Cells[v.Id] = RawVertex{Id: v.Id, X: v.X, Y: v.Y, City: v.CityId}
		}

		// FIXME(jfs): iterate in case of a truncated result
		edges, err := cliReg.Edges(context.Background(), &region.ListEdgesReq{})
		if err != nil {
			ctx.Error(502, err.Error())
			return
		}
		for _, e := range edges.Items {
			m.Roads = append(m.Roads, RawEdge{Src: e.Src, Dst: e.Dst})
		}

		ctx.JSON(200, m)
	}
}

func serveRegionCities(f *FrontService) StatelessPage {
	return func(ctx *macaron.Context) {
		id := ctx.Query("id")
		if id != "calaquyr" {
			ctx.Error(400, "Invalid region")
			return
		}

		tab := make([]RawCity, 0)
		cli := region.NewMapClient(f.cnxRegion)

		// FIXME(jfs): iterate in case of a truncated result
		cities, err := cli.Cities(context.Background(), &region.CitiesReq{})
		if err != nil {
			ctx.Error(502, err.Error())
			return
		}
		for _, v := range cities.Items {
			tab = append(tab, RawCity{Id: v.Id, Name: v.Name, Cell: v.Cell})
		}

		ctx.JSON(200, tab)
	}
}
