// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import "sort"

func (s SetOfArmies) Len() int           { return len(s) }
func (s SetOfArmies) Less(i, j int) bool { return s[i].Id < s[j].Id }
func (s SetOfArmies) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s SetOfArmies) Get(id uint64) *Army {
	for _, a := range s {
		if a.Id == id {
			return a
		}
	}
	return nil
}

func (s *SetOfArmies) Add(a *Army) {
	*s = append(*s, a)
	sort.Sort(s)
}

func (w *World) ArmyCreate(c *City, name string) (uint64, error) {
	a := &Army{
		Id: w.getNextId(), City: c.Id, Cell: c.Cell,
		Name: name, Units: make(SetOfUnits, 0),
		Targets: make([]Command, 0),
	}
	w.Live.Armies.Add(a)
	c.armies.Add(a)
	return a.Id, nil
}

func (w *World) ArmyGet(id uint64) *Army {
	return w.Live.Armies.Get(id)
}

func (a *Army) PopCommand() {
	a.Targets = a.Targets[1:]
}

func (a *Army) ApplyAgressivity(w *World) {
	// FIXME(jfs): NYI
}

func (a *Army) Move(w *World) error {
	if len(a.Targets) <= 0 {
		return nil
	}

	cmd := a.Targets[0]
	src := a.Cell
	nxt, err := w.Places.NodeGetStep(src, cmd.Cell)
	if err != nil {
		return err
	}
	a.Cell = nxt
	if nxt == cmd.Cell {
		var preventPopping bool
		switch cmd.Action {
		case CmdPause:
			// Just a stop on the way
		case CmdCityAttack:
			// FIXME(jfs): NYI
		case CmdCityDefend:
			preventPopping = true
			// FIXME(jfs): NYI
		case CmdCityOverlord:
			// FIXME(jfs): NYI
		case CmdCityBreak:
			// FIXME(jfs): NYI
		case CmdCityMassacre:
			// FIXME(jfs): NYI
		}
		if !preventPopping {
			a.PopCommand()
		}
		a.ApplyAgressivity(w)
	}
	return nil
}
