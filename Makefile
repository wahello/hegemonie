BASE=github.com/jfsmig/hegemonie
GO=go

all:
	protoc -I pkg/auth   pkg/auth/service.proto --go_out=plugins=grpc:pkg/auth/proto
	protoc -I pkg/region pkg/region/city.proto  --go_out=plugins=grpc:pkg/region/proto_city
	protoc -I pkg/region pkg/region/army.proto  --go_out=plugins=grpc:pkg/region/proto_army
	protoc -I pkg/region pkg/region/admin.proto --go_out=plugins=grpc:pkg/region/proto_admin
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

