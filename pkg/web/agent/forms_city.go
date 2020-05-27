// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"github.com/go-macaron/session"
	region "github.com/jfsmig/hegemonie/pkg/region/proto"
	"gopkg.in/macaron.v1"
	"strings"
)

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

	Location uint64 `form:"location" binding:"Required"`
	Action   string `form:"action" binding:"Required"`
}

type FormCityArmyDisband struct {
	CharacterId uint64 `form:"cid" binding:"Required"`
	CityId      uint64 `form:"lid" binding:"Required"`
	ArmyId      uint64 `form:"aid" binding:"Required"`
}

type FormCityArmyCancel struct {
	CharacterId uint64 `form:"cid" binding:"Required"`
	CityId      uint64 `form:"lid" binding:"Required"`
	ArmyId      uint64 `form:"aid" binding:"Required"`
}

func doCityStudy(f *FrontService, m *macaron.Macaron) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityStudy) {
		_, _, err := f.authenticateCharacterFromSession(ctx, sess, info.CharacterId)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		cliReg := region.NewCityClient(f.cnxRegion)
		_, err = cliReg.Study(contextMacaronToGrpc(ctx, sess),
			&region.StudyReq{City: info.CityId, Character: info.CharacterId, KnowledgeType: info.KnowledgeId})
		if err != nil {
			flash.Warning(err.Error())
		}

		ctx.Redirect("/game/land/knowledges?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}
}

func doCityBuild(f *FrontService, m *macaron.Macaron) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityBuild) {
		_, _, err := f.authenticateCharacterFromSession(ctx, sess, info.CharacterId)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		cliReg := region.NewCityClient(f.cnxRegion)
		_, err = cliReg.Build(contextMacaronToGrpc(ctx, sess),
			&region.BuildReq{City: info.CityId, Character: info.CharacterId, BuildingType: info.BuildingId})
		if err != nil {
			flash.Warning(err.Error())
		}

		ctx.Redirect("/game/land/buildings?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}
}

func doCityTrain(f *FrontService, m *macaron.Macaron) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityTrain) {
		_, _, err := f.authenticateCharacterFromSession(ctx, sess, info.CharacterId)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		cliReg := region.NewCityClient(f.cnxRegion)
		_, err = cliReg.Train(contextMacaronToGrpc(ctx, sess),
			&region.TrainReq{City: info.CityId, Character: info.CharacterId, UnitType: info.UnitId})
		if err != nil {
			flash.Warning(err.Error())
		}

		ctx.Redirect("/game/land/units?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}
}

func doCityCreateArmy(f *FrontService, m *macaron.Macaron) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityArmyCreate) {
		ctx.Redirect("/game/land/overview?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}
}

func doCityTransferUnit(f *FrontService, m *macaron.Macaron) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityUnitTransfer) {
		ctx.Redirect("/game/land/overview?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}
}

func doCityDisbandArmy(f *FrontService, m *macaron.Macaron) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityArmyDisband) {
		ctx.Redirect("/game/land/overview?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}
}

func doCityCancelArmy(f *FrontService, m *macaron.Macaron) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityArmyCancel) {
		ctx.Redirect("/game/land/overview?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}
}

func doCityCommandArmy(f *FrontService, m *macaron.Macaron) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityArmyCommand) {
		_, _, err := f.authenticateCharacterFromSession(ctx, sess, info.CharacterId)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		url := "/game/army?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId) + "&aid=" + utoa(info.ArmyId)

		actionId := region.ArmyCommandType_Move
		switch strings.ToLower(info.Action) {
		case "move":
			actionId = region.ArmyCommandType_Move
		case "wait":
			actionId = region.ArmyCommandType_Wait
		case "attack":
			actionId = region.ArmyCommandType_Attack
		case "defend":
			actionId = region.ArmyCommandType_Defend
		case "overlord":
			actionId = region.ArmyCommandType_Overlord
		case "break":
			actionId = region.ArmyCommandType_Break
		case "massacre":
			actionId = region.ArmyCommandType_Massacre
		case "deposit":
			actionId = region.ArmyCommandType_Deposit
		case "disband":
			actionId = region.ArmyCommandType_Disband
		default:
			flash.Warning("Invalid action name")
			ctx.Redirect(url)
			return
		}

		cli := region.NewArmyClient(f.cnxRegion)
		cmd := &region.ArmyCommandReq{
			Id: &region.ArmyId{
				Character: info.CharacterId,
				City:      info.CityId,
				Army:      info.ArmyId,
			},
			Command: &region.ArmyCommand{
				Action: actionId,
				Target: info.Location,
			},
		}
		_, err = cli.Command(contextMacaronToGrpc(ctx, sess), cmd)
		if err != nil {
			flash.Warning(err.Error())
		}

		ctx.Redirect(url)
	}
}
