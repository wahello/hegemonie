# Copyright (c) 2018-2021 Contributors as noted in the AUTHORS file
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

BASE=github.com/jfsmig/hegemonie
GO=go
PROTOC=protoc
COV_OUT=coverage.txt

AUTO=
# gen-set
AUTO+= pkg/gen-set/genset_auto_test.go
AUTO+= pkg/map/graph/map_auto.go
AUTO+= pkg/region/model/world_auto.go
# grpc
AUTO+= pkg/map/proto/map_grpc.pb.go
AUTO+= pkg/map/proto/map.pb.go
AUTO+= pkg/event/proto/event_grpc.pb.go
AUTO+= pkg/event/proto/event.pb.go
AUTO+= pkg/region/proto/region_grpc.pb.go
AUTO+= pkg/region/proto/region.pb.go
AUTO+= pkg/healthcheck/healthcheck_grpc.pb.go
AUTO+= pkg/healthcheck/healthcheck.pb.go

default: hege

all: prepare hege

gen-set: pkg/gen-set/gen-set.go
	$(GO) install $(BASE)/pkg/gen-set

hege: gen-set
	$(GO) install $(BASE)/pkg/hege

.PHONY: all default prepare clean clean-auto clean-coverage test bench fmt docker try hege

prepare: $(AUTO)

pkg/gen-set/genset_auto_test.go: pkg/gen-set/genset_test.go gen-set
	-rm $@
	$(GO) generate github.com/jfsmig/hegemonie/pkg/gen-set

pkg/map/graph/map_auto.go: pkg/map/graph/map.go gen-set
	-rm $@
	$(GO) generate github.com/jfsmig/hegemonie/pkg/map/graph

pkg/region/model/world_auto.go: pkg/region/model/types.go gen-set
	-rm $@
	$(GO) generate github.com/jfsmig/hegemonie/pkg/region/model

define grpc_generate
	$(PROTOC) -I api --go_out=$(1) --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative --go-grpc_out=$(1) $(2)
endef

pkg/map/proto/%.pb.go: api/map.proto
	$(call grpc_generate,pkg/map/proto,api/map.proto)

pkg/region/proto/%.pb.go: api/region.proto
	$(call grpc_generate,pkg/region/proto,api/region.proto)

pkg/event/proto/%.pb.go: api/event.proto
	$(call grpc_generate,pkg/event/proto,api/event.proto)

pkg/healthcheck/%.pb.go: api/healthcheck.proto
	$(call grpc_generate,pkg/healthcheck,api/healthcheck.proto)

clean: clean-auto clean-coverage

fmt:
	go list ./... | grep -v -e attic -e vendor | while read D ; do go fmt $$D ; done

clean-auto:
	-rm $(AUTO)

clean-coverage:
	-rm profile.out $(COV_OUT)

test: all clean-coverage
	set -e ; go list ./... | grep -v -e attic -e vendor | while read D ; do go test -race -coverprofile=profile.out -covermode=atomic $$D ; if [ -f profile.out ] ; then cat profile.out >> $(COV_OUT) ; fi ; done

benchmark: all clean-coverage
	set -e ; go list ./... | grep -v -e attic -e vendor | while read D ; do go test -race -coverprofile=profile.out -covermode=atomic -bench=$$D $$D ; if [ -f profile.out ] ; then cat profile.out >> $(COV_OUT) ; fi ;  done

