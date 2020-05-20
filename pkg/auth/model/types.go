// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package auth

import "time"

type Character struct {
	Id      uint64
	Region  string
	Name    string
	Off     bool `json:",omitempty"`
	Deleted bool `json:",omitempty"`
}

type User struct {
	Id         uint64
	Name       string
	Email      string
	Password   string
	Characters []Character
	Admin      bool `json:",omitempty"`
	Inactive   bool `json:",omitempty"`
	Suspended  bool `json:",omitempty"`
	Deleted    bool `json:",omitempty"`

	CTime time.Time
	MTime time.Time
	ATime time.Time
}

type Db struct {
	UsersById   []*User
	UsersByMail map[string]*User

	NextId uint64
	Salt   string
}
