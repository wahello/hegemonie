// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"github.com/go-macaron/session"
	auth "github.com/jfsmig/hegemonie/pkg/auth/proto"
	region "github.com/jfsmig/hegemonie/pkg/region/proto"
	"gopkg.in/macaron.v1"
)

type ActionPage func(*macaron.Context, session.Store, *session.Flash)

type NoFlashPage func(*macaron.Context, session.Store)

type StatelessPage func(*macaron.Context)

func serveRoot(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	ctx.Data["Title"] = "Hegemonie"
	ctx.Data["userid"] = sess.Get("userid")
	ctx.HTML(200, "index")
}

func serveGameUser(f *FrontService) ActionPage {
	return func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		uView, err := f.authenticateUserFromSession(ctx, sess)
		if err != nil {
			flash.Error(err.Error())
			ctx.Redirect("/")
			return
		}

		cliReg := region.NewCityClient(f.cnxRegion)
		for _, c := range uView.Characters {
			l, err := cliReg.List(contextMacaronToGrpc(ctx, sess), &region.ListReq{Character: c.Id})
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
		uView, cView, err := f.authenticateCharacterFromSession(ctx, sess, atou(ctx.Query("cid")))
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		// Load the Cities managed by the current Character
		cliReg := region.NewCityClient(f.cnxRegion)
		list, err := cliReg.List(contextMacaronToGrpc(ctx, sess), &region.ListReq{Character: cView.Id})
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
