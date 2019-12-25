// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

func (w *World) UnitGet(id uint64) *Unit {
	for _, c := range w.Units {
		if c.Id == id {
			return c
		}
	}
	return nil
}

func (w *World) UnitGetType(id uint64) *UnitType {
	for _, c := range w.UnitTypes {
		if c.Id == id {
			return c
		}
	}
	return nil
}
