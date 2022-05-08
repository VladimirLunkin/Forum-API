FROM golang:1.15 AS build

ADD . /app
WORKDIR /app
RUN go build ./cmd/forum/main.go

FROM ubuntu:20.04

WORKDIR /usr/src/app

COPY . .
COPY --from=build /app/main/ .

EXPOSE 5000
CMD ./main
