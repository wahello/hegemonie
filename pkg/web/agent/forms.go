// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"errors"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
	"github.com/google/uuid"
	auth "github.com/jfsmig/hegemonie/pkg/auth/proto"
	region "github.com/jfsmig/hegemonie/pkg/region/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"gopkg.in/macaron.v1"
	"strings"
)

func (f *FrontService) authenticateUserFromSession(ctx *macaron.Context, sess session.Store) (*auth.UserView, error) {
	// Validate the session data
	userid := ptou(sess.Get("userid"))
	if userid == 0 {
		return nil, errors.New("Not authenticated")
	}

	// Authorize the character with the user
	cliAuth := auth.NewAuthClient(f.cnxAuth)
	return cliAuth.UserShow(contextMacaronToGrpc(ctx, sess),
		&auth.UserShowReq{Id: userid})
}

func (f *FrontService) authenticateAdminFromSession(ctx *macaron.Context, sess session.Store) (*auth.UserView, error) {
	uView, err := f.authenticateUserFromSession(ctx, sess)
	if err != nil {
		return nil, err
	}
	if !uView.Admin {
		return nil, errors.New("No administration permissions")
	}
	return uView, nil
}

func (f *FrontService) authenticateCharacterFromSession(ctx *macaron.Context, sess session.Store, idChar uint64) (*auth.UserView, *auth.CharacterView, error) {
	// Validate the session data
	userid := ptou(sess.Get("userid"))
	if userid == 0 || idChar == 0 {
		return nil, nil, errors.New("Not authenticated")
	}

	// Authorize the character with the user
	cliAuth := auth.NewAuthClient(f.cnxAuth)
	uView, err := cliAuth.CharacterShow(contextMacaronToGrpc(ctx, sess),
		&auth.CharacterShowReq{User: userid, Character: idChar})
	if err != nil {
		return nil, nil, err
	}

	return uView, uView.Characters[0], nil
}

func (f *FrontService) routeForms(m *macaron.Macaron) {
	doLogIn := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormLogin) {
		// Cleanup a previous session
		sess.Flush()

		sessionId := uuid.New().String()
		sess.Set("session-id", sessionId)

		// Authorize the character with the user
		cliAuth := auth.NewAuthClient(f.cnxAuth)
		uView, err := cliAuth.UserAuth(contextMacaronToGrpc(ctx, sess),
			&auth.UserAuthReq{Mail: info.UserMail, Pass: info.UserPass})

		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/")
		} else {
			strId := utoa(uView.Id)
			ctx.SetSecureCookie("session", strId)
			sess.Set("userid", strId)
			ctx.Redirect("/game/user")
		}
	}

	doLogOut := func(ctx *macaron.Context, s session.Store) {
		ctx.SetSecureCookie("session", "")
		s.Flush()
		ctx.Redirect("/")
	}

	doMove := func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		_, err := f.authenticateAdminFromSession(ctx, sess)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/")
			return
		}

		cliReg := region.NewAdminClient(f.cnxRegion)
		_, err = cliReg.Move(contextMacaronToGrpc(ctx, sess), &region.None{})
		if err != nil {
			flash.Warning(err.Error())
		}
		ctx.Redirect("/game/admin")
	}

	doProduce := func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		_, err := f.authenticateAdminFromSession(ctx, sess)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		cliReg := region.NewAdminClient(f.cnxRegion)
		_, err = cliReg.Produce(contextMacaronToGrpc(ctx, sess), &region.None{})
		if err != nil {
			flash.Warning(err.Error())
		}
		ctx.Redirect("/game/admin")
	}

	doCityStudy := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityStudy) {
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

	doCityBuild := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityBuild) {
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

	doCityTrain := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityTrain) {
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

	doCityCreateArmy := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityArmyCreate) {
		ctx.Redirect("/game/land/overview?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}

	doCityTransferUnit := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityUnitTransfer) {
		ctx.Redirect("/game/land/overview?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}

	doCityCommandArmy := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityArmyCommand) {
		url := "/game/army?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId) + "&aid=" + utoa(info.ArmyId)

		_, _, err := f.authenticateCharacterFromSession(ctx, sess, info.CharacterId)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect(url)
			return
		}

		actionId := region.ArmyCommandType_Move
		switch strings.ToLower(info.Action) {
		case "move":
			actionId = region.ArmyCommandType_Move
		case "attack":
			actionId = region.ArmyCommandType_Attack
		case "defend":
			actionId = region.ArmyCommandType_Defend
		case "wait":
			actionId = region.ArmyCommandType_Wait
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
		utils.Logger.Warn().Interface("req", cmd).Send()
		_, err = cli.Command(contextMacaronToGrpc(ctx, sess), cmd)
		if err != nil {
			flash.Warning(err.Error())
		}

		ctx.Redirect(url)
	}

	doCityDisbandArmy := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityArmyDisband) {
		ctx.Redirect("/game/land/overview?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}

	doCityCancelArmy := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityArmyCancel) {
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
	m.Post("/action/army/cancel", binding.Bind(FormCityArmyDisband{}), doCityCancelArmy)
	m.Post("/action/army/disband", binding.Bind(FormCityArmyDisband{}), doCityDisbandArmy)
	m.Post("/action/army/command", binding.Bind(FormCityArmyCommand{}), doCityCommandArmy)
	m.Post("/action/army/create", binding.Bind(FormCityArmyCreate{}), doCityCreateArmy)
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
