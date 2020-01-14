// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import "github.com/juju/errors"

func (s SetOfId) Len() int           { return len(s) }
func (s SetOfId) Less(i, j int) bool { return s[i] < s[j] }
func (s SetOfId) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s SetOfFights) Len() int           { return len(s) }
func (s SetOfFights) Less(i, j int) bool { return s[i].Id < s[j].Id }
func (s SetOfFights) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// Leave the Fight as a loser
func (f *Fight) Flea(w *World, a *Army) error {
	return errors.New("NYI")
}

// Change the side in the Fight.
// If the Army was defending, it becomes an attacker, if it was an attacker
// it becomes a defender.
func (f *Fight) Flip(w *World, a *Army) error {
	return errors.New("NYI")
}
