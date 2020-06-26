// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"errors"
	"github.com/jfsmig/hegemonie/pkg/utils"
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
	if a.Fight != 0 {
		return
	}

	if len(a.Targets) <= 0 {
		// The Army has no command pending. It just stays.
	} else {
		cmd := a.Targets[0]
		src := a.Cell
		dst := cmd.Cell

		// pSrc := w.Places.CellGet(src)
		var pLocalCity *City
		pLocal := w.Places.CellGet(a.Cell)
		if pLocal != nil {
			pLocalCity = w.CityGet(pLocal.City)
		}

		nxt, err := w.Places.PathNextStep(src, dst)
		if err != nil || nxt == 0 {
			if err != nil {
				utils.Logger.Warn().Err(err).Uint64("src", src).Uint64("dst", dst).Send()
			}
			w.notifier.Army(a.City).Item(a).NoRoute(src, dst).Send()
		} else {
			a.Cell = nxt
			w.notifier.Army(a.City).Item(a).Move(src, dst).Send()
			if pLocalCity != nil && a.City.ID != pLocalCity.ID {
				w.notifier.Army(pLocalCity).Item(a).Move(src, dst).Send()
			}
		}

		if nxt == dst {
			var preventPopping bool
			switch cmd.Action {
			case CmdMove:
				// Just a stop on the way
			case CmdCityAttack:
				a.JoinCityAttack(w, pLocalCity)
			case CmdCityDefend:
				if a.JoinCityDefence(w, pLocalCity) {
					preventPopping = true
				}
			case CmdCityDisband:
				a.Disband(w, pLocalCity, true)
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

func (a *Army) Disband(w *World, pCity *City, shouldNotify bool) {
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

		if shouldNotify {
			// FIXME(jfs): Notify pCity the arrival of 'nb' units
			// FIXME(jfs): Notify a.City the transfer of 'nb' units
		}
	}
}

func (a *Army) BreakBuilding(w *World, pCity *City) {
	if pCity == nil {
		panic("Impossible action: nil city")
	}

	idx := rand.Intn(len(pCity.Buildings))
	pCity.Buildings.Remove(pCity.Buildings[idx])

	// FIXME(jfs): Popularities
	// FIXME(jfs): Notify pLocalCity
	// FIXME(jfs): Notify a.City
}

func (a *Army) Conquer(w *World, pCity *City) {
	if pCity == nil {
		panic("Impossible action: nil city")
	}
	a.City.ConquerCity(w, pCity)
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

	a.Fight = pCity.Assault.ID
	pCity.Assault.Defense[a.ID] = a

	return true
}

func (a *Army) JoinCityAttack(w *World, pCity *City) {
	if pCity.Assault == nil {
		pCity.Assault = &Fight{
			ID: w.getNextID(), Cell: pCity.Cell,
			Defense: make(SetOfArmies, 0),
			Attack:  make(SetOfArmies, 0)}
		if def, _ := pCity.CreateArmyDefence(w); def != nil {
			def.Fight = pCity.Assault.ID
			pCity.Assault.Defense[def.ID] = def
		}
	}

	if pCity.Assault.Cell != a.Cell {
		panic("inconsistency")
	}
	if pCity.Assault.Cell != pCity.Cell {
		panic("inconsistency")
	}

	a.Fight = pCity.Assault.ID
	pCity.Assault.Attack[a.ID] = a
}

// Leave the Fight as a loser
func (a *Army) Flea(w *World) error {
	return errors.New("Flea NYI")
}

// Change the side in the Fight.
// If the Army was defending, it becomes an attacker, if it was an attacker
// it becomes a defender.
func (a *Army) Flip(w *World) error {
	return errors.New("Flip NYI")
}

func (a *Army) DeferAttack(w *World, t *MapVertex) error {
	a.Targets = append(a.Targets, Command{Action: CmdCityAttack, Cell: t.ID})
	return nil
}

func (a *Army) DeferDefend(w *World, t *MapVertex) error {
	a.Targets = append(a.Targets, Command{Action: CmdCityDefend, Cell: t.ID})
	return nil
}

func (a *Army) DeferDisband(w *World, t *MapVertex) error {
	a.Targets = append(a.Targets, Command{Action: CmdCityDisband, Cell: t.ID})
	return nil
}

func (a *Army) DeferMove(w *World, t *MapVertex) error {
	a.Targets = append(a.Targets, Command{Action: CmdMove, Cell: t.ID})
	return nil
}

func (a *Army) DeferWait(w *World, t *MapVertex) error {
	a.Targets = append(a.Targets, Command{Action: CmdWait, Cell: t.ID})
	return nil
}
