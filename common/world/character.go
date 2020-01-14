// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"errors"
	"sort"
)

func (s SetOfCharacters) Len() int           { return len(s) }
func (s SetOfCharacters) Less(i, j int) bool { return s[i].Id < s[j].Id }
func (s SetOfCharacters) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s SetOfCharacters) Get(id uint64) *Character {
	// TODO(jfs): lookup in the sorted array
	for _, c := range s {
		if c.Id == id {
			return c
		}
	}
	return nil
}

func (s *SetOfCharacters) Add(c *Character) {
	*s = append(*s, c)
	sort.Sort(*s)
}

func (w *World) CharacterGet(id uint64) *Character {
	return w.Auth.Characters.Get(id)
}

// Notify the caller of the cities managed by the given Character.
func (w *World) CharacterGetCities(id uint64, owner func(*City), deputy func(*City)) {
	if id <= 0 {
		return
	}

	w.rw.RLock()
	defer w.rw.RUnlock()

	for _, c := range w.Live.Cities {
		if c.Owner == id {
			owner(c)
		} else if c.Deputy == id {
			deputy(c)
		}
	}
}

func (w *World) CharacterShow(uid, cid uint64) (view CharacterView, err error) {
	if cid <= 0 || uid <= 0 {
		err = errors.New("EINVAL")
	} else {
		w.rw.RLock()
		defer w.rw.RUnlock()

		pChar := w.CharacterGet(cid)
		pUser := w.UserGet(uid)
		if pChar == nil || pUser == nil {
			err = errors.New("Not Found")
		} else if pChar.User != pUser.Id {
			err = errors.New("Forbidden")
		} else {
			view.Id = pChar.Id
			view.Name = pChar.Name
			view.User = NamedItem{Id: pUser.Id, Name: pUser.Name}
			view.DeputyOf = make([]NamedItem, 0)
			view.OwnerOf = make([]NamedItem, 0)
			w.CharacterGetCities(cid,
				func(city *City) {
					view.OwnerOf = append(view.OwnerOf, NamedItem{Id: city.Id, Name: city.Name})
				},
				func(city *City) {
					view.DeputyOf = append(view.DeputyOf, NamedItem{Id: city.Id, Name: city.Name})
				})
		}
	}
	return view, err
}
