// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"context"
	"fmt"
	"github.com/go-macaron/session"
	"github.com/jfsmig/hegemonie/pkg/auth/proto"
	"github.com/jfsmig/hegemonie/pkg/region/proto_city"
	"gopkg.in/macaron.v1"
	"io/ioutil"
	"net/http"
)

type ActionPage func(*macaron.Context, session.Store, *session.Flash)

type NoFlashPage func(*macaron.Context, session.Store)

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

	m.Get("/game/map/region", serveRegionMap(f))
	m.Get("/game/map/city", serveCityMap(f))
}

func serveRoot(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	ctx.Data["Title"] = "Hegemonie"
	ctx.Data["userid"] = sess.Get("userid")
	ctx.HTML(200, "index")
}

func serveGameAdmin(f *FrontService) ActionPage {
	return func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		// Validate the input
		userid := ptou(sess.Get("userid"))
		if userid == 0 {
			flash.Error("Invalid session")
			ctx.Redirect("/")
			return
		}

		// Authorize the character with the user
		cliAuth := hegemonie_auth_proto.NewAuthClient(f.cnxAuth)
		view, err := cliAuth.UserShow(context.Background(),
			&hegemonie_auth_proto.UserShowReq{Id: userid})
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		if !view.Admin {
			flash.Warning("Insufficient permission")
			ctx.Redirect("/game/user")
			return
		}

		ctx.Data["Title"] = view.Name
		ctx.Data["userid"] = utoa(userid)
		ctx.Data["User"] = view
		//ctx.Data["Score"] = &sReply.Board
		ctx.HTML(200, "admin")
	}
}

func serveGameUser(f *FrontService) ActionPage {
	return func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		// Validate the input
		userid := ptou(sess.Get("userid"))
		if userid == 0 {
			flash.Error("Invalid session")
			ctx.Redirect("/")
			return
		}

		// Authorize the character with the user
		cliAuth := hegemonie_auth_proto.NewAuthClient(f.cnxAuth)
		view, err := cliAuth.UserShow(context.Background(),
			&hegemonie_auth_proto.UserShowReq{Id: userid})

		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/")
			return
		}
		ctx.Data["Title"] = view.Name
		ctx.Data["userid"] = utoa(userid)
		ctx.Data["User"] = view
		ctx.HTML(200, "user")
	}
}

func serveGameCharacter(f *FrontService) ActionPage {
	return func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		var err error

		// Validate the input
		userid := ptou(sess.Get("userid"))
		charid := atou(ctx.Query("cid"))
		if userid == 0 || charid == 0 {
			flash.Error("Invalid session")
			ctx.Redirect("/")
			return
		}

		// Authorize the character with the user
		cliAuth := hegemonie_auth_proto.NewAuthClient(f.cnxAuth)
		view, err := cliAuth.CharacterShow(context.Background(),
			&hegemonie_auth_proto.CharacterShowReq{User: userid, Character: charid})
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		// Load the Cities managed by the current Character
		cliReg := hegemonie_region_proto_city.NewCityClient(f.cnxRegion)
		list, err := cliReg.List(context.Background(),
			&hegemonie_region_proto_city.ListReq{Character: view.Id})
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		// Query the World server for the Character
		ctx.Data["Title"] = view.Name
		ctx.Data["userid"] = utoa(userid)
		ctx.Data["cid"] = utoa(charid)
		ctx.Data["Character"] = view
		ctx.Data["Cities"] = list.Items
		ctx.HTML(200, "character")
	}
}

func serveGameCityPage(f *FrontService, template string) ActionPage {
	return func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		// Validate the input
		userid := ptou(sess.Get("userid"))
		charid := atou(ctx.Query("cid"))
		landid := atou(ctx.Query("lid"))
		if userid == 0 || charid == 0 || landid == 0 {
			flash.Error("Invalid session")
			ctx.Redirect("/")
			return
		}

		// Authorize the character with the user
		cliAuth := hegemonie_auth_proto.NewAuthClient(f.cnxAuth)
		cView, err := cliAuth.CharacterShow(context.Background(),
			&hegemonie_auth_proto.CharacterShowReq{User: userid, Character: charid})
		if err != nil {
			flash.Warning("Auth error: " + err.Error())
			ctx.Redirect("/game/character?cid=" + fmt.Sprint(charid))
			return
		}

		// Load the chosen City
		cliReg := hegemonie_region_proto_city.NewCityClient(f.cnxRegion)
		lView, err := cliReg.Show(context.Background(),
			&hegemonie_region_proto_city.CityId{Character: charid, City: landid})
		if err != nil {
			flash.Warning("Region error: " + err.Error())
			ctx.Redirect("/game/character?cid=" + fmt.Sprint(charid))
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
		for _, item := range lView.Assets.Knowledge {
			item.Type = f.knowledge[item.IdType]
		}
		f.rw.RUnlock()

		ctx.Data["userid"] = utoa(userid)
		ctx.Data["cid"] = utoa(charid)
		ctx.Data["lid"] = utoa(landid)
		ctx.Data["Title"] = lView.Name
		ctx.Data["Character"] = cView
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

func serveGameCityBudget(f *FrontService) ActionPage {
	return serveGameCityPage(f, "land_budget")
}

func serveRegionMap(f *FrontService) NoFlashPage {
	return func(ctx *macaron.Context, s session.Store) {
		// gameMap, overlay, err := mapper.Generate()
		// if err != nil {
		// 	ctx.Resp.WriteHeader(500)
		// 	return
		// }
		// ctx.Data["map"] = gameMap
		// ctx.Data["overlay"] = overlay
		// ctx.HTML(200, "map")

		// TODO: VDO: handle error
		resp, _ := http.Get("http://" + f.endpointRegion + "/cmd_back_region/places")
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			// Backend error
			ctx.Resp.WriteHeader(503)
			return
		}
		mapBytes, _ := ioutil.ReadAll(resp.Body)

		resp2, _ := http.Get("http://" + f.endpointRegion + "/cmd_back_region/cities")
		defer resp2.Body.Close()
		if resp2.StatusCode != http.StatusOK {
			// Backend error
			ctx.Resp.WriteHeader(503)
			return
		}
		mapCities, _ := ioutil.ReadAll(resp2.Body)

		ctx.Data["map"] = string(mapBytes)
		ctx.Data["cities"] = string(mapCities)

		ctx.HTML(200, "map")
	}
}

func serveCityMap(f *FrontService) NoFlashPage {
	return func(ctx *macaron.Context, s session.Store) {
		// gameMap, overlay, err := mapper.Generate()
		// if err != nil {
		// 	ctx.Resp.WriteHeader(500)
		// 	return
		// }
		// ctx.Data["map"] = gameMap
		// ctx.Data["overlay"] = overlay
		// ctx.HTML(200, "map")

		// TODO: VDO: handle error
		resp, _ := http.Get("http://" + f.endpointRegion + "/cmd_back_region/places")
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			// Backend error
			ctx.Resp.WriteHeader(503)
			return
		}
		mapBytes, _ := ioutil.ReadAll(resp.Body)

		resp2, _ := http.Get("http://" + f.endpointRegion + "/cmd_back_region/cities")
		defer resp2.Body.Close()
		if resp2.StatusCode != http.StatusOK {
			// Backend error
			ctx.Resp.WriteHeader(503)
			return
		}
		mapCities, _ := ioutil.ReadAll(resp2.Body)

		ctx.Data["map"] = string(mapBytes)
		ctx.Data["cities"] = string(mapCities)

		ctx.HTML(200, "map")
	}
}
