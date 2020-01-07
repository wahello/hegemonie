// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/pongo2"
	"github.com/go-macaron/session"
	. "github.com/jfsmig/hegemonie/common/client"
	"github.com/jfsmig/hegemonie/common/mapper"
	"gopkg.in/macaron.v1"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type FormLogin struct {
	UserMail string `form:"email" binding:"Required"`
	UserPass string `form:"password" binding:"Required"`
}

type FormCityStudy struct {
	CharacterId uint64 `form:"cid" binding:"Required"`
	CityId      uint64 `form:"lid" binding:"Required"`
	KnowledgeId uint64 `form:"kid" binding:"Required"`
}

type FormCityBuild struct {
	CharacterId uint64 `form:"cid" binding:"Required"`
	CityId      uint64 `form:"lid" binding:"Required"`
	BuildingId  uint64 `form:"bid" binding:"Required"`
}

type FormCityTrain struct {
	CharacterId uint64 `form:"cid" binding:"Required"`
	CityId      uint64 `form:"lid" binding:"Required"`
	UnitId      uint64 `form:"uid" binding:"Required"`
}

type front struct {
	endpointNorth string
	endpointWorld string
	dirTemplates  string
	dirStatic     string

	region *RegionClientTcp
}

func utoa(u uint64) string {
	return strconv.FormatUint(u, 10)
}

func atou(s string) uint64 {
	u, err := strconv.ParseUint(s, 10, 63)
	if err != nil {
		return 0
	} else {
		return u
	}
}

func ptou(p interface{}) uint64 {
	if p == nil {
		return 0
	}
	return atou(p.(string))
}

func (f *front) routePages(m *macaron.Macaron) {
	m.Get("/",
		func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
			ctx.Data["Title"] = "Hegemonie"
			ctx.Data["userid"] = sess.Get("userid")
			ctx.HTML(200, "index")
		})
	m.Get("/admin",
		func(ctx *macaron.Context) {

		})
	m.Get("/game/user",
		func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
			// Validate the input
			userid := ptou(sess.Get("userid"))
			if userid == 0 {
				flash.Error("Invalid session")
				ctx.Redirect("/")
				return
			}

			// Query the World server for the user
			args := UserShowArgs{UserId: userid}
			reply := UserShowReply{}
			err := f.region.UserShow(&args, &reply)
			if err != nil {
				flash.Warning("Backend error error: " + err.Error())
				ctx.Redirect("/game/user")
			} else {
				ctx.Data["Title"] = reply.View.Name
				ctx.Data["userid"] = utoa(userid)
				ctx.Data["User"] = &reply.View
				ctx.HTML(200, "user")
			}
		})
	m.Get("/game/character",
		func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
			// Validate the input
			userid := ptou(sess.Get("userid"))
			charid := atou(ctx.Query("cid"))
			if userid == 0 || charid == 0 {
				flash.Error("Invalid session")
				ctx.Redirect("/")
				return
			}

			// Query the World server for the Character
			args := CharacterShowArgs{UserId: userid, CharacterId: charid}
			reply := CharacterShowReply{}
			err := f.region.CharacterShow(&args, &reply)
			if err != nil {
				flash.Warning("Backend error: " + err.Error())
				ctx.Redirect("/game/user")
			} else {
				ctx.Data["Title"] = reply.View.Name
				ctx.Data["userid"] = utoa(userid)
				ctx.Data["cid"] = utoa(charid)
				ctx.Data["Character"] = &reply.View
				ctx.HTML(200, "character")
			}
		})
	m.Get("/game/land",
		func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
			// Validate the input
			userid := ptou(sess.Get("userid"))
			charid := atou(ctx.Query("cid"))
			landid := atou(ctx.Query("lid"))
			if userid == 0 || charid == 0 || landid == 0 {
				flash.Error("Invalid session")
				ctx.Redirect("/")
				return
			}

			// Query the World server for the Character
			args := CityShowArgs{UserId: userid, CharacterId: charid, CityId: landid}
			reply := CityShowReply{}
			err := f.region.CityShow(&args, &reply)
			if err != nil {
				flash.Warning("Character error: " + err.Error())
				ctx.Redirect("/game/user")
			} else {
				ctx.Data["Title"] = reply.View.Name
				ctx.Data["userid"] = utoa(userid)
				ctx.Data["cid"] = utoa(charid)
				ctx.Data["lid"] = utoa(landid)
				ctx.Data["Land"] = &reply.View
				ctx.HTML(200, "land")
			}
		})
	m.Get("/game/map",
		func(ctx *macaron.Context, s session.Store) {
			// gameMap, overlay, err := mapper.Generate()
			// if err != nil {
			// 	ctx.Resp.WriteHeader(500)
			// 	return
			// }
			// ctx.Data["map"] = gameMap
			// ctx.Data["overlay"] = overlay
			// ctx.HTML(200, "map")

			// TODO: VDO: handle error
			resp, _ := http.Get("http://" + f.endpointWorld + "/world/places")
			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				// Backend error
				ctx.Resp.WriteHeader(503)
				return
			}
			mapBytes, _ := ioutil.ReadAll(resp.Body)

			resp2, _ := http.Get("http://" + f.endpointWorld + "/world/cities")
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
		})
	// TODO: VDO: disable these routes when DEBUG=false
	m.Get("/debug/map/map",
		func(ctx *macaron.Context, s session.Store) {
			gameMap, _, err := mapper.Generate()
			if err != nil {
				ctx.Resp.WriteHeader(500)
				return
			}
			ctx.JSON(200, gameMap)
		})
	m.Get("/debug/map/overlay",
		func(ctx *macaron.Context, s session.Store) {
			_, overlay, err := mapper.Generate()
			if err != nil {
				ctx.Resp.WriteHeader(500)
				return
			}
			ctx.JSON(200, overlay)
		})
}

func (f *front) routeForms(m *macaron.Macaron) {
	doLogIn := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormLogin) {
		// Cleanup a previous session
		sess.Flush()

		// Authenticate the user by the world-server
		reply := AuthReply{}
		args := AuthArgs{UserMail: info.UserMail, UserPass: info.UserPass}
		err := f.region.Auth(&args, &reply)
		if err != nil {
			flash.Error("Authentication error: " + err.Error())
			ctx.Redirect("/")
		} else {
			// Establish a session for the user
			strid := utoa(reply.Id)
			ctx.SetSecureCookie("session", strid)
			sess.Set("userid", strid)
			ctx.Redirect("/game/user")
		}
	}

	doLogOut := func(ctx *macaron.Context, s session.Store) {
		ctx.SetSecureCookie("session", "")
		s.Flush()
		ctx.Redirect("/")
	}

	doMove := func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		err := f.region.RoundMove(&RoundMoveArgs{}, &RoundMoveReply{})
		if err != nil {
			flash.Error("Action error: " + err.Error())
		}
		ctx.Redirect("/game/user")
	}

	doProduce := func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		err := f.region.RoundProduce(&RoundProduceArgs{}, &RoundProduceReply{})
		if err != nil {
			flash.Error("Action error: " + err.Error())
		}
		ctx.Redirect("/game/user")
	}

	doCityStudy := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityStudy) {
		reply := CityStudyReply{}
		args := CityStudyArgs{
			UserId:      ptou(sess.Get("userid")),
			CharacterId: info.CharacterId,
			CityId:      info.CityId,
			KnowledgeId: info.KnowledgeId,
		}
		err := f.region.CityStudy(&args, &reply)
		if err != nil {
			flash.Error("Action error: " + err.Error())
		}
		ctx.Redirect("/game/land?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}

	doCityBuild := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityBuild) {
		reply := CityBuildReply{}
		args := CityBuildArgs{
			UserId:      ptou(sess.Get("userid")),
			CharacterId: info.CharacterId,
			CityId:      info.CityId,
			BuildingId:  info.BuildingId,
		}
		err := f.region.CityBuild(&args, &reply)
		if err != nil {
			flash.Error("Action error: " + err.Error())
		}
		ctx.Redirect("/game/land?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}

	doCityTrain := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityTrain) {
		reply := CityTrainReply{}
		args := CityTrainArgs{
			UserId:      ptou(sess.Get("userid")),
			CharacterId: info.CharacterId,
			CityId:      info.CityId,
			UnitId:      info.UnitId,
		}
		err := f.region.CityTrain(&args, &reply)
		if err != nil {
			flash.Error("Action error: " + err.Error())
		}
		ctx.Redirect("/game/land?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}

	m.Post("/action/login", binding.Bind(FormLogin{}), doLogIn)
	m.Post("/action/logout", doLogOut)
	m.Get("/action/logout", doLogOut)
	m.Post("/action/move", doMove)
	m.Post("/action/produce", doProduce)
	m.Post("/action/city/study", binding.Bind(FormCityStudy{}), doCityStudy)
	m.Post("/action/city/build", binding.Bind(FormCityBuild{}), doCityBuild)
	m.Post("/action/city/train", binding.Bind(FormCityTrain{}), doCityTrain)
}

func (f *front) routeMiddlewares(m *macaron.Macaron) {
	// TODO(jfs): The secret has to be shared among all the running instances
	m.SetDefaultCookieSecret(randomSecret())
	m.Use(macaron.Static(f.dirStatic, macaron.StaticOptions{
		Prefix: "static",
	}))
	m.Use(pongo2.Pongoer(pongo2.Options{
		Directory:       f.dirTemplates,
		Extensions:      []string{".tpl", ".html", ".tmpl"},
		HTMLContentType: "text/html",
		Charset:         "UTF-8",
		IndentJSON:      true,
		IndentXML:       true,
	}))
	m.Use(session.Sessioner())
	m.Use(func(ctx *macaron.Context, s session.Store) {
		auth := func() {
			uid := s.Get("userid")
			if uid == "" {
				ctx.Redirect("/index.html")
			}
		}
		// Pages under the /game/* prefix require an established authentication
		switch {
		case strings.HasPrefix(ctx.Req.URL.Path, "/game/"),
			strings.HasPrefix(ctx.Req.URL.Path, "/action/"):
			auth()
		}
	})
}

func randomSecret() string {
	var sb strings.Builder
	sb.WriteString(strconv.FormatInt(time.Now().UnixNano(), 16))
	sb.WriteRune('-')
	sb.WriteString(strconv.FormatUint(uint64(rand.Uint32()), 16))
	sb.WriteRune('-')
	sb.WriteString(strconv.FormatUint(uint64(rand.Uint32()), 16))
	return sb.String()
}

func main() {
	var err error
	var f front

	flag.StringVar(&f.endpointNorth, "north", "127.0.0.1:8080", "TCP/IP North endpoint")
	flag.StringVar(&f.endpointWorld, "world", "127.0.0.1:8081", "World Server to be contacted")
	flag.StringVar(&f.dirTemplates, "templates", "/var/lib/hegemonie/templates", "Directory with the HTML tmeplates")
	flag.StringVar(&f.dirStatic, "static", "/var/lib/hegemonie/static", "Directory with the static files")
	flag.Parse()

	m := macaron.Classic()
	f.region = DialClientTcp(f.endpointWorld)
	f.routeMiddlewares(m)
	f.routeForms(m)
	f.routePages(m)

	err = http.ListenAndServe(f.endpointNorth, m)
	if err != nil {
		log.Printf("Server error: %s", err.Error())
	}
}
