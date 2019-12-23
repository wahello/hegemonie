// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

func (w *World) UserCreate(mail, pass string) (uint64, error) {
	if !validMail(mail) || !validPass(pass) {
		return 0, errors.New("EINVAL")
	}

	w.rw.Lock()
	defer w.rw.Unlock()

	h := w.hashPassword(pass)
	for _, u := range w.Users {
		if u.Email == mail && u.Password == h {
			return 0, errors.New("User exists")
		}
	}

	id := w.getNextId()
	u := User{Id: id, Name: "No-Name", Email: mail, Password: pass}
	w.Users = append(w.Users, u)
	return id, nil
}

func (w *World) UserGet(id uint64) *User {
	// TODO(jfs): lookup in the sorted array
	for _, u := range w.Users {
		if u.Id == id {
			return &u
		}
	}

	return nil
}

func (w *World) UserAuth(mail, pass string) (uint64, error) {
	if mail == "" || pass == "" {
		return 0, errors.New("EINVAL")
	}

	w.rw.RLock()
	defer w.rw.RUnlock()

	h := w.hashPassword(pass)
	for _, u := range w.Users {
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

func (w *World) UserGetCharacters(id uint64, hook func(*Character)) {
	for _, c := range w.Characters {
		if c.User == id {
			hook(&c)
		}
	}
}

func (w *World) hashPassword(pass string) string {
	checksum := sha256.New()
	checksum.Write([]byte(w.Salt))
	checksum.Write([]byte(pass))
	return hex.EncodeToString(checksum.Sum(nil))
}

func validMail(m string) bool {
	// TODO(jfs): Not yet implemented
	return len(m) > 0
}

func validPass(p string) bool {
	// TODO(jfs): Not yet implemented
	return len(p) > 0
}

func (s *SetOfUsers) Len() int {
	return len(*s)
}

func (s *SetOfUsers) Less(i, j int) bool {
	return (*s)[i].Id < (*s)[j].Id
}

func (s *SetOfUsers) Swap(i, j int) {
	tmp := (*s)[i]
	(*s)[i] = (*s)[j]
	(*s)[j] = tmp
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
