// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"bytes"
	"encoding/json"
	"github.com/go-macaron/session"
	auth "github.com/jfsmig/hegemonie/pkg/auth/proto"
	event "github.com/jfsmig/hegemonie/pkg/event/proto"
	region "github.com/jfsmig/hegemonie/pkg/region/proto"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gopkg.in/macaron.v1"
	"strings"
)

type ActionPage func(*macaron.Context, session.Store, *session.Flash)

type NoFlashPage func(*macaron.Context, session.Store)

type StatelessPage func(*macaron.Context)

func serveRoot(ctx *macaron.Context, sess session.Store, flash *session.Flash) {
	ctx.Data["Title"] = "Hegemonie"
	ctx.Data["userid"] = sess.Get("userid")
	ctx.HTML(200, "index")
}

func serveGameUser(f *frontService) ActionPage {
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

func serveGameCharacter(f *frontService) ActionPage {
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

		charLog := make([]string, 0)

		// Load the log of events for the given character
		cliEvt := event.NewConsumerClient(f.cnxEvent)
		log, err := cliEvt.List(contextMacaronToGrpc(ctx, sess), &event.ListReq{CharId: cView.Id, Max: 25})
		if err != nil {
			utils.Logger.Warn().Err(err).Msg("log query")
		} else {
			var sb strings.Builder
			jsonEncoder := json.NewEncoder(&sb)
			localizer := i18n.NewLocalizer(f.translations, "en", "en")

			for _, logEvent := range log.Items {
				buf := bytes.NewBuffer(logEvent.Payload)
				decoder := json.NewDecoder(buf)
				var params map[string]interface{}
				if err = decoder.Decode(&params); err != nil {
					utils.Logger.Warn().Bytes("event", logEvent.Payload).Err(err).Msg("log format")
				} else {
					t, _ := params["action"].(string)
					msg, err := localizer.Localize(&i18n.LocalizeConfig{
						MessageID: t, TemplateData: params, PluralCount: 2,
					})

					// The rendering of the message failed, we dispaly the raw Json
					// instead of the expected rendered
					if err != nil {
						utils.Logger.Warn().Err(err).Msg("localizing-1")
						sb.Reset()
						err = jsonEncoder.Encode(&params)
						if err != nil {
							utils.Logger.Warn().Err(err).Msg("json marshalling")
						} else {
							params["Json"] = sb.String()
							msg, err = localizer.Localize(&i18n.LocalizeConfig{
								DefaultMessage: &i18n.Message{ID: t, One: "{{.Json}}", Other: "{{.Json}}"},
								TemplateData:   params,
							})
							if err != nil {
								utils.Logger.Warn().Err(err).Msg("localizing-2")
							}
						}
					}

					if msg != "" {
						charLog = append(charLog, msg)
					}
				}
			}
		}

		// Query the World server for the Character
		ctx.Data["Title"] = uView.Name + "|" + cView.Name
		ctx.Data["userid"] = utoa(uView.Id)
		ctx.Data["User"] = uView
		ctx.Data["cid"] = utoa(cView.Id)
		ctx.Data["Character"] = cView
		ctx.Data["Cities"] = list.Items
		ctx.Data["Log"] = charLog
		ctx.HTML(200, "character")
	}
}
