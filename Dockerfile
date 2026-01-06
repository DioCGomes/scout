FROM golang:1.23-alpine AS builder

ARG VERSION=dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN apk --no-cache add build-base

ENV CGO_ENABLED=1
WORKDIR /app/cmd/scout

RUN go build -ldflags="-s -w -X main.version=${VERSION}" -o scout

FROM alpine:3.21

RUN apk --no-cache add libgcc

WORKDIR /scan

COPY --from=builder /app/cmd/scout/scout /usr/local/bin/scout

RUN chmod +x /usr/local/bin/scout

ENTRYPOINT ["scout"]