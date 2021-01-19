#------------------------------------------------------------------------------
# Install the system dependencies

FROM golang:1.15-buster AS dependencies
LABEL maintainer="Jean-Francois SMIGIELSKI <jf.smigielski@gmail.com>"
USER 0
WORKDIR /
RUN apt-get update -y \
&& apt-get install -y --no-install-recommends \
  make \
  protobuf-compiler \
  librocksdb-dev \
  librocksdb5.17 \
&& apt-get clean \
&& rm -rf /var/lib/apt/lists/*



#------------------------------------------------------------------------------
# Build the binaries

FROM dependencies AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=1 \
    GOPATH=/gopath

USER 0
WORKDIR /gopath/src/github.com/jfsmig/hegemonie

# Build & Install the code, then extract all the system deps. Inspired by:
# https://dev.to/ivan/go-build-a-minimal-docker-image-in-just-three-steps-514i
COPY go.sum go.mod /gopath/src/github.com/jfsmig/hegemonie/
RUN go mod download

COPY Makefile LICENSE AUTHORS.md /gopath/src/github.com/jfsmig/hegemonie/
COPY pkg /gopath/src/github.com/jfsmig/hegemonie/pkg
COPY api /gopath/src/github.com/jfsmig/hegemonie/api
RUN make \
&& mkdir /dist \
&& cp -p -v /gopath/bin/hege /dist/ \
&& cp -p -v /gopath/bin/hege /usr/bin/

WORKDIR /dist
RUN set -ex \
&& mkdir -p /dist/lib64 \
&& ldd /dist/hege | tr -s '[:blank:]' '\n' | grep '^/' | sort | uniq \
 | xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;' \
&& cp /lib64/ld-linux-x86-64.so.2 /dist/lib64/

# Prepare the env for the demo env.
# Despite this isn't usable by the current container, it benefits from
# the rich Bash env of the complete container. It helps building a thin
# container for the demo with just the binary tools and the sample data,
# without any distribution burden.
COPY bin       /usr/bin
COPY docs/maps /etc/hegemonie/maps

RUN set -ex \
&& hege-map-transform.sh /etc/hegemonie/maps /etc/hegemonie/maps

RUN set -ex \
&& mkdir /etc/hegemonie/pki \
&& hege-pki-ca.sh  /etc/hegemonie/pki \
&& hege-pki-srv.sh /etc/hegemonie/pki maps \
&& hege-pki-srv.sh /etc/hegemonie/pki regions \
&& hege-pki-srv.sh /etc/hegemonie/pki events



#------------------------------------------------------------------------------
# Create the minimal runtime image

FROM scratch as runtime

USER 0
COPY --chown=0:0 --from=builder /dist /

EXPOSE 6000

USER 0
WORKDIR /
ENTRYPOINT ["/hege"]



#------------------------------------------------------------------------------
# Create a bundle with minimal runtime plus demonsration data

FROM runtime as demo

USER 0
COPY --from=builder /etc/hegemonie    /etc/hegemonie
COPY                /docs/definitions /etc/hegemonie/definitions
COPY                /docs/lang        /etc/hegemonie/lang

USER 0
WORKDIR /
ENTRYPOINT ["/hege"]



#------------------------------------------------------------------------------
# Create the debug runtime image
# Bigger than the minimal runtime in the way it provides a complete ubuntu
# distribution with a rich bash environment. At least can you troubleshoot
# Missing files if it happens

FROM ubuntu:20.04 as debug

USER 0
WORKDIR /

RUN apt-get update -y \
&& apt-get install -y --no-install-recommends \
  librocksdb5.17 \
&& apt-get autoremove --purge \
&& apt-get clean \
&& rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/bin/hege    /usr/bin
COPY --from=builder /etc/hegemonie   /etc/hegemonie
COPY                docs/definitions /etc/hegemonie/definitions
COPY                docs/lang        /etc/hegemonie/lang
COPY                bin              /usr/bin

USER 0
WORKDIR /root
ENTRYPOINT ["/bin/bash"]



#------------------------------------------------------------------------------
# Specialize a Prometheus instance for the in-game timeseries that is pre-
# configured to work with the demonstration environment.

FROM prom/prometheus as demo-prometheus
COPY etc/prometheus/prometheus.yml /etc/prometheus/prometheus.yml
RUN cat /etc/prometheus/prometheus.yml


