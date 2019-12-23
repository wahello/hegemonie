BASE=github.com/jfsmig/hegemonie

main: hege-front hege-world
all: hege-front hege-world world client mapper

clean:
	-rm -f hege-front hege-world hege-ticker

install: all
	install /usr/local/bin hege-front hege-world

.PHONY: all clean install test try \
	world client mapper hege-front hege-world

mapper:
	go install $(BASE)/common/mapper
world:
	go install $(BASE)/common/world
client:
	go install $(BASE)/common/client
hege-front:
	go install $(BASE)/hege-front
hege-world:
	go install $(BASE)/hege-world

test:
	go test $(BASE)/common/world
	go test $(BASE)/common/client
	go test $(BASE)/common/mapper
	go test $(BASE)/hege-world
	go test $(BASE)/hege-front

try: hege-front hege-world
	ci/run.sh $$PWD/ci/bootstrap-empty.json
