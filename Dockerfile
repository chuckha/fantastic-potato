# syntax=docker/dockerfile:1

FROM golang:1.17-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY data ./data
COPY internal ./internal

RUN go build -o /geogame ./cmd/web

EXPOSE 8888

CMD [ "/geogame" ]
