FROM golang:1.17-alpine as builder

ENV SERVICE_BIN /opt/service
ENV MIGRATE_BIN /opt/migrate
ENV CLI_BIN /opt/cli

ENV CODE_DIR /go/src/

ENV GOOS=linux
ENV CGO_ENABLED=0

WORKDIR ${CODE_DIR}

COPY go.mod go.sum ./

RUN go mod download

COPY ./ ./

ARG LDFLAGS
RUN go build -ldflags "$LDFLAGS" -o ${SERVICE_BIN} cmd/service/*
RUN go build -ldflags "$LDFLAGS" -o ${MIGRATE_BIN} cmd/migrate/*
RUN go build -ldflags "$LDFLAGS" -o ${CLI_BIN} cmd/cli/*

FROM alpine:3.15

ENV BIN_DIR /opt/
COPY --from=builder ${BIN_DIR} ${BIN_DIR}

COPY migrations /migrations
COPY .env /
