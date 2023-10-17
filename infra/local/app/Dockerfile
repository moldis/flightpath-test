#syntax=docker/dockerfile:1.4
FROM golang:1.19-alpine AS base

RUN --mount=type=cache,target=/root/.cache \
    apk --update add git ca-certificates gcc musl-dev

WORKDIR /app
COPY --link go.mod go.sum ./
COPY --link cmd cmd
COPY --link pkg pkg
COPY --link vendor vendor

ARG CGO_ENABLED=0
ARG GOOS=linux

RUN --mount=type=cache,target=/root/.cache/go-build \
    go vet ./... \
    && go build -o api -v artemb/flights-path/cmd/api

FROM scratch AS api

WORKDIR /app
COPY --from=base --link /app/api ./

ENTRYPOINT ["/app/api"]
