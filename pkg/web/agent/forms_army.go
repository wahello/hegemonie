// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"github.com/go-macaron/session"
	region "github.com/jfsmig/hegemonie/pkg/region/proto"
	"gopkg.in/macaron.v1"
)

type FormArmyId struct {
	CharacterID uint64 `form:"cid" binding:"Required"`
	CityID      uint64 `form:"lid" binding:"Required"`
	ArmyID      uint64 `form:"aid" binding:"Required"`
}

type FormArmyTarget struct {
	CharacterID uint64 `form:"cid" binding:"Required"`
	CityID      uint64 `form:"lid" binding:"Required"`
	ArmyID      uint64 `form:"aid" binding:"Required"`
	TargetID    uint64 `form:"location" binding:"Required"`
}

func doArmyCancel(f *frontService) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormArmyId) {
		client, err := f.authAndConnect(ctx, sess, info.CharacterID)
		if err == nil {
			_, err = client.Cancel(contextMacaronToGrpc(ctx, sess),
				&region.ArmyId{Character: info.CharacterID, City: info.CityID, Army: info.ArmyID})
		}
		if err != nil {
			flash.Warning(err.Error())
		}
		ctx.Redirect(info.Url("/game/land/overview"))
	}
}

func doArmyFlip(f *frontService) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormArmyId) {
		client, err := f.authAndConnect(ctx, sess, info.CharacterID)
		if err == nil {
			_, err = client.Flip(contextMacaronToGrpc(ctx, sess),
				&region.ArmyId{Character: info.CharacterID, City: info.CityID, Army: info.ArmyID})
		}
		if err != nil {
			flash.Warning(err.Error())
		}
		ctx.Redirect(info.Url("/game/army"))
	}
}

func doArmyFlea(f *frontService) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormArmyId) {
		client, err := f.authAndConnect(ctx, sess, info.CharacterID)
		if err == nil {
			_, err = client.Flea(contextMacaronToGrpc(ctx, sess),
				&region.ArmyId{Character: info.CharacterID, City: info.CityID, Army: info.ArmyID})
		}
		if err != nil {
			flash.Warning(err.Error())
		}
		ctx.Redirect(info.Url("/game/army"))
	}
}

func doArmyDisband(f *frontService) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormArmyTarget) {
		client, err := f.authAndConnect(ctx, sess, info.CharacterID)
		if err == nil {
			_, err = client.Disband(contextMacaronToGrpc(ctx, sess), &region.ArmyTarget{
				Id:     &region.ArmyId{Character: info.CharacterID, City: info.CityID, Army: info.ArmyID},
				Target: info.TargetID,
			})
		}
		if err != nil {
			flash.Warning(err.Error())
		}
		ctx.Redirect(info.Url("/game/army"))
	}
}

func doArmyMove(f *frontService) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormArmyTarget) {
		client, err := f.authAndConnect(ctx, sess, info.CharacterID)
		if err == nil {
			_, err = client.Move(contextMacaronToGrpc(ctx, sess), &region.ArmyMoveReq{
				Id:     &region.ArmyId{Character: info.CharacterID, City: info.CityID, Army: info.ArmyID},
				Target: info.TargetID,
				Args:   &region.ArmyMoveArgs{},
			})
		}
		if err != nil {
			flash.Warning(err.Error())
		}
		ctx.Redirect(info.Url("/game/army"))
	}
}

func doArmyWait(f *frontService) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormArmyTarget) {
		client, err := f.authAndConnect(ctx, sess, info.CharacterID)
		if err == nil {
			_, err = client.Wait(contextMacaronToGrpc(ctx, sess), &region.ArmyTarget{
				Id:     &region.ArmyId{Character: info.CharacterID, City: info.CityID, Army: info.ArmyID},
				Target: info.TargetID,
			})
		}
		if err != nil {
			flash.Warning(err.Error())
		}
		ctx.Redirect(info.Url("/game/army"))
	}
}

func doArmyDefend(f *frontService) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormArmyTarget) {
		client, err := f.authAndConnect(ctx, sess, info.CharacterID)
		if err == nil {
			_, err = client.Defend(contextMacaronToGrpc(ctx, sess), &region.ArmyTarget{
				Id:     &region.ArmyId{Character: info.CharacterID, City: info.CityID, Army: info.ArmyID},
				Target: info.TargetID,
			})
		}
		if err != nil {
			flash.Warning(err.Error())
		}
		ctx.Redirect(info.Url("/game/army"))
	}
}

func doArmyAssault(f *frontService) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormArmyTarget) {
		client, err := f.authAndConnect(ctx, sess, info.CharacterID)
		if err == nil {
			_, err = client.Attack(contextMacaronToGrpc(ctx, sess), &region.ArmyAssaultReq{
				Id:     &region.ArmyId{Character: info.CharacterID, City: info.CityID, Army: info.ArmyID},
				Target: info.TargetID,
				Args:   &region.ArmyAssaultArgs{},
			})
		}
		if err != nil {
			flash.Warning(err.Error())
		}
		ctx.Redirect(info.Url("/game/army"))
	}
}

func (f FormArmyTarget) Url(page string) string {
	return page + "?cid=" + utoa(f.CharacterID) + "&lid=" + utoa(f.CityID) + "&aid=" + utoa(f.ArmyID)
}

func (f FormArmyId) Url(page string) string {
	return page + "?cid=" + utoa(f.CharacterID) + "&lid=" + utoa(f.CityID) + "&aid=" + utoa(f.ArmyID)
}

func (f *frontService) authAndConnect(ctx *macaron.Context, sess session.Store, cid uint64) (region.ArmyClient, error) {
	_, _, err := f.authenticateCharacterFromSession(ctx, sess, cid)
	if err != nil {
		return nil, err
	}
	return region.NewArmyClient(f.cnxRegion), nil
}
