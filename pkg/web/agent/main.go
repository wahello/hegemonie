// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/pongo2"
	"github.com/go-macaron/session"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"gopkg.in/macaron.v1"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jfsmig/hegemonie/pkg/auth/proto"
	mapper "github.com/jfsmig/hegemonie/pkg/mapper"
	"github.com/jfsmig/hegemonie/pkg/region/proto_city"
)

func Command() *cobra.Command {
	front := FrontService{}
	agent := &cobra.Command{
		Use:     "agent",
		Aliases: []string{"srv", "server", "service", "worker"},
		Short:   "Web service",
		RunE: func(cmd *cobra.Command, args []string) error {
			if front.endpointRegion == "" {
				return errors.New("Missing region URL")
			}

			if fi, err := os.Stat(front.dirTemplates); err != nil || !fi.IsDir() {
				return errors.New("Invalid path for the directory of templates")
			}
			if fi, err := os.Stat(front.dirStatic); err != nil || !fi.IsDir() {
				return errors.New("Invalid path for the directory of static files")
			}

			m := macaron.Classic()
			front.routeMiddlewares(m)
			front.routeForms(m)
			front.routePages(m)

			var err error

			front.cnxAuth, err = grpc.Dial(front.endpointAuth, grpc.WithInsecure())
			if err != nil {
				return err
			}

			front.cnxRegion, err = grpc.Dial(front.endpointRegion, grpc.WithInsecure())
			if err != nil {
				return err
			}

			return http.ListenAndServe(front.endpointNorth, m)
		},
	}
	agent.Flags().StringVar(&front.endpointNorth, "endpoint", ":8080", "TCP/IP North endpoint")
	agent.Flags().StringVar(&front.endpointRegion, "region", "", "World Server to be contacted")
	agent.Flags().StringVar(&front.endpointAuth, "auth", "", "Auth Server to be contacted")
	agent.Flags().StringVar(&front.dirTemplates, "templates", "/data/templates", "Directory with the HTML templates")
	agent.Flags().StringVar(&front.dirStatic, "static", "/data/static", "Directory with the static files")
	return agent
}

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

type FormCityUnitTransfer struct {
	CharacterId uint64 `form:"cid" binding:"Required"`
	CityId      uint64 `form:"lid" binding:"Required"`
	UnitId      uint64 `form:"uid" binding:"Required"`
	ArmyId      uint64 `form:"aid" binding:"Required"`
}

type FormCityArmyCreate struct {
	CharacterId uint64 `form:"cid" binding:"Required"`
	CityId      uint64 `form:"lid" binding:"Required"`
	Name        string `form:"name" binding:"Required"`
}

type FormCityArmyCommand struct {
	CharacterId uint64 `form:"cid" binding:"Required"`
	CityId      uint64 `form:"lid" binding:"Required"`
	ArmyId      uint64 `form:"aid" binding:"Required"`

	Cell   uint64 `form:"cell" binding:"Required"`
	Action uint64 `form:"what" binding:"Required"`
}

type FrontService struct {
	dirTemplates   string
	dirStatic      string
	endpointNorth  string
	endpointRegion string
	endpointAuth   string

	cnxRegion *grpc.ClientConn
	cnxAuth   *grpc.ClientConn
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

func serveRoot(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	ctx.Data["Title"] = "Hegemonie"
	ctx.Data["userid"] = sess.Get("userid")
	ctx.HTML(200, "index")
}

type ActionPage func(*macaron.Context, session.Store, *session.Flash)
type NoFlashPage func(*macaron.Context, session.Store)

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

func serveGameCity(f *FrontService) ActionPage {
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

		ctx.Data["userid"] = utoa(userid)
		ctx.Data["cid"] = utoa(charid)
		ctx.Data["lid"] = utoa(landid)
		ctx.Data["Title"] = lView.Name
		ctx.Data["Character"] = cView
		ctx.Data["Land"] = lView
		ctx.HTML(200, "land")
	}
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

func (f *FrontService) routePages(m *macaron.Macaron) {
	m.Get("/", serveRoot)
	m.Get("/game/admin", serveGameAdmin(f))
	m.Get("/game/user", serveGameUser(f))
	m.Get("/game/character", serveGameCharacter(f))
	m.Get("/game/land", serveGameCity(f))

	m.Get("/game/map/region", serveRegionMap(f))
	m.Get("/game/map/city", serveCityMap(f))

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

func (f *FrontService) routeForms(m *macaron.Macaron) {
	doLogIn := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormLogin) {
		// Cleanup a previous session
		sess.Flush()

		// Authorize the character with the user
		cliAuth := hegemonie_auth_proto.NewAuthClient(f.cnxAuth)
		view, err := cliAuth.UserAuth(context.Background(),
			&hegemonie_auth_proto.UserAuthReq{Mail: info.UserMail, Pass: info.UserPass})

		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/")
		} else {
			strid := utoa(view.Id)
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
		/*
			err := f.region.RoundMove(&hclient.RoundMoveArgs{}, &hclient.RoundMoveReply{})
			if err != nil {
				flash.Error("Action error: " + err.Error())
			}
		*/
		ctx.Redirect("/game/user")
	}

	doProduce := func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		/*
			err := f.region.RoundProduce(&hclient.RoundProduceArgs{}, &hclient.RoundProduceReply{})
			if err != nil {
				flash.Error("Action error: " + err.Error())
			}
		*/
		ctx.Redirect("/game/user")
	}

	doCityStudy := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityStudy) {
		/*
			reply := hclient.CityStudyReply{}
			args := hclient.CityStudyArgs{
				UserId:      ptou(sess.Get("userid")),
				CharacterId: info.CharacterId,
				CityId:      info.CityId,
				KnowledgeId: info.KnowledgeId,
			}
			err := f.region.CityStudy(&args, &reply)
			if err != nil {
				flash.Error("Action error: " + err.Error())
			}
		*/
		ctx.Redirect("/game/land?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}

	doCityBuild := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityBuild) {
		/*
			reply := hclient.CityBuildReply{}
			args := hclient.CityBuildArgs{
				UserId:      ptou(sess.Get("userid")),
				CharacterId: info.CharacterId,
				CityId:      info.CityId,
				BuildingId:  info.BuildingId,
			}
			err := f.region.CityBuild(&args, &reply)
			if err != nil {
				flash.Error("Action error: " + err.Error())
			}
		*/
		ctx.Redirect("/game/land?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}

	doCityTrain := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityTrain) {
		/*
			reply := hclient.CityTrainReply{}
			args := hclient.CityTrainArgs{
				UserId:      ptou(sess.Get("userid")),
				CharacterId: info.CharacterId,
				CityId:      info.CityId,
				UnitId:      info.UnitId,
			}
			err := f.region.CityTrain(&args, &reply)
			if err != nil {
				flash.Error("Action error: " + err.Error())
			} else {
				flash.Info("Started!")
			}
		*/
		ctx.Redirect("/game/land?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}

	doCityCreateArmy := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityArmyCreate) {
		/*
			reply := hclient.CityCreateArmyReply{}
			args := hclient.CityCreateArmyArgs{
				UserId:      ptou(sess.Get("userid")),
				CharacterId: info.CharacterId,
				CityId:      info.CityId,
				Name:        info.Name,
			}
			err := f.region.CityCreateArmy(&args, &reply)
			if err != nil {
				flash.Error("Action error: " + err.Error())
			}
		*/
		ctx.Redirect("/game/land?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}

	doCityTransferUnit := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityUnitTransfer) {
		/*
			reply := hclient.CityTransferUnitReply{}
			args := hclient.CityTransferUnitArgs{
				UserId:      ptou(sess.Get("userid")),
				CharacterId: info.CharacterId,
				CityId:      info.CityId,
				UnitId:      info.UnitId,
				ArmyId:      info.ArmyId,
			}
			err := f.region.CityTransferUnit(&args, &reply)
			if err != nil {
				flash.Error("Action error: " + err.Error())
			}
		*/
		ctx.Redirect("/game/land?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}

	doCityCommandArmy := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityArmyCommand) {
		/*
			reply := hclient.CityCommandArmyReply{}
			args := hclient.CityCommandArmyArgs{
				UserId:      ptou(sess.Get("userid")),
				CharacterId: info.CharacterId,
				CityId:      info.CityId,
				ArmyId:      info.ArmyId,
				Cell:        info.Cell,
				Action:      info.Action,
			}
			err := f.region.CityCommandArmy(&args, &reply)
			if err != nil {
				flash.Error("Action error: " + err.Error())
			}
		*/
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
	m.Post("/action/city/army/command", binding.Bind(FormCityArmyCommand{}), doCityCommandArmy)
	m.Post("/action/city/army/create", binding.Bind(FormCityArmyCreate{}), doCityCreateArmy)
	m.Post("/action/city/unit/transfer", binding.Bind(FormCityUnitTransfer{}), doCityTransferUnit)
}

func (f *FrontService) routeMiddlewares(m *macaron.Macaron) {
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
