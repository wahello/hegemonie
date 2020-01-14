// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"sort"
	"testing"
)

func TestSetOfUsers(t *testing.T) {
	u2 := &User{Id: 2}
	s := SetOfUsers{}
	s.Add(&User{Id: 1})
	s.Add(&User{Id: 3})
	s.Add(u2)
	s.Add(&User{Id: 4})
	if len(s) != 4 {
		t.Fatal()
	}
	if !sort.IsSorted(s) {
		t.Fatal()
	}
	if u2 != s.Get(2) {
		t.Fatal()
	}
}

func TestUserPasswordNotInClear(t *testing.T) {
	salt := "llkjlkj"
	s := SetOfUsers{}
	s.Create(1, "plop", "pass", salt)
	s.Create(3, "plep", "piss", salt)
	s.Create(2, "plip", "puss", salt)
	s.Create(4, "plup", "poss", salt)
	if !sort.IsSorted(s) {
		t.Fatal()
	}
	if s.Get(1).Password == "pass" {
		t.Fatal()
	}
	if s.Get(3) != s.Lookup("plep") {
		t.Fatal()
	}
	if id, err := s.Auth("plip", "puss", salt); id <= 0 || err != nil {
		t.Fatal()
	}
}
