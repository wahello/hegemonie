// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_web_agent

import (
	"errors"
	"github.com/go-macaron/session"
	"github.com/google/uuid"
	auth "github.com/jfsmig/hegemonie/pkg/auth/proto"
	"gopkg.in/macaron.v1"
)

func (f *frontService) authenticateUserFromSession(ctx *macaron.Context, sess session.Store) (*auth.UserView, error) {
	// Validate the session data
	userID := ptou(sess.Get("userid"))
	if userID == 0 {
		return nil, errors.New("Not authenticated")
	}

	// Authorize the character with the user
	cliAuth := auth.NewAuthClient(f.cnxAuth)
	return cliAuth.UserShow(contextMacaronToGrpc(ctx, sess),
		&auth.UserShowReq{Id: userID})
}

func (f *frontService) authenticateAdminFromSession(ctx *macaron.Context, sess session.Store) (*auth.UserView, error) {
	uView, err := f.authenticateUserFromSession(ctx, sess)
	if err != nil {
		return nil, err
	}
	if !uView.Admin {
		return nil, errors.New("No administration permissions")
	}
	return uView, nil
}

func (f *frontService) authenticateCharacterFromSession(ctx *macaron.Context, sess session.Store, idChar uint64) (*auth.UserView, *auth.CharacterView, error) {
	// Validate the session data
	userID := ptou(sess.Get("userid"))
	if userID == 0 || idChar == 0 {
		return nil, nil, errors.New("Not authenticated")
	}

	// Authorize the character with the user
	cliAuth := auth.NewAuthClient(f.cnxAuth)
	uView, err := cliAuth.CharacterShow(contextMacaronToGrpc(ctx, sess),
		&auth.CharacterShowReq{User: userID, Character: idChar})
	if err != nil {
		return nil, nil, err
	}

	return uView, uView.Characters[0], nil
}

type FormLogin struct {
	UserMail string `form:"email" binding:"Required"`
	UserPass string `form:"password" binding:"Required"`
}

func doLogin(f *frontService, m *macaron.Macaron) macaron.Handler {
	return func(ctx *macaron.Context, flash *session.Flash, sess session.Store, info FormLogin) {
		// Cleanup a previous session
		sess.Flush()

		sessionID := uuid.New().String()
		sess.Set("session-id", sessionID)

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
}

func doLogout(f *frontService, m *macaron.Macaron) macaron.Handler {
	return func(ctx *macaron.Context, s session.Store) {
		ctx.SetSecureCookie("session", "")
		s.Flush()
		ctx.Redirect("/")
	}
}
