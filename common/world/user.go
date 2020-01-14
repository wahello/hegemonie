// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sort"
)

func (s SetOfUsers) Len() int           { return len(s) }
func (s SetOfUsers) Less(i, j int) bool { return s[i].Id < s[j].Id }
func (s SetOfUsers) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s SetOfUsers) Get(id uint64) *User {
	// FIXME(jfs): lookup in the sorted array
	for _, u := range s {
		if u.Id == id {
			return u
		}
	}
	return nil
}

func (s SetOfUsers) Lookup(mail string) *User {
	for _, u := range s {
		if u.Email == mail {
			return u
		}
	}
	return nil
}

func (s SetOfUsers) Auth(mail, pass, salt string) (uint64, error) {
	h := hashPassword(pass, salt)
	for _, u := range s {
		if u.Email == mail {
			if u.Password == h {
				// Hashed password matches
				return u.Id, nil
			} else if u.Password[0] == ':' && u.Password[1:] == pass {
				// Clear password matches
				return u.Id, nil
			} else {
				return 0, nil
			}
		}
	}

	return 0, errors.New("User not found")
}

func (s *SetOfUsers) Add(user *User) {
	*s = append(*s, user)
	sort.Sort(s)
}

func (s *SetOfUsers) Create(id uint64, mail, pass, salt string) *User {
	u := &User{Id: id, Name: "No-Name", Email: mail, Password: hashPassword(pass, salt)}
	s.Add(u)
	return u
}

func hashPassword(pass, salt string) string {
	checksum := sha256.New()
	checksum.Write([]byte(salt))
	checksum.Write([]byte(pass))
	return hex.EncodeToString(checksum.Sum(nil))
}

func validMail(m string) bool {
	// FIXME(jfs): Not yet implemented
	return len(m) > 0
}

func validPass(p string) bool {
	// FIXME(jfs): Not yet implemented
	return len(p) > 0
}

func (w *World) UserCreate(mail, pass string) (uint64, error) {
	if !validMail(mail) || !validPass(pass) {
		return 0, errors.New("EINVAL")
	}

	w.rw.Lock()
	defer w.rw.Unlock()

	if nil != w.Auth.Users.Lookup(mail) {
		return 0, errors.New("User exists")
	}

	id := w.getNextId()
	w.Auth.Users.Create(id, mail, pass, w.Salt)
	return id, nil
}

func (w *World) UserGet(id uint64) *User {
	if id <= 0 {
		return nil
	}
	return w.Auth.Users.Get(id)
}

func (w *World) UserAuth(mail, pass string) (uint64, error) {
	if mail == "" || pass == "" {
		return 0, errors.New("EINVAL")
	}

	w.rw.RLock()
	defer w.rw.RUnlock()
	return w.Auth.Users.Auth(mail, pass, w.Salt)
}

func (w *World) UserGetCharacters(id uint64, hook func(*Character)) {
	// TODO(jfs): Maybe an index would be good, but this is not called often in a session.
	for _, c := range w.Auth.Characters {
		if c.User == id {
			hook(c)
		}
	}
}

func (w *World) hashPassword(pass string) string {
	return hashPassword(pass, w.Salt)
}

func (w *World) UserShow(id uint64) (view UserView, err error) {
	if id <= 0 {
		err = errors.New("EINVAL")
	} else {
		w.rw.RLock()
		defer w.rw.RUnlock()

		if u := w.UserGet(id); u == nil {
			err = errors.New("User not found")
		} else {
			view.Id = u.Id
			view.Name = u.Name
			view.Email = u.Email
			view.Inactive = u.Inactive
			view.Admin = u.Admin
			view.Characters = make([]NamedItem, 0)
			w.UserGetCharacters(id, func(c *Character) {
				view.Characters = append(view.Characters, NamedItem{Id: c.Id, Name: c.Name})
			})
		}
	}
	return view, err
}
