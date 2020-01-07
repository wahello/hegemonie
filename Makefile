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

test:
	$(GO) test $(BASE)

try: all
	./ci/run.sh $$PWD/ci/bootstrap
