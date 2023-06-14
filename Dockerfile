FROM golang:1.20-alpine AS builder

RUN apk add bash ca-certificates
WORKDIR /go/src/github.com/openware/wio

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go install .

FROM alpine

RUN apk add ca-certificates
COPY --from=builder /go/bin/wio /bin/wio

ENTRYPOINT ["/bin/wio"]
