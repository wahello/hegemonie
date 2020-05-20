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
	if a.Fight != 0 {
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
			case CmdCityLiberate:
				a.Liberate(w, pLocalCity)
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

func (a *Army) Liberate(w *World, pCity *City) {
	if pCity == nil {
		panic("Impossible action: nil city")
	}

	pOverlord := w.CityGet(a.City)
	if pOverlord == nil {
		panic("Impossible action: nil overlord")
	}

	pOverlord.LiberateCity(w, pCity)
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

	a.Fight = pCity.Assault.Id
	pCity.Assault.Defense[a.Id] = a
	delete(w.Live.Armies, a.Id)

	return true
}

func (a *Army) JoinCityAttack(w *World, pCity *City) {
	if pCity.Assault == nil {
		pCity.Assault = &Fight{
			Id: w.getNextId(), Cell: pCity.Cell,
			Defense: make(map[uint64]*Army),
			Attack:  make(map[uint64]*Army)}
		def := pCity.MakeDefence(w)
		def.Fight = pCity.Assault.Id
		pCity.Assault.Defense[def.Id] = def
	}

	if pCity.Assault.Cell != a.Cell {
		panic("inconsistency")
	}
	if pCity.Assault.Cell != pCity.Cell {
		panic("inconsistency")
	}

	a.Fight = pCity.Assault.Id
	pCity.Assault.Attack[a.Id] = a
	delete(w.Live.Armies, a.Id)
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

func (a *Army) DeferAttack(w *World, t *City) error {
	//FIXME(jfs):
	return errors.New("NYI")
}

func (a *Army) DeferDefend(w *World, t *City) error {
	//FIXME(jfs):
	return errors.New("NYI")
}

func (a *Army) DeferBreak(w *World, t *City) error {
	//FIXME(jfs):
	return errors.New("NYI")
}

func (a *Army) DeferDeposit(w *World, t *City) error {
	//FIXME(jfs):
	return errors.New("NYI")
}

func (a *Army) DeferDisband(w *World, t *City) error {
	//FIXME(jfs):
	return errors.New("NYI")
}

func (a *Army) DeferMassacre(w *World, t *City) error {
	//FIXME(jfs):
	return errors.New("NYI")
}

func (a *Army) DeferConquer(w *World, t *City) error {
	//FIXME(jfs):
	return errors.New("NYI")
}

func (a *Army) DeferLiberate(w *World, t *City) error {
	//FIXME(jfs):
	return errors.New("NYI")
}
