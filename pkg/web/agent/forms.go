// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"context"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
	"github.com/jfsmig/hegemonie/pkg/auth/proto"
	"gopkg.in/macaron.v1"
)

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
		ctx.Redirect("/game/land/overview?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
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
		ctx.Redirect("/game/land/overview?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
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
		ctx.Redirect("/game/land/overview?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
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
		ctx.Redirect("/game/land/overview?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
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
		ctx.Redirect("/game/land/overview?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
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
		ctx.Redirect("/game/land/overview?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
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
