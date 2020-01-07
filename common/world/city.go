// Copyright (C) 2018-2019 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package world

import (
	"errors"
)

func (s *SetOfCities) Len() int {
	return len(*s)
}

func (s *SetOfCities) Less(i, j int) bool {
	return (*s)[i].Id < (*s)[j].Id
}

func (s *SetOfCities) Swap(i, j int) {
	tmp := (*s)[i]
	(*s)[i] = (*s)[j]
	(*s)[j] = tmp
}

func (w *World) CityGet(id uint64) *City {
	for _, c := range w.Live.Cities {
		if c.Id == id {
			return c
		}
	}
	return nil
}

func (w *World) CityCheck(id uint64) bool {
	return w.CityGet(id) != nil
}

func (w *World) CityCreate(loc uint64) (uint64, error) {
	w.rw.Lock()
	defer w.rw.Unlock()

	c := &City{
		Id: w.getNextId(), Cell: loc,
		Units:      make(SetOfUnits, 0),
		Buildings:  make(SetOfBuildings, 0),
		Knowledges: make(SetOfKnowledges, 0),
	}
	w.Live.Cities = append(w.Live.Cities, c)
	return c.Id, nil
}

func (w *World) CityTrain(userId, characterId, cityId, uId uint64) (uint64, error) {
	w.rw.Lock()
	defer w.rw.Unlock()

	pCity, _, _, err := w.CityGetAndCheck(userId, characterId, cityId)
	if err != nil {
		return 0, err
	}
	return pCity.Train(w, uId)
}

func (w *World) CityBuild(userId, characterId, cityId, bId uint64) (uint64, error) {
	w.rw.Lock()
	defer w.rw.Unlock()

	pCity, _, _, err := w.CityGetAndCheck(userId, characterId, cityId)
	if err != nil {
		return 0, err
	}
	return pCity.Build(w, bId)
}

func (w *World) CityStudy(userId, characterId, cityId, kId uint64) (uint64, error) {
	w.rw.Lock()
	defer w.rw.Unlock()

	pCity, _, _, err := w.CityGetAndCheck(userId, characterId, cityId)
	if err != nil {
		return 0, err
	}
	return pCity.Study(w, kId)
}

func (c *City) Unit(id uint64) *Unit {
	for _, b := range c.Units {
		if id == b.Id {
			return b
		}
	}
	return nil
}

func (c *City) Building(id uint64) *Building {
	for _, b := range c.Buildings {
		if id == b.Id {
			return b
		}
	}
	return nil
}

func (c *City) Knowledge(id uint64) *Knowledge {
	for _, b := range c.Knowledges {
		if id == b.Id {
			return b
		}
	}
	return nil
}

func (w *World) CityGetAndCheck(userId, characterId, cityId uint64) (*City, *Character, *Character, error) {
	// Fetch + sanity checks about the city
	pCity := w.CityGet(cityId)
	if pCity == nil {
		return nil, nil, nil, errors.New("Not Found")
	}
	if pCity.Deputy != characterId && pCity.Owner != characterId {
		return nil, nil, nil, errors.New("Forbidden")
	}

	// Fetch + senity checks about the City
	pOwner := w.CharacterGet(pCity.Owner)
	pDeputy := w.CharacterGet(pCity.Deputy)
	if pOwner == nil || pDeputy == nil {
		return nil, nil, nil, errors.New("Not Found")
	}
	if pOwner.User != userId && pDeputy.User != userId {
		return nil, nil, nil, errors.New("Forbidden")
	}

	return pCity, pOwner, pDeputy, nil
}

func (w *World) CityShow(userId, characterId, cityId uint64) (view CityView, err error) {
	w.rw.RLock()
	defer w.rw.RUnlock()

	pCity, pOwner, pDeputy, e := w.CityGetAndCheck(userId, characterId, cityId)
	if e != nil {
		err = e
		return
	}

	view = pCity.Show(w)
	view.Owner.Name = pOwner.Name
	view.Deputy.Name = pDeputy.Name
	return
}

func (c *City) Show(w *World) (view CityView) {
	view.Id = c.Id
	view.Name = c.Name
	view.Owner.Id = c.Owner
	view.Deputy.Id = c.Deputy
	view.Buildings = make([]BuildingView, 0, len(c.Buildings))
	view.Units = make([]UnitView, 0, len(c.Units))
	view.Armies = make([]NamedItem, 0)
	view.Knowledges = make([]KnowledgeView, 0)
	view.KFrontier = make([]KnowledgeType, 0)
	view.BFrontier = make([]BuildingType, 0)
	view.UFrontier = make([]UnitType, 0)

	for _, a := range c.armies {
		view.Armies = append(view.Armies, NamedItem{Id: a.Id, Name: a.Name})
	}

	// Compute the modifiers
	for i := 0; i < ResourceMax; i++ {
		view.Production.Buildings.Mult[i] = 1.0
		view.Production.Knowledge.Mult[i] = 1.0
		view.Production.Troops.Mult[i] = 1.0
		view.Stock.Buildings.Mult[i] = 1.0
		view.Stock.Knowledge.Mult[i] = 1.0
		view.Stock.Troops.Mult[i] = 1.0
	}

	for _, k := range c.Knowledges {
		v := KnowledgeView{}
		v.Id = k.Id
		v.Type = *w.KnowledgeTypeGet(k.Type)
		view.Knowledges = append(view.Knowledges, v)
	}
	for _, b := range c.Buildings {
		v := BuildingView{}
		v.Id = b.Id
		v.Type = *w.GetBuildingType(b.Type)
		view.Buildings = append(view.Buildings, v)
		for i := 0; i < ResourceMax; i++ {
			view.Production.Buildings.Plus[i] += v.Type.Prod.Plus[i]
			view.Production.Buildings.Mult[i] *= v.Type.Prod.Mult[i]
			view.Stock.Buildings.Plus[i] += v.Type.Stock.Plus[i]
			view.Stock.Buildings.Mult[i] *= v.Type.Stock.Mult[i]
		}
	}
	for _, u := range c.Units {
		v := UnitView{}
		v.Id = u.Id
		v.Type = *w.UnitTypeGet(u.Type)
		view.Units = append(view.Units, v)
		for i := 0; i < ResourceMax; i++ {
			view.Production.Troops.Plus[i] += v.Type.Prod.Plus[i]
			view.Production.Troops.Mult[i] *= v.Type.Prod.Mult[i]
		}
	}

	for _, kt := range c.KnowledgeFrontier(w) {
		view.KFrontier = append(view.KFrontier, *kt)
	}
	for _, bt := range c.BuildingFrontier(w) {
		view.BFrontier = append(view.BFrontier, *bt)
	}
	for _, ut := range c.UnitFrontier(w) {
		view.UFrontier = append(view.UFrontier, *ut)
	}

	// Apply all the modifiers on the production
	view.Production.Base = c.Production
	view.Production.Actual = c.Production
	for i := 0; i < ResourceMax; i++ {
		v := float64(view.Production.Base[i])
		v = v * view.Production.Troops.Mult[i]
		v = v * view.Production.Buildings.Mult[i]
		v = v * view.Production.Knowledge.Mult[i]

		vi := int64(v)
		vi = vi + view.Production.Troops.Plus[i]
		vi = vi + view.Production.Buildings.Plus[i]
		vi = vi + view.Production.Knowledge.Plus[i]

		view.Production.Actual[i] = uint64(vi)
	}

	// Apply all the modifiers on the stock
	view.Stock.Base = c.StockCapacity
	view.Stock.Actual = c.StockCapacity
	view.Stock.Usage = c.Stock
	for i := 0; i < ResourceMax; i++ {
		v := float64(view.Stock.Base[i])
		v = v * view.Stock.Troops.Mult[i]
		v = v * view.Stock.Buildings.Mult[i]
		v = v * view.Stock.Knowledge.Mult[i]

		vi := int64(v)
		vi = vi + view.Stock.Troops.Plus[i]
		vi = vi + view.Stock.Buildings.Plus[i]
		vi = vi + view.Stock.Knowledge.Plus[i]

		view.Stock.Actual[i] = uint64(vi)
	}

	return
}

func (c *City) Produce(w *World) {
	// Pre-compute the modified values of Stock and Production
	view := c.Show(w)
	post := view.Stock.Usage

	post.Add(&view.Production.Actual)

	for _, b := range c.Buildings {
		if b.Ticks > 0 {
			bt := w.GetBuildingType(b.Id)
			if post.GreaterOrEqualTo(&bt.Cost) {
				post.Remove(&bt.Cost)
				b.Ticks--
			}
		}
	}

	for _, u := range c.Units {
		if u.Ticks > 0 {
			ut := w.UnitTypeGet(u.Type)
			if post.GreaterOrEqualTo(&ut.Cost) {
				post.Remove(&ut.Cost)
				u.Ticks--
			}
		}
	}

	post.TrimTo(&view.Stock.Actual)
	c.Stock = post
}

// Transfer a Unit from the City to the given Army.
// No check is performed if the City controls the Army.
func (c *City) TransferUnit(w *World, a *Army, idUnit uint64) error {
	pUnit := c.Unit(idUnit)
	if pUnit == nil {
		return errors.New("Unit not found")
	}

	c.Units.Remove(pUnit)
	a.Units.Add(pUnit)
	return nil
}

// Transfer a Unit from the City to the given Army.
// The chain of ownerships (Charcater, City, Army, Unit) is checked.
func (w *World) CityTransferUnit(idUser, idChar, idCity uint64, idUnit, idArmy uint64) error {
	w.rw.Lock()
	defer w.rw.Unlock()

	pCity, _, _, err := w.CityGetAndCheck(idUser, idChar, idCity)
	if err != nil {
		return err
	}

	pArmy := w.ArmyGet(idArmy)
	if pArmy == nil {
		return errors.New("Army not found")
	}

	if pArmy.City != pCity.Id {
		return errors.New("")
	}

	return pCity.TransferUnit(w, pArmy, idUnit)
}

func (w *World) CityCreateArmy(idUser, idChar, idCity uint64, name string) (uint64, error) {
	w.rw.Lock()
	defer w.rw.Unlock()

	pCity, _, _, err := w.CityGetAndCheck(idUser, idChar, idCity)
	if err != nil {
		return 0, err
	}

	return w.ArmyCreate(pCity, name)
}

func (c *City) KnowledgeFrontier(w *World) []*KnowledgeType {
	return w.KnowledgeGetFrontier(c.Knowledges)
}

func (c *City) BuildingFrontier(w *World) []*BuildingType {
	return w.BuildingGetFrontier(c.Buildings, c.Knowledges)
}

// Return a collection of UnitType that may be trained by the current City
// because all the requirements are met.
// Each UnitType 'p' returned validates 'c.UnitAllowed(p)'.
func (c *City) UnitFrontier(w *World) []*UnitType {
	return w.UnitGetFrontier(c.Buildings)
}

// Check the current City has all the requirements to train a Unti of the
// given UnitType.
func (c *City) UnitAllowed(pType *UnitType) bool {
	if pType.RequiredBuilding == 0 {
		return true
	}
	for _, b := range c.Buildings {
		if b.Type == pType.RequiredBuilding {
			return true
		}
	}
	return false
}

// Create a Unit of the given UnitType.
// No check is performed to verify the City has all the requirements.
func (c *City) UnitCreate(w *World, pType *UnitType) uint64 {
	id := w.getNextId()
	u := &Unit{Id: id, Type: pType.Id, Ticks: pType.Ticks, Health: pType.Health}
	c.Units.Add(u)
	return id
}

// Start the training of a Unit of the given UnitType (id).
// The whole chain of requirements will be checked.
func (c *City) Train(w *World, idType uint64) (uint64, error) {
	pType := w.UnitTypeGet(idType)
	if pType == nil {
		return 0, errors.New("Unit Type not found")
	}
	if !c.UnitAllowed(pType) {
		return 0, errors.New("Precondition Failed: no suitable building")
	}

	return c.UnitCreate(w, pType), nil
}

func (c *City) Study(w *World, kId uint64) (uint64, error) {
	pType := w.KnowledgeTypeGet(kId)
	if pType == nil {
		return 0, errors.New("Knowledge Type not found")
	}
	owned := make(map[uint64]bool)
	for _, k := range c.Knowledges {
		if kId == k.Type {
			return 0, errors.New("Already started")
		}
		owned[k.Type] = true
	}
	for _, k := range pType.Conflicts {
		if owned[k] {
			return 0, errors.New("Conflict")
		}
	}
	for _, k := range pType.Requires {
		if !owned[k] {
			return 0, errors.New("Precondition Failed")
		}
	}

	id := w.getNextId()
	c.Knowledges.Add(&Knowledge{Id: id, Type: kId, Ticks: pType.Ticks})
	return id, nil
}

func (c *City) Build(w *World, bId uint64) (uint64, error) {
	pType := w.GetBuildingType(bId)
	if pType == nil {
		return 0, errors.New("Building Type not found")
	}
	if pType.Unique {
		for _, b := range c.Buildings {
			if b.Type == bId {
				return 0, errors.New("Building already present")
			}
		}
	}

	// Check the knowledge requirements are met
	owned := make(map[uint64]bool)
	for _, k := range c.Knowledges {
		owned[k.Type] = true
	}
	for _, k := range pType.Conflicts {
		if owned[k] {
			return 0, errors.New("Conflict")
		}
	}
	for _, k := range pType.Requires {
		if !owned[k] {
			return 0, errors.New("Precondition Failed")
		}
	}

	id := w.getNextId()
	c.Buildings.Add(&Building{Id: id, Type: bId, Ticks: pType.Ticks})
	return id, nil
}
