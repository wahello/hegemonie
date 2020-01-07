#https://dev.to/ivan/go-build-a-minimal-docker-image-in-just-three-steps-514i

FROM golang:1.13.0-stretch AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=1

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build

WORKDIR /dist
RUN ls -l && pwd && ls -l /
RUN cp /build/hegemonie ./hegemonie

RUN ldd /build/hegemonie | tr -s '[:blank:]' '\n' | grep '^/' | \
    xargs -I % sh -c 'mkdir -p $(dirname ./%); cp % ./%;'
RUN mkdir -p lib64 && cp /lib64/ld-linux-x86-64.so.2 lib64/

RUN mkdir /data

# Create the minimal runtime image
FROM scratch
COPY --chown=0:0 --from=builder /dist /
COPY --chown=65534:0 --from=builder /data /data
USER 65534
WORKDIR /data
ENTRYPOINT ["/hegemonie"]
