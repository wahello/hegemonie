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

func doMove(f *frontService, m *macaron.Macaron) macaron.Handler {
	return func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		_, err := f.authenticateAdminFromSession(ctx, sess)
		if err != nil {
			flash.Warning(err.Error())
			ctx.Redirect("/game/user")
			return
		}

		cliReg := region.NewAdminClient(f.cnxRegion)
		_, err = cliReg.Move(contextMacaronToGrpc(ctx, sess), &region.None{})
		if err != nil {
			flash.Warning(err.Error())
		}
		ctx.Redirect("/game/admin")
	}
}

func doProduce(f *frontService, m *macaron.Macaron) macaron.Handler {
	return func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
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
}
