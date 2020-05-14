// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"context"
	"errors"
	"github.com/go-macaron/binding"
	"github.com/go-macaron/session"
	auth "github.com/jfsmig/hegemonie/pkg/auth/proto"
	region "github.com/jfsmig/hegemonie/pkg/region/proto"
	"gopkg.in/macaron.v1"
)

func (f *FrontService) authenticateUserFromSession(sess session.Store) (*auth.UserView, error) {
	// Validate the session data
	userid := ptou(sess.Get("userid"))
	if userid == 0 {
		return nil, errors.New("Not authenticated")
	}

	// Authorize the character with the user
	cliAuth := auth.NewAuthClient(f.cnxAuth)
	return cliAuth.UserShow(context.Background(),
		&auth.UserShowReq{Id: userid})
}

func (f *FrontService) authenticateAdminFromSession(sess session.Store) (*auth.UserView, error) {
	uView, err := f.authenticateUserFromSession(sess)
	if err != nil {
		return nil, err
	}
	if !uView.Admin {
		return nil, errors.New("No administration permissions")
	}
	return uView, nil
}

func (f *FrontService) authenticateCharacterFromSession(sess session.Store, idChar uint64) (*auth.UserView, *auth.CharacterView, error) {
	// Validate the session data
	userid := ptou(sess.Get("userid"))
	if userid == 0 || idChar == 0 {
		return nil, nil, errors.New("Not authenticated")
	}

	// Authorize the character with the user
	cliAuth := auth.NewAuthClient(f.cnxAuth)
	uView, err := cliAuth.CharacterShow(context.Background(),
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

		// Authorize the character with the user
		cliAuth := auth.NewAuthClient(f.cnxAuth)
		uView, err := cliAuth.UserAuth(context.Background(),
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
		_, err := f.authenticateAdminFromSession(sess)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/")
			return
		}

		cliReg := region.NewAdminClient(f.cnxRegion)
		_, err = cliReg.Move(context.Background(), &region.None{})
		if err != nil {
			flash.Warning(err.Error())
		}
		ctx.Redirect("/game/admin")
	}

	doProduce := func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		_, err := f.authenticateAdminFromSession(sess)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		cliReg := region.NewAdminClient(f.cnxRegion)
		_, err = cliReg.Produce(context.Background(), &region.None{})
		if err != nil {
			flash.Warning(err.Error())
		}
		ctx.Redirect("/game/admin")
	}

	doCityStudy := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityStudy) {
		_, _, err := f.authenticateCharacterFromSession(sess, info.CharacterId)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		cliReg := region.NewCityClient(f.cnxRegion)
		_, err = cliReg.Study(context.Background(),
			&region.StudyReq{City: info.CityId, Character: info.CharacterId, KnowledgeType: info.KnowledgeId})
		if err != nil {
			flash.Warning(err.Error())
		}

		ctx.Redirect("/game/land/knowledges?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}

	doCityBuild := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityBuild) {
		_, _, err := f.authenticateCharacterFromSession(sess, info.CharacterId)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		cliReg := region.NewCityClient(f.cnxRegion)
		_, err = cliReg.Build(context.Background(),
			&region.BuildReq{City: info.CityId, Character: info.CharacterId, BuildingType: info.BuildingId})
		if err != nil {
			flash.Warning(err.Error())
		}

		ctx.Redirect("/game/land/buildings?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
	}

	doCityTrain := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityTrain) {
		_, _, err := f.authenticateCharacterFromSession(sess, info.CharacterId)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		cliReg := region.NewCityClient(f.cnxRegion)
		_, err = cliReg.Train(context.Background(),
			&region.TrainReq{City: info.CityId, Character: info.CharacterId, UnitType: info.UnitId})
		if err != nil {
			flash.Warning(err.Error())
		}

		ctx.Redirect("/game/land/units?cid=" + utoa(info.CharacterId) + "&lid=" + utoa(info.CityId))
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

	doCityDisbandArmy := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityArmyDisband) {
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

	doCityCancelArmy := func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormCityArmyCancel) {
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

	Cell   uint64 `form:"cell" binding:"Required"`
	Action uint64 `form:"what" binding:"Required"`
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
