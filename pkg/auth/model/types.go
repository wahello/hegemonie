// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package auth

import "time"

type Character struct {
	ID      uint64 `json:"Id"`
	Region  string
	Name    string
	Off     bool `json:",omitempty"`
	Deleted bool `json:",omitempty"`
}

type User struct {
	ID         uint64 `json:"Id"`
	Name       string
	Email      string
	Password   string
	Characters []*Character
	Admin      bool `json:",omitempty"`
	Inactive   bool `json:",omitempty"`
	Suspended  bool `json:",omitempty"`
	Deleted    bool `json:",omitempty"`

	// Time of the creation of the User
	CTime time.Time
	// Time of the last update of the User
	MTime time.Time
	// Time of the last Inactive update
	ITime time.Time
}

type Db struct {
	UsersByID   []*User
	UsersByMail map[string]*User

	NextID uint64
	Salt   string
}
