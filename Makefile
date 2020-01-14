BASE=github.com/jfsmig/hegemonie
GO=go

all:
	$(GO) install $(BASE)

clean:
	$(GO) clean $(BASE)

.PHONY: all clean test fmt try \
	world client mapper hege-front hege-world

fmt:
	$(GO) fmt $(BASE)
	$(GO) fmt $(BASE)/common/client
	$(GO) fmt $(BASE)/common/mapper
	$(GO) fmt $(BASE)/common/world
	$(GO) fmt $(BASE)/front
	$(GO) fmt $(BASE)/region

test:
	$(GO) test -v $(BASE)
	$(GO) test -v $(BASE)/common/client
	$(GO) test -v $(BASE)/common/mapper
	$(GO) test -v $(BASE)/common/world
	$(GO) test -v $(BASE)/front
	$(GO) test -v $(BASE)/region

bench:
	$(GO) test -bench=. -v $(BASE)
	$(GO) test -bench=. -v $(BASE)/common/client
	$(GO) test -bench=. -v $(BASE)/common/mapper
	$(GO) test -bench=. -v $(BASE)/common/world
	$(GO) test -bench=. -v $(BASE)/front
	$(GO) test -bench=. -v $(BASE)/region

try: all
	./ci/run.sh $$PWD/ci/bootstrap
