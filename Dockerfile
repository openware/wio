FROM golang:1.15-alpine AS builder

RUN apk add bash ca-certificates
WORKDIR /go/src/github.com/openware/wio

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go install .

USER root
RUN apk add ca-certificates

ENTRYPOINT ["/bin/wio"]
