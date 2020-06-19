// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_auth_agent

import (
	"context"
	"github.com/jfsmig/hegemonie/pkg/auth/model"
	proto "github.com/jfsmig/hegemonie/pkg/auth/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func userView(u *auth.User) *proto.UserView {
	return &proto.UserView{
		Id: u.ID, Name: u.Name, Mail: u.Email,
		Admin: u.Admin, Inactive: u.Inactive, Suspended: u.Suspended,
	}
}

func userViewFull(u *auth.User) *proto.UserView {
	rep := userView(u)
	for _, c := range u.Characters {
		if c.Off || c.Deleted {
			continue
		}
		rep.Characters = append(rep.Characters,
			&proto.CharacterView{Id: c.ID, Name: c.Name, Region: c.Region})
	}
	return rep
}

func (srv *authService) UserShow(ctx context.Context, req *proto.UserShowReq) (*proto.UserView, error) {
	var u *auth.User
	if req.Id > 0 {
		u = srv.db.UserGet(req.Id)
	} else if len(req.Mail) > 0 {
		u = srv.db.UserLookup(req.Mail)
	} else {
		return nil, status.Error(codes.InvalidArgument, "Missing ID and Name")
	}

	if u == nil || u.Deleted {
		return nil, status.Error(codes.NotFound, "No such User")
	}
	return userViewFull(u), nil
}

func (srv *authService) UserList(ctx context.Context, req *proto.UserListReq) (*proto.UserListRep, error) {
	if req.Limit <= 0 {
		req.Limit = 1024
	}

	rep := proto.UserListRep{}
	for _, u := range srv.db.UsersByID {
		if uint64(len(rep.Items)) >= req.Limit {
			break
		}
		if u.Deleted || u.ID <= req.Marker {
			continue
		}

		rep.Items = append(rep.Items, userView(u))
	}
	return &rep, nil
}

func (srv *authService) UserCreate(ctx context.Context, req *proto.UserCreateReq) (*proto.UserView, error) {
	var err error
	u := srv.db.UserLookup(req.Mail)
	if u != nil {
		return nil, status.Error(codes.AlreadyExists, "User already registered")
	} else {
		u, err = srv.db.CreateUser(req.Mail)
	}
	return userViewFull(u), err
}

func (srv *authService) UserUpdate(ctx context.Context, req *proto.UserUpdateReq) (*proto.None, error) {
	u := srv.db.UserGet(req.Id)
	if u != nil || u.Deleted {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	if req.Pass != "" {
		srv.db.SetPassword(u, req.Pass)
		u.Name = req.Name
	}
	return &proto.None{}, nil
}

func (srv *authService) UserSuspend(ctx context.Context, req *proto.UserSuspendReq) (*proto.None, error) {
	u := srv.db.UserGet(req.Id)
	if u != nil || u.Deleted {
		return nil, status.Error(codes.NotFound, "User not found")
	}
	u.Suspended = true
	return &proto.None{}, nil
}

func (srv *authService) UserAuth(ctx context.Context, req *proto.UserAuthReq) (*proto.UserView, error) {
	u := srv.db.UserLookup(req.Mail)
	if u == nil {
		return nil, status.Error(codes.NotFound, "No such User")
	}
	err := srv.db.AuthBasic(u, req.Pass)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, "User suspended")
	}
	return userView(u), nil
}

func (srv *authService) CharacterShow(ctx context.Context, req *proto.CharacterShowReq) (*proto.UserView, error) {
	if req.User <= 0 || req.Character <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid user/role ID")
	}

	if u := srv.db.UserGet(req.User); u == nil {
		return nil, status.Error(codes.NotFound, "Not such User")
	} else {
		if u.Suspended || u.Deleted {
			return nil, status.Error(codes.PermissionDenied, "User suspended")
		}
		for _, c := range u.Characters {
			if c.ID == req.Character {
				uView := userView(u)
				uView.Characters = append(uView.Characters, &proto.CharacterView{
					Id: c.ID, Region: c.Region, Name: c.Name, Off: c.Off,
				})
				return uView, nil
			}
		}
		return nil, status.Error(codes.PermissionDenied, "Character mismatch")
	}
}

func (srv *authService) CharacterUpdate(ctx context.Context, req *proto.CharacterUpdateReq) (*proto.None, error) {
	if req.User <= 0 || req.Character <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	u := srv.db.UserGet(req.User)
	if u == nil {
		return nil, status.Error(codes.NotFound, "Not such User")
	}

	c := u.GetCharacter(req.Character)
	if c == nil {
		return nil, status.Error(codes.NotFound, "No such Role")
	}
	c.Name = req.Name
	c.Off = req.Off
	c.Deleted = req.Deleted

	return &proto.None{}, nil
}

func (srv *authService) CharacterMigrate(ctx context.Context, req *proto.CharacterMigrationReq) (*proto.None, error) {
	if req.User <= 0 || req.Character <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Invalid user ID")
	}

	u := srv.db.UserGet(req.User)
	if u == nil {
		return nil, status.Error(codes.NotFound, "Not such User")
	}

	c := u.GetCharacter(req.Character)
	if c == nil {
		return nil, status.Error(codes.NotFound, "No such Role")
	}
	c.Region = req.Region

	return &proto.None{}, nil
}
