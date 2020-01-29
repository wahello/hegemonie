// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync/atomic"
)

func (db *Db) Init() {
	db.UsersById = make([]*User, 0)
	db.ReHash()
}

func (db *Db) Check() error {
	return nil
}

func (db *Db) ReHash() error {
	db.UsersByMail = make(map[string]*User, 0)
	db.NextId = 0
	// FIXME(jfs): Sort the array
	for _, u := range db.UsersById {
		db.UsersByMail[u.Email] = u
		if u.Id > db.NextId {
			db.NextId = u.Id
		}
	}
	db.NextId++
	return nil
}

func (db *Db) UserGet(id uint64) *User {
	for _, u := range db.UsersById {
		if u.Id == id {
			return u
		}
	}
	return nil
}

func (db *Db) UserLookup(mail string) *User {
	if u, ok := db.UsersByMail[mail]; !ok {
		return nil
	} else {
		return u
	}
}

func (db *Db) Create(email string) *User {
	id := atomic.AddUint64(&db.NextId, 1)
	u := &User{
		Id:         id,
		Name:       "NOT-SET",
		Email:      email,
		Password:   "",
		Characters: make([]Character, 0),
	}
	db.UsersById = append(db.UsersById, u)
	db.UsersByMail[u.Email] = u
	return u
}

func (db *Db) SetPassword(u *User, pass string) {
	u.Password = hashPassword(pass, db.Salt)
}

func (db *Db) AuthBasic(u *User, pass string) error {
	if !u.Valid() {
		return errors.New("User suspended")
	}
	if len(u.Password) <= 0 {
		return errors.New("Permission Denied")
	}

	if u.Password[0] == ':' {
		if u.Password[1:] != pass {
			return errors.New("Permission Denied")
		}
	} else {
		if u.Password[:] != hashPassword(pass, db.Salt) {
			return errors.New("Permission Denied")
		}
	}
	return nil
}

func (u *User) Valid() bool {
	return u != nil && !u.Suspended && !u.Deleted
}

func hashPassword(pass, salt string) string {
	checksum := sha256.New()
	checksum.Write([]byte(salt))
	checksum.Write([]byte(pass))
	return hex.EncodeToString(checksum.Sum(nil))
}
