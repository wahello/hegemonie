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
	db.UsersByID = make([]*User, 0)
	db.ReHash()
}

func (db *Db) Check() error {
	return nil
}

func (db *Db) ReHash() error {
	db.UsersByMail = make(map[string]*User, 0)
	db.NextID = 0
	// FIXME(jfs): Sort the array
	for _, u := range db.UsersByID {
		db.UsersByMail[u.Email] = u
		if u.ID > db.NextID {
			db.NextID = u.ID
		}
	}
	db.NextID++
	return nil
}

func (db *Db) UserGet(id uint64) *User {
	for _, u := range db.UsersByID {
		if u.ID == id {
			return u
		}
	}
	return nil
}

func (db *Db) UserLookup(mail string) *User {
	if u, ok := db.UsersByMail[mail]; ok {
		return u
	}
	return nil
}

func (db *Db) CreateUser(email string) (*User, error) {
	// FIXME(jfs): Verify the format of the email
	id := atomic.AddUint64(&db.NextID, 1)
	u := &User{
		ID:         id,
		Name:       "NOT-SET",
		Email:      email,
		Password:   "",
		Characters: make([]*Character, 0),
	}
	db.UsersByID = append(db.UsersByID, u)
	db.UsersByMail[u.Email] = u
	return u, nil
}

func (db *Db) CreateCharacter(idUser uint64, name, region string) (*Character, error) {
	idChar := atomic.AddUint64(&db.NextID, 1)
	u := db.UserGet(idUser)
	if u == nil {
		return nil, errors.New("No such User")
	}
	return u.CreateCharacter(idChar, name, region)
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

func (u *User) Rename(n string) *User {
	// FIXME(jfs): validate the name
	u.Name = n
	return u
}

func (u *User) Promote() *User {
	u.Admin = true
	return u
}

func (u *User) Demote() *User {
	u.Admin = false
	return u
}

func (u *User) SetRawPassword(p string) *User {
	// FIXME(jfs): validate the password
	u.Password = p
	return u
}

func (u *User) CreateCharacter(id uint64, name, region string) (*Character, error) {
	// FIXME(jfs): Verify the format of the name
	c := &Character{ID: id, Name: name, Region: region}
	u.Characters = append(u.Characters, c)
	return c, nil
}

func (u *User) GetCharacter(idChar uint64) *Character {
	for _, c := range u.Characters {
		if c.ID == idChar {
			return c
		}
	}
	return nil
}

func hashPassword(pass, salt string) string {
	checksum := sha256.New()
	checksum.Write([]byte(salt))
	checksum.Write([]byte(pass))
	return hex.EncodeToString(checksum.Sum(nil))
}
