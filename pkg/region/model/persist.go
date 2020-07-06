// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"errors"
	"sort"

	"github.com/jfsmig/hegemonie/pkg/utils"
)

func (w *World) SetNotifier(n Notifier) {
	w.notifier = LogEvent(n)
}

func (w *World) Init() {
	w.WLock()
	defer w.WUnlock()

	w.SetNotifier(&noEvt{})
	w.Regions = make(SetOfRegions, 0)
	w.Definitions.Units = make(SetOfUnitTypes, 0)
	w.Definitions.Buildings = make(SetOfBuildingTypes, 0)
	w.Definitions.Knowledges = make(SetOfKnowledgeTypes, 0)
}

func (w *World) Check() error {
	w.RLock()
	defer w.RUnlock()
	if w.notifier == nil || w.mapView == nil {
		return errInvalidState
	}
	if err := w.Definitions.Check(); err != nil {
		return err
	}
	for _, r := range w.Regions {
		if err := r.Check(); err != nil {
			return err
		}
	}
	return nil
}

func (defs *DefinitionsBase) Check() error {
	if !sort.IsSorted(&defs.Knowledges) {
		return errors.New("knowledge types unsorted")
	}
	if !sort.IsSorted(&defs.Buildings) {
		return errors.New("building types unsorted")
	}
	if !sort.IsSorted(&defs.Units) {
		return errors.New("unit types unsorted")
	}

	return nil
}

func (reg *Region) Check() error {
	if !sort.IsSorted(&reg.Cities) {
		return errors.New("cities unsorted")
	}
	if !sort.IsSorted(&reg.Fights) {
		return errors.New("fights unsorted")
	}

	for _, a := range reg.Fights {
		if !sort.IsSorted(&a.Attack) {
			return errors.New("fight attack unsorted")
		}
		if !sort.IsSorted(&a.Defense) {
			return errors.New("fight defense unsorted")
		}
	}
	for _, a := range reg.Cities {
		if !sort.IsSorted(&a.Knowledges) {
			return errors.New("knowledge unsorted")
		}
		if !sort.IsSorted(&a.Buildings) {
			return errors.New("building unsorted")
		}
		if !sort.IsSorted(&a.Units) {
			return errors.New("unit sequence: unsorted")
		}
		if !sort.IsSorted(&a.lieges) {
			return errors.New("city lieges unsorted")
		}
		if !sort.IsSorted(&a.Armies) {
			return errors.New("city armies unsorted")
		}
		for _, a := range a.Armies {
			if !sort.IsSorted(&a.Units) {
				return errors.New("units unsorted")
			}
		}
	}

	return nil
}

func (defs *DefinitionsBase) PostLoad() error {
	sort.Sort(&defs.Knowledges)
	sort.Sort(&defs.Buildings)
	sort.Sort(&defs.Units)
	return nil
}

func (reg *Region) PostLoad() error {
	// Sort all the lookup arrays
	sort.Sort(&reg.Cities)
	sort.Sort(&reg.Fights)

	for _, c := range reg.Cities {
		sort.Sort(&c.Knowledges)
		sort.Sort(&c.Buildings)
		sort.Sort(&c.Units)
		if c.Armies == nil {
			c.Armies = make(SetOfArmies, 0)
		} else {
			sort.Sort(&c.Armies)
		}
		if c.lieges == nil {
			c.lieges = make(SetOfCities, 0)
		} else {
			sort.Sort(&c.lieges)
		}

		for _, a := range c.Armies {
			// Link Armies to their City
			a.City = c
			// FIXME: Link each Army to its Fight
		}
	}

	return nil
}

func (w *World) PostLoad() error {
	w.Definitions.PostLoad()
	sort.Sort(&w.Regions)
	for _, r := range w.Regions {
		r.PostLoad()
	}
	return nil
}

func (defs *DefinitionsBase) Sections(p string) utils.PersistencyMapping {
	if p == "" {
		panic("Invalid path")
	}
	return []utils.CfgSection{
		{p + "/units.json", &defs.Units},
		{p + "/buildings.json", &defs.Buildings},
		{p + "/knowledge.json", &defs.Knowledges},
	}
}

func (reg *Region) Sections(p string) utils.PersistencyMapping {
	if p == "" {
		panic("Invalid path")
	}
	return []utils.CfgSection{
		{p + "/cities.json", &reg.Cities},
		{p + "/fights.json", &reg.Fights},
	}
}

func (w *World) Sections(p string) utils.PersistencyMapping {
	if p == "" {
		panic("Invalid path")
	}
	sections := []utils.CfgSection{
		{p + "/config.json", &w.Config},
	}
	sections = append(sections, w.Definitions.Sections(p+"/_defs")...)
	for _, r := range w.Regions {
		sections = append(sections, r.Sections(p+"/"+r.Name)...)
	}
	return sections
}
