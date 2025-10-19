# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git gcc musl-dev openssh-client

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o netconf-checker \
    ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates openssh-client

WORKDIR /app

COPY --from=builder /build/netconf-checker /app/netconf-checker

RUN chmod +x /app/netconf-checker

USER nobody:nobody

ENTRYPOINT ["/app/netconf-checker"]
