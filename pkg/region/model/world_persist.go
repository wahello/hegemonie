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

func (w *World) Init() {
	w.WLock()
	defer w.WUnlock()

	w.Places.Init()

	if w.NextId <= 0 {
		w.NextId = 1
	}
	w.Live.Armies = make(SetOfArmies, 0)
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

	if !sort.IsSorted(&w.Live.Armies) {
		return errors.New("armies unsorted")
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
	for _, a := range w.Live.Armies {
		if !sort.IsSorted(&a.Units) {
			return errors.New("units unsorted")
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
		if !sort.IsSorted(&a.armies) {
			return errors.New("city armies unsorted")
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
	sort.Sort(&w.Live.Armies)
	sort.Sort(&w.Live.Cities)
	sort.Sort(&w.Live.Fights)
	for _, a := range w.Live.Armies {
		sort.Sort(&a.Units)
	}
	for _, c := range w.Live.Cities {
		sort.Sort(&c.Knowledges)
		sort.Sort(&c.Buildings)
		sort.Sort(&c.Units)
		if c.armies == nil {
			c.armies = make(SetOfArmies, 0)
		} else {
			sort.Sort(&c.armies)
		}
		if c.lieges == nil {
			c.lieges = make(SetOfCities, 0)
		} else {
			sort.Sort(&c.lieges)
		}
	}

	// Link Armies and Cities
	for _, a := range w.Live.Armies {
		if a.City == 0 {
			return errors.New(fmt.Sprintf("Army %v points to no City", a))
		} else if c := w.CityGet(a.City); c == nil {
			return errors.New(fmt.Sprintf("Army %v points to ghost City", a))
		} else {
			c.armies.Add(a)
		}
	}

	// Compute the highest unique ID
	maxId := w.NextId
	for _, u := range w.Definitions.Units {
		if u.Id > maxId {
			maxId = u.Id
		}
	}
	for _, u := range w.Definitions.Buildings {
		if u.Id > maxId {
			maxId = u.Id
		}
	}
	for _, u := range w.Definitions.Knowledges {
		if u.Id > maxId {
			maxId = u.Id
		}
	}
	for _, u := range w.Live.Armies {
		if u.Id > maxId {
			maxId = u.Id
		}
	}
	for _, u := range w.Live.Cities {
		if u.Id > maxId {
			maxId = u.Id
		}
	}

	for _, c := range w.Live.Cities {
		for _, u := range c.Units {
			if u.Id > maxId {
				maxId = u.Id
			}
		}
		for _, u := range c.Knowledges {
			if u.Id > maxId {
				maxId = u.Id
			}
		}
		for _, u := range c.Buildings {
			if u.Id > maxId {
				maxId = u.Id
			}
		}
	}
	for _, a := range w.Live.Armies {
		for _, u := range a.Units {
			if u.Id > maxId {
				maxId = u.Id
			}
		}
	}

	w.NextId = maxId + 1
	return nil
}

func makeSaveFilename() string {
	now := time.Now().Round(1 * time.Second)
	return "save-" + now.Format("20060102_030405")
}

func (w *World) SaveLiveToFiles(basePath string) (string, error) {
	if basePath == "" {
		return "", errors.New("No save path configured")
	}

	p := basePath + "/" + makeSaveFilename()
	if err := os.MkdirAll(p, 0755); err != nil {
		return p, err
	}

	type cfgSection struct {
		path string
		obj  interface{}
	}
	cfgSections := []cfgSection{
		{p + "/map.json", &w.Places},
		{p + "/armies.json", &w.Live.Armies},
		{p + "/cities.json", &w.Live.Cities},
		{p + "/fights.json", &w.Live.Fights},
	}
	for _, section := range cfgSections {
		out, err := os.Create(section.path)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Failed to save the World in [%s]: %s", section.path, err.Error()))
		}
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", " ")
		err = encoder.Encode(section.obj)
		_ = out.Close()
		if err != nil {
			return "", errors.New(fmt.Sprintf("Failed to save the World in [%s]: %s", section.path, err.Error()))
		}
	}

	return p, nil
}

func (w *World) LoadLiveFromFiles(basePath string) error {
	if basePath == "" {
		return errors.New("No save path configured")
	}

	type cfgSection struct {
		path string
		obj  interface{}
	}
	cfgSections := []cfgSection{
		{basePath + "/map.json", &w.Places},
		{basePath + "/armies.json", &w.Live.Armies},
		{basePath + "/cities.json", &w.Live.Cities},
		{basePath + "/fights.json", &w.Live.Fights},
	}
	for _, section := range cfgSections {
		in, err := os.Open(section.path)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to load the World from [%s]: %s", section.path, err.Error()))
		}
		err = json.NewDecoder(in).Decode(section.obj)
		in.Close()
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to load the World from [%s]: %s", section.path, err.Error()))
		}
	}

	return nil
}

func (w *World) SaveDefinitionsToFiles(basePath string) (string, error) {
	if basePath == "" {
		return "", errors.New("No save path configured")
	}

	p := basePath
	if err := os.MkdirAll(p, 0755); err != nil {
		return p, err
	}

	type cfgSection struct {
		path string
		obj  interface{}
	}
	cfgSections := []cfgSection{
		{p + "/config.json", &w.Config},
		{p + "/units.json", &w.Definitions.Units},
		{p + "/buildings.json", &w.Definitions.Buildings},
		{p + "/knowledge.json", &w.Definitions.Knowledges},
	}
	for _, section := range cfgSections {
		out, err := os.Create(section.path)
		if err != nil {
			return "", errors.New(fmt.Sprintf("Failed to save the World in [%s]: %s", section.path, err.Error()))
		}
		encoder := json.NewEncoder(out)
		encoder.SetIndent("", " ")
		err = encoder.Encode(section.obj)
		_ = out.Close()
		if err != nil {
			return "", errors.New(fmt.Sprintf("Failed to save the World in [%s]: %s", section.path, err.Error()))
		}
	}

	return p, nil
}

func (w *World) LoadDefinitionsFromFiles(basePath string) error {
	if basePath == "" {
		return errors.New("No save path configured")
	}

	type cfgSection struct {
		path string
		obj  interface{}
	}
	cfgSections := []cfgSection{
		{basePath + "/config.json", &w.Config},
		{basePath + "/units.json", &w.Definitions.Units},
		{basePath + "/buildings.json", &w.Definitions.Buildings},
		{basePath + "/knowledge.json", &w.Definitions.Knowledges},
	}
	for _, section := range cfgSections {
		in, err := os.Open(section.path)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to load the World from [%s]: %s", section.path, err.Error()))
		}
		err = json.NewDecoder(in).Decode(section.obj)
		in.Close()
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to load the World from [%s]: %s", section.path, err.Error()))
		}
	}

	return nil
}
