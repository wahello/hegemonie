BASE=github.com/jfsmig/hegemonie
GO=go
PROTOC=protoc

AUTO= pkg/region/model/world_auto.go
AUTO+= pkg/auth/proto/auth.pb.go
AUTO+= pkg/region/proto/region.pb.go

all: prepare
	$(GO) install $(BASE)
	$(GO) install $(BASE)/cmd/hege-mapper

prepare: $(AUTO)

pkg/region/model/%_auto.go: pkg/region/model/core.go cmd/gen-set/main.go
	-rm $@
	$(GO) generate github.com/jfsmig/hegemonie/pkg/region/model

pkg/auth/proto/%.pb.go: pkg/auth/auth.proto
	$(PROTOC) -I pkg/auth   pkg/auth/auth.proto --go_out=plugins=grpc:pkg/auth/proto

pkg/region/proto/%.pb.go: pkg/region/region.proto
	$(PROTOC) -I pkg/region pkg/region/region.proto  --go_out=plugins=grpc:pkg/region/proto

clean:
	-rm $(AUTO)

.PHONY: all prepare clean test bench fmt try

fmt:
	find * -type f -name '*.go' \
		| grep -v -e '_auto.go$$' -e '.pb.go$$' \
		| while read F ; do dirname $$F ; done \
		| sort | uniq | while read D ; do ( set -x ; cd $$D && go fmt ) done

test: all
	find * -type f -name '*_test.go' \
		| while read F ; do dirname $$F ; done \
		| sort | uniq | while read D ; do ( set -x ; cd $$D && go test ) done

bench: all
	find * -type f -name '*_test.go' \
		| while read F ; do dirname $$F ; done \
		| sort | uniq | while read D ; do ( set -x ; cd $$D && go -bench=. test ) done

try: all
	./ci/run.sh $$PWD/ci/bootstrap

img_tag:
	 ( export L='(C) Quentin Minten / CC BY-NC-SA 3.0' ; \
		for F in website/www/static/img0/quentin-minten*/*.jpg ; do \
			BN=$(basename $$F) ; \
			convert img0/$$BN -gravity south -stroke '#000C' -strokewidth 2 -annotate 0 "$L" -stroke  none -fill yellow -annotate 0 "$L" website/www/static/img/$$BN ; \
		done )
