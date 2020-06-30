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
	"sort"
	"time"
)

func (w *World) SetNotifier(n Notifier) {
	w.notifier = LogEvent(n)
}

func (w *World) Init() {
	w.WLock()
	defer w.WUnlock()

	w.SetNotifier(&noEvt{})
	w.Places.Init()

	w.nextID = 1
	w.Live.Cities = make(SetOfCities, 0)
	w.Live.Fights = make(SetOfFights, 0)
	w.Definitions.Units = make(SetOfUnitTypes, 0)
	w.Definitions.Buildings = make(SetOfBuildingTypes, 0)
	w.Definitions.Knowledges = make(SetOfKnowledgeTypes, 0)
}

func (w *World) Check() error {
	if !sort.IsSorted(&w.Places.Cells) {
		return errors.New("locations unsorted")
	}
	if !sort.IsSorted(&w.Places.Roads) {
		return errors.New("roads unsorted")
	}

	if !sort.IsSorted(&w.Definitions.Knowledges) {
		return errors.New("knowledge types unsorted")
	}
	if !sort.IsSorted(&w.Definitions.Buildings) {
		return errors.New("building types unsorted")
	}
	if !sort.IsSorted(&w.Definitions.Units) {
		return errors.New("unit types unsorted")
	}

	if !sort.IsSorted(&w.Live.Cities) {
		return errors.New("cities unsorted")
	}
	if !sort.IsSorted(&w.Live.Fights) {
		return errors.New("fights unsorted")
	}

	for _, a := range w.Live.Fights {
		if !sort.IsSorted(&a.Attack) {
			return errors.New("fight attack unsorted")
		}
		if !sort.IsSorted(&a.Defense) {
			return errors.New("fight defense unsorted")
		}
	}
	for _, a := range w.Live.Cities {
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

func (w *World) PostLoad() error {
	// Sort all the lookup arrays
	sort.Sort(&w.Places.Cells)
	sort.Sort(&w.Places.Roads)
	sort.Sort(&w.Definitions.Knowledges)
	sort.Sort(&w.Definitions.Buildings)
	sort.Sort(&w.Definitions.Units)
	sort.Sort(&w.Live.Cities)
	sort.Sort(&w.Live.Fights)
	for _, c := range w.Live.Cities {
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

	// Compute the highest unique ID
	maxID := w.nextID
	for _, u := range w.Definitions.Units {
		if u.ID > maxID {
			maxID = u.ID
		}
	}
	for _, u := range w.Definitions.Buildings {
		if u.ID > maxID {
			maxID = u.ID
		}
	}
	for _, u := range w.Definitions.Knowledges {
		if u.ID > maxID {
			maxID = u.ID
		}
	}
	for _, c := range w.Live.Cities {
		if c.ID > maxID {
			maxID = c.ID
		}
		for _, a := range c.Armies {
			if a.ID > maxID {
				maxID = a.ID
			}
			for _, u := range a.Units {
				if u.ID > maxID {
					maxID = u.ID
				}
			}
		}
		for _, u := range c.Units {
			if u.ID > maxID {
				maxID = u.ID
			}
		}
		for _, u := range c.Knowledges {
			if u.ID > maxID {
				maxID = u.ID
			}
		}
		for _, u := range c.Buildings {
			if u.ID > maxID {
				maxID = u.ID
			}
		}
	}

	w.nextID = maxID + 1
	return nil
}

func makeSaveFilename() string {
	now := time.Now().Round(1 * time.Second)
	return "save-" + now.Format("20060102_030405")
}

type persistencyMapping []cfgSection

type cfgSection struct {
	path string
	obj  interface{}
}

func liveSections(p string, w *World) persistencyMapping {
	return []cfgSection{
		{p + "/map.json", &w.Places},
		{p + "/cities.json", &w.Live.Cities},
		{p + "/fights.json", &w.Live.Fights},
	}
}

func defsSections(p string, w *World) persistencyMapping {
	return []cfgSection{
		{p + "/config.json", &w.Config},
		{p + "/units.json", &w.Definitions.Units},
		{p + "/buildings.json", &w.Definitions.Buildings},
		{p + "/knowledge.json", &w.Definitions.Knowledges},
	}
}

func (p persistencyMapping) dump() error {
	for _, section := range p {
		out, err := os.Create(section.path)
		if err != nil {
			return fmt.Errorf("Failed to save the World in [%s]: %s", section.path, err.Error())
		}
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", " ")
		err = encoder.Encode(section.obj)
		_ = out.Close()
		if err != nil {
			return fmt.Errorf("Failed to save the World in [%s]: %s", section.path, err.Error())
		}
	}
	return nil
}

func (p persistencyMapping) load() error {
	for _, section := range p {
		in, err := os.Open(section.path)
		if err != nil {
			return fmt.Errorf("Failed to load the World from [%s]: %s", section.path, err.Error())
		}
		err = json.NewDecoder(in).Decode(section.obj)
		in.Close()
		if err != nil {
			return fmt.Errorf("Failed to load the World from [%s]: %s", section.path, err.Error())
		}
	}
	return nil
}

func (w *World) SaveLiveToFiles(basePath string) (string, error) {
	if basePath == "" {
		return "", errors.New("No save path configured")
	}

	p := basePath + "/" + makeSaveFilename()
	err := os.MkdirAll(p, 0755)
	if err != nil {
		return p, err
	}

	return p, liveSections(p, w).dump()
}

func (w *World) LoadLiveFromFiles(basePath string) error {
	if basePath == "" {
		return errors.New("No save path configured")
	}

	return liveSections(basePath, w).load()
}

func (w *World) SaveDefinitionsToFiles(basePath string) (string, error) {
	if basePath == "" {
		return "", errors.New("No save path configured")
	}

	err := os.MkdirAll(basePath, 0755)
	if err != nil {
		return basePath, err
	}
	return basePath, defsSections(basePath, w).dump()
}

func (w *World) LoadDefinitionsFromFiles(basePath string) error {
	if basePath == "" {
		return errors.New("No save path configured")
	}
	return defsSections(basePath, w).load()
}
