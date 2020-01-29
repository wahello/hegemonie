#https://dev.to/ivan/go-build-a-minimal-docker-image-in-just-three-steps-514i

FROM golang:1.13.0-stretch AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOPATH=/gopath

WORKDIR /dist

#RUN mkdir -p /gopath/src/github.com/jfsmig/hegemonie
COPY . /gopath/src/github.com/jfsmig/hegemonie

# Build & Install the code
RUN set -x \
&& cd /gopath/src/github.com/jfsmig/hegemonie \
&& go mod download \
&& go build -o /dist/hegemonie

# Install the dependencies
RUN set -x \
&& mkdir -p /dist/lib64 \
&& ldd ./hegemonie | tr -s '[:blank:]' '\n' | grep '^/' | \
   xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;' \
&& cp /lib64/ld-linux-x86-64.so.2 /dist/lib64/

# Install static data
RUN set -x \
&& cd /gopath/src/github.com/jfsmig/hegemonie \
&& mkdir -p /data/templates \
&& cp -r pkg/web/templates /data/ \
&& cp -r pkg/web/static /data/static

# Create the minimal runtime image
FROM scratch
COPY --chown=0:0 --from=builder /dist /
COPY --chown=65534:0 --from=builder /data /data
EXPOSE 8080/tcp
USER 65534
WORKDIR /data
ENTRYPOINT ["/hegemonie"]
