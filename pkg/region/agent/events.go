// Copyright (C) 2018-2020 Hegemonie's AUTHORS
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package hegemonie_region_agent

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	hegemonie_rpevent_proto "github.com/jfsmig/hegemonie/pkg/event/proto"
	"github.com/jfsmig/hegemonie/pkg/region/model"
	"google.golang.org/grpc"
)

type EventStore struct {
	cnx *grpc.ClientConn
}

type EventArmy struct {
	store  *EventStore `json:"-"`
	charId uint64      `json:"-"`

	SourceCityId uint64 `json:"sourceCityId"`
	SourceCity   string `json:"sourceCity"`

	ArmyId   uint64 `json:"armyId"`
	ArmyName string `json:"army"`

	ArmyCityId   uint64 `json:"armyCityId"`
	ArmyCityName string `json:"armyCity"`

	Src uint64 `json:"src"`
	Dst uint64 `json:"dst"`

	Action string `json:"action"`
}

type EventKnowledge struct {
	store *EventStore
}

type EventUnits struct {
	store *EventStore
}

func (es *EventStore) Army(log *region.City) region.EventArmy {
	return &EventArmy{
		store:        es,
		charId:       log.Owner,
		SourceCity:   log.Name,
		SourceCityId: log.Id,
	}
}

func (es *EventStore) Knowledge(log *region.City) region.EventKnowledge {
	return &EventKnowledge{store: es}
}

func (es *EventStore) Units(log *region.City) region.EventUnits {
	return &EventUnits{store: es}
}

func (evt *EventArmy) Item(a *region.Army) region.EventArmy {
	evt.ArmyId = a.Id
	evt.ArmyName = a.Name
	evt.ArmyCityId = a.City.Id
	evt.ArmyCityName = a.City.Name
	return evt
}

func (evt *EventArmy) Move(src, dst uint64) region.EventArmy {
	evt.Src, evt.Dst = src, dst
	evt.Action = "Move"
	return evt
}

func (evt *EventArmy) NoRoute(src, dst uint64) region.EventArmy {
	evt.Src, evt.Dst = src, dst
	evt.Action = "NoRoute"
	return evt
}

func (evt *EventArmy) Send() {
	var buffer bytes.Buffer
	enc := json.NewEncoder(&buffer)
	enc.SetIndent("", "")
	enc.Encode(evt)

	client := hegemonie_rpevent_proto.NewProducerClient(evt.store.cnx)
	client.Push1(context.Background(), &hegemonie_rpevent_proto.Push1Req{
		CharId:  evt.charId,
		EvtId:   uuid.New().String(),
		Payload: buffer.Bytes(),
	})
}

func (evt *EventKnowledge) Item(c *region.City, kt *region.KnowledgeType) region.EventKnowledge {
	// TODO FIXME
	return evt
}

func (evt *EventKnowledge) Step(current, max uint64) region.EventKnowledge {
	// TODO FIXME
	return evt
}

func (evt *EventKnowledge) Send() {
	// TODO FIXME
}

func (evt *EventUnits) Item(c *region.City, ut *region.UnitType) region.EventUnits {
	// TODO FIXME
	return evt
}

func (evt *EventUnits) Step(current, max uint64) region.EventUnits {
	// TODO FIXME
	return evt
}

func (evt *EventUnits) Send() {
	// TODO FIXME
}
