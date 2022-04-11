# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.17-alpine as build

WORKDIR /src

# Download mod
COPY src/go.mod .
COPY src/go.sum .
RUN go mod download

# Copy source
COPY src .

# Build image
RUN CGO_ENABLED=0 go build -o /speed-test

##
## Generate latest ca-certificates
##

FROM debian:buster-slim AS certs

RUN \
  apt update && \
  apt install -y ca-certificates && \
  cat /etc/ssl/certs/* > /ca-certificates.crt

##
## Deploy
##
FROM scratch

COPY --from=build /speed-test /usr/local/bin/speed-test
COPY --from=certs /ca-certificates.crt /etc/ssl/certs/

ENV HOME /root
ENV USER root

WORKDIR /workspace

ENTRYPOINT [ "speed-test" ]
