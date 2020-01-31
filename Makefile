BASE=github.com/jfsmig/hegemonie
GO=go

all:
	protoc -I pkg/auth   pkg/auth/auth.proto --go_out=plugins=grpc:pkg/auth/proto
	protoc -I pkg/region pkg/region/region.proto  --go_out=plugins=grpc:pkg/region/proto
	$(GO) install $(BASE)

clean:
	$(GO) clean $(BASE)

.PHONY: all clean test bench fmt try

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

