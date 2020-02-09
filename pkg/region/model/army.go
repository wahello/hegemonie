// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"errors"
	"log"
	"math/rand"
	"sort"
)

func (s SetOfArmies) Len() int           { return len(s) }
func (s SetOfArmies) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s SetOfArmies) Less(i, j int) bool { return s[i].Id < s[j].Id }

func (s SetOfArmies) Check() error {
	if !sort.IsSorted(s) {
		return errors.New("Unsorted")
	}
	var lastId uint64
	for _, a := range s {
		if lastId == a.Id {
			return errors.New("Dupplicate ID")
		}
		lastId = a.Id
	}
	return nil
}

func (s SetOfArmies) GetIndex(id uint64) int {
	for idx, a := range s {
		if a.Id == id {
			return idx
		}
	}
	return len(s)
}

func (s SetOfArmies) Get(id uint64) *Army {
	idx := s.GetIndex(id)
	if idx < len(s) {
		return s[idx]
	}
	return nil
}

func (s *SetOfArmies) Add(a *Army) {
	*s = append(*s, a)
	sort.Sort(s)
}

func (s *SetOfArmies) Remove(a *Army) {
	idx := s.GetIndex(a.Id)
	if idx >= 0 && idx < len(*s) {
		if len(*s) == 1 {
			*s = (*s)[:0]
		} else {
			s.Swap(idx, s.Len()-1)
			*s = (*s)[:s.Len()-1]
			sort.Sort(*s)
		}
	}
}

func (a *Army) GetId() uint64 { return a.Id }

func (a *Army) PopCommand() {
	a.Targets = a.Targets[1:]
}

func (a *Army) ApplyAgressivity(w *World) {
	// FIXME(jfs): NYI
}

func (a *Army) Move(w *World) {
	if a.Deleted {
		return
	}

	if len(a.Targets) <= 0 {
		// The Army has no command pending. It just stays.
	} else {
		cmd := a.Targets[0]
		src := a.Cell
		dst := cmd.Cell

		nxt, err := w.Places.PathNextStep(src, dst)
		if err != nil {
			log.Println("Map error:", err.Error())
		} else if nxt == 0 {
			// FIXME(jfs): Notify the City that there is no route
		} else {
			a.Cell = nxt
			// FIXME(jfs): Notify a.City of the movement
			// FIXME(jfs): Notify the local City of the passage
		}

		// pSrc := w.Places.CellGet(src)
		var pLocalCity *City
		pLocal := w.Places.CellGet(a.Cell)
		if pLocal != nil {
			pLocalCity = w.CityGet(pLocal.City)
		}

		if nxt == dst {
			var preventPopping bool
			switch cmd.Action {
			case CmdPause:
				// Just a stop on the way
			case CmdCityAttack:
				a.JoinCityAttack(w, pLocalCity)
			case CmdCityDefend:
				if a.JoinCityDefence(w, pLocalCity) {
					preventPopping = true
				}
			case CmdCityOverlord:
				a.Conquer(w, pLocalCity)
			case CmdCityBreak:
				a.BreakBuilding(w, pLocalCity)
			case CmdCityMassacre:
				a.Massacre(w, pLocalCity)
			case CmdCityDeposit:
				a.Deposit(w, pLocalCity)
			case CmdCityDisband:
				a.Disband(w, pLocalCity)
			}
			if !preventPopping {
				a.PopCommand()
			}
		}
	}
	a.ApplyAgressivity(w)
}

func (a *Army) Deposit(w *World, pCity *City) {
	if pCity == nil {
		panic("Impossible action: nil city")
	}

	pCity.Stock.Add(a.Stock)
	a.Stock.Zero()

	// FIXME(jfs): Popularities

	// FIXME(jfs): Notify pLocalCity
	// FIXME(jfs): Notify a.City
}

func (a *Army) Massacre(w *World, pCity *City) {
	if pCity == nil {
		panic("Impossible action: nil city")
	}

	pCity.TicksMassacres++

	// FIXME(jfs): Popularities
	// FIXME(jfs): Notify pLocalCity
	// FIXME(jfs): Notify a.City
}

func (a *Army) Disband(w *World, pCity *City) {
	if pCity == nil {
		panic("Impossible action: nil city")
	}

	nb := len(a.Units)
	if nb > 0 {
		for _, u := range a.Units {
			pCity.Units = append(pCity.Units, u)
		}
		sort.Sort(pCity.Units)
		a.Units = a.Units[:0]
		a.Deleted = true

		// FIXME(jfs): Notify pCity the arrival of 'nb' units
		// FIXME(jfs): Notify a.City the transfer of 'nb' units
	}
}

func (a *Army) BreakBuilding(w *World, pCity *City) {
	if pCity == nil {
		panic("Impossible action: nil city")
	}

	idx := rand.Intn(len(pCity.Buildings))
	pCity.Buildings[idx].Deleted = true

	// FIXME(jfs): Popularities
	// FIXME(jfs): Notify pLocalCity
	// FIXME(jfs): Notify a.City
}

func (a *Army) Conquer(w *World, pCity *City) {
	if pCity == nil {
		panic("Impossible action: nil city")
	}

	pOverlord := w.CityGet(a.City)
	if pOverlord == nil {
		panic("Impossible action: nil overlord")
	}

	pOverlord.ConquerCity(w, pCity)
}

func (a *Army) JoinCityDefence(w *World, pCity *City) bool {
	if pCity.Assault == nil {
		return false
	}
	if pCity.Assault.Cell != a.Cell {
		panic("inconsistency")
	}
	if pCity.Assault.Cell != pCity.Cell {
		panic("inconsistency")
	}

	w.Live.Armies.Remove(a)
	pCity.Assault.Defense.Add(a)
	a.Fight = pCity.Assault.Id
	return true
}

func (a *Army) JoinCityAttack(w *World, pCity *City) {
	if pCity.Assault == nil {
		pCity.Assault = &Fight{
			Id: w.getNextId(), Cell: pCity.Cell,
			Defense: make(SetOfArmies, 0),
			Attack:  make(SetOfArmies, 0)}
		pCity.Assault.Defense.Add(pCity.MakeDefence(w))
	}

	if pCity.Assault.Cell != a.Cell {
		panic("inconsistency")
	}
	if pCity.Assault.Cell != pCity.Cell {
		panic("inconsistency")
	}

	w.Live.Armies.Remove(a)
	pCity.Assault.Attack.Add(a)
	a.Fight = pCity.Assault.Id
}

// Leave the Fight as a loser
func (a *Army) Flea(w *World) error {
	return errors.New("NYI")
}

// Change the side in the Fight.
// If the Army was defending, it becomes an attacker, if it was an attacker
// it becomes a defender.
func (a *Army) Flip(w *World) error {
	return errors.New("NYI")
}

func (w *World) ArmyCreate(c *City, name string) (*Army, error) {
	a := &Army{
		Id: w.getNextId(), City: c.Id, Cell: c.Cell,
		Name: name, Units: make(SetOfUnits, 0),
		Targets: make([]Command, 0),
	}
	w.Live.Armies.Add(a)
	c.armies.Add(a)
	return a, nil
}

func (w *World) ArmyGet(id uint64) *Army {
	return w.Live.Armies.Get(id)
}
