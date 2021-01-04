// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package region

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

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

func walkJson(path string, hook func(path string, decoder *json.Decoder) error) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if !strings.HasSuffix(path, ".json") {
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		if info.Size() <= 2 {
			return nil
		}
		f, err := os.Open(path)
		decoder := json.NewDecoder(f)
		return hook(path, decoder)
	})

}

func (defs *DefinitionsBase) loadUnits(basedir string) error {
	return walkJson(basedir, func(_ string, decoder *json.Decoder) error {
		tmp := make([]*UnitType, 0)
		if err := decoder.Decode(&tmp); err != nil {
			return err
		}
		defs.Units = append(defs.Units, tmp...)
		return nil
	})
}

func (defs *DefinitionsBase) loadKnowledge(basedir string) (err error) {
	return walkJson(basedir, func(_ string, decoder *json.Decoder) error {
		tmp := make([]*KnowledgeType, 0)
		if err = decoder.Decode(&tmp); err != nil {
			return err
		}
		defs.Knowledges = append(defs.Knowledges, tmp...)
		return nil
	})
}

func (defs *DefinitionsBase) loadBuildings(basedir string) (err error) {
	return walkJson(basedir, func(_ string, decoder *json.Decoder) error {
		tmp := make([]*BuildingType, 0)
		if err = decoder.Decode(&tmp); err != nil {
			return err
		}
		defs.Buildings = append(defs.Buildings, tmp...)
		return nil
	})
}

func (defs *DefinitionsBase) load(path string) (err error) {
	err = defs.loadUnits(path + "/units")
	if err == nil {
		err = defs.loadKnowledge(path + "/knowledge")
	}
	if err == nil {
		err = defs.loadBuildings(path + "/buildings")
	}

	if err == nil {
		sort.Sort(&defs.Knowledges)
		sort.Sort(&defs.Buildings)
		sort.Sort(&defs.Units)
	}
	return err
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
	sort.Sort(&w.Regions)
	for _, r := range w.Regions {
		if err := r.PostLoad(); err != nil {
			return err
		}
	}
	return nil
}

func (w *World) LoadDefinitions(basedir string) (err error) {
	err = w.Definitions.load(basedir)
	if err != nil {
		return fmt.Errorf("invalid world from [%s]: %v", basedir, err)
	}

	err = w.Definitions.Check()
	if err != nil {
		return fmt.Errorf("inconsistent world from [%s]: %v", basedir, err)
	}

	return nil
}

func (w *World) LoadRegions(basedir string) error {
	return walkJson(basedir, func(path string, decoder *json.Decoder) error {
		reg := &Region{}
		err := decoder.Decode(&reg)
		if err != nil {
			return fmt.Errorf("region decoding error [%s]: %s", path, err.Error())
		}
		w.Regions.Add(reg)
		return nil
	})
}
