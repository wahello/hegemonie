// Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package evtbacklocal

import (
	"bytes"
	"fmt"
	"github.com/jfsmig/hegemonie/pkg/utils"
	"github.com/tecbot/gorocksdb"
	"math"
	"strconv"
	"strings"
	"time"
)

type Backend struct {
	db *gorocksdb.DB
}

type Item struct {
	CharID  string
	When    uint64
	ID      string
	Payload []byte
}

func Open(path string) (*Backend, error) {
	options := gorocksdb.NewDefaultOptions()
	options.SetCreateIfMissing(true)

	db, err := gorocksdb.OpenDb(options, path)
	if err != nil {
		return nil, err
	}

	return &Backend{db: db}, nil
}

func (b *Backend) Push1(charID string, id string, payload []byte) error {
	opts := gorocksdb.NewDefaultWriteOptions()
	opts.SetSync(false)
	defer opts.Destroy()

	when := math.MaxUint64 - uint64(time.Now().UnixNano())
	k := fmt.Sprintf("%s/%16X/%s", charID, when, id)
	utils.Logger.Warn().Bytes("key", []byte(k)).Msg("PUSH")
	return b.db.Put(opts, []byte(k), payload)
}

func (b *Backend) Ack1(charID string, when uint64, id string) error {
	opts := gorocksdb.NewDefaultWriteOptions()
	opts.SetSync(false)
	defer opts.Destroy()

	w := math.MaxUint64 - when
	k := fmt.Sprintf("%s/%16X/%s", charID, w, id)
	utils.Logger.Warn().Bytes("key", []byte(k)).Msg("DEL")
	return b.db.Delete(opts, []byte(k))
}

func (b *Backend) List(charID string, when uint64, max uint32) ([]Item, error) {
	if max <= 0 {
		max = 100
	} else if max > 1000 {
		max = 1000
	}

	var w uint64
	if when == 0 {
		w = 0
	} else {
		w = math.MaxUint64 - when
	}

	prefix := []byte(fmt.Sprintf("%s/", charID))
	needle := []byte(fmt.Sprintf("%s/%016X/", charID, w))

	opts := gorocksdb.NewDefaultReadOptions()
	defer opts.Destroy()
	opts.SetFillCache(true)
	opts.SetVerifyChecksums(false)
	iterator := b.db.NewIterator(opts)
	iterator.Seek(needle)

	var err error
	out := make([]Item, 0)
	for ; iterator.Valid(); iterator.Next() {
		k := iterator.Key().Data()
		if !bytes.HasPrefix(k, prefix) {
			break
		}
		tokens := strings.SplitN(string(k), "/", 3)
		when, err = strconv.ParseUint(tokens[1], 16, 64)
		out = append(out, Item{
			CharID:  charID,
			When:    when,
			ID:      tokens[2],
			Payload: iterator.Value().Data(),
		})
	}

	if err != nil {
		return nil, err
	}
	return out, nil
}
