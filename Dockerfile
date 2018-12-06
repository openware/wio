FROM golang:1.11-alpine AS builder

# Install some dependencies needed to build the project
RUN apk add bash ca-certificates git gcc g++ libc-dev
WORKDIR /go/src/github.com/openware/wio

ENV GO111MODULE=on

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go install -a -tags netgo -ldflags '-w -extldflags "-static"'

FROM alpine

RUN apk add ca-certificates
COPY --from=builder /go/bin/wio /bin/wio

ENTRYPOINT ["/bin/wio"]
