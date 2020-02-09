BASE=github.com/jfsmig/hegemonie
GO=go

all: prepare
	$(GO) install $(BASE)

prepare:
	protoc -I pkg/auth   pkg/auth/auth.proto --go_out=plugins=grpc:pkg/auth/proto
	protoc -I pkg/region pkg/region/region.proto  --go_out=plugins=grpc:pkg/region/proto

clean:
	$(GO) clean $(BASE)

.PHONY: all prepare clean test bench fmt try

fmt: all
	find * -type f -name '*.go' \
		| grep -v '.pb.go$$' | while read F ; do dirname $$F ; done \
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
