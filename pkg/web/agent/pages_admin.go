// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"github.com/go-macaron/session"
	region "github.com/jfsmig/hegemonie/pkg/region/proto"
	"gopkg.in/macaron.v1"
	"sort"
)

func serveGameAdmin(f *FrontService) ActionPage {
	return func(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
		uView, err := f.authenticateAdminFromSession(ctx, sess)
		if err != nil {
			flash.Error(err.Error())
			ctx.Redirect("/")
			return
		}

		cli := region.NewAdminClient(f.cnxRegion)
		scoreBoard, err := cli.GetScores(contextMacaronToGrpc(ctx, sess), &region.None{})
		if err != nil {
			flash.Warning("Region error: " + err.Error())
			ctx.Redirect("/game/admin")
			return
		}
		t := scoreBoard.Items
		sort.Slice(scoreBoard.Items, func(i, j int) bool {
			return t[i].Score > t[j].Score || (t[i].Score == t[j].Score && t[i].Id < t[j].Id)
		})

		ctx.Data["Scores"] = scoreBoard.Items
		ctx.Data["Title"] = uView.Name
		ctx.Data["userid"] = utoa(uView.Id)
		ctx.Data["User"] = uView
		ctx.HTML(200, "admin")
	}
}
