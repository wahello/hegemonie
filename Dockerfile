#------------------------------------------------------------------------------
# Install the system dependencies

FROM golang:1.15-buster AS dependencies
LABEL maintainer="Jean-Francois SMIGIELSKI <jf.smigielski@gmail.com>"
USER 0
RUN set -x \
&& apt-get update -y \
&& apt-get install -y --no-install-recommends \
  make \
  protobuf-compiler \
  librocksdb-dev \
  librocksdb5.17

WORKDIR /
ENTRYPOINT ["/bin/bash"]


#------------------------------------------------------------------------------
# Build the binaries

FROM dependencies AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOPATH=/gopath
USER 0
WORKDIR /dist

COPY Makefile LICENSE AUTHORS.md go.sum go.mod \
  /gopath/src/github.com/jfsmig/hegemonie/
COPY pkg /gopath/src/github.com/jfsmig/hegemonie/pkg
COPY api /gopath/src/github.com/jfsmig/hegemonie/api
COPY cmd /gopath/src/github.com/jfsmig/hegemonie/cmd

# Build & Install the code
RUN set -x \
&& cd /gopath/src/github.com/jfsmig/hegemonie \
&& go mod download

RUN set -x \
&& cd /gopath/src/github.com/jfsmig/hegemonie \
&& make \
&& cp -p -v /gopath/bin/hege /dist

# Install the dependencies.
# Inspired by https://dev.to/ivan/go-build-a-minimal-docker-image-in-just-three-steps-514i
RUN set -x \
&& mkdir -p /dist/lib64 \
&& ldd /dist/hege | tr -s '[:blank:]' '\n' \
 | grep '^/' | grep -v '^/dist' | sort | uniq \
 | xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;' \
&& cp /lib64/ld-linux-x86-64.so.2 /dist/lib64/

# Mangle the maps to build ther raw shape based on the seed definitions
# JFS: we do this here because it is very fast to execute and it benefits
#      from the rich bash environment.
COPY docs/maps        /data/maps
COPY docs/definitions /data/defs
COPY docs/lang        /data/lang
RUN set -ex \
&& D=/data/maps \
&& HEGE=/gopath/bin/hege \
&& ls "$D" | \
   grep '.seed.json$' | \
   while read F ; do echo "$F" "$F" ; done | \
   sed -r 's/^(\S+).seed.json /\1.final.json /' | \
   while read FINAL SEED ; do \
    echo "$D" "$SEED" "$FINAL" ; \
    "$HEGE" tools map init < "$D/$SEED" | "$HEGE" tools map normalize > "$D/$FINAL" ; \
  done

WORKDIR /
ENTRYPOINT ["/bin/bash"]



#------------------------------------------------------------------------------
# Create the minimal runtime image

FROM scratch as runtime
USER 0
COPY --chown=0:0 --from=builder /dist /
# Expose each default port of each module. --> pkg/utils/constants.go
EXPOSE 8080/tcp
EXPOSE 8081/tcp
EXPOSE 8082/tcp
EXPOSE 8083/tcp
EXPOSE 8084/tcp
USER 65534
WORKDIR /
ENTRYPOINT ["/hege"]



#------------------------------------------------------------------------------
# Create a bundle with minimal runtime plus demonsration data

FROM runtime as demo
USER 0
COPY --chown=65534:0 --from=builder /data/defs/        /data/defs/
COPY --chown=65534:0 --from=builder /data/lang/        /data/lang/
COPY --chown=65534:0 --from=builder /data/maps/        /data/maps/
COPY --chown=65534:0                pkg/web/templates  /data/templates/
COPY --chown=65534:0                pkg/web/static     /data/static/
USER 65534
WORKDIR /data
ENTRYPOINT ["/hege"]

