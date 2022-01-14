FROM golang:1.15 AS build

ADD . /app
WORKDIR /app
RUN go build ./cmd/forum/main.go

FROM ubuntu:20.04

RUN apt-get -y update && apt-get install -y tzdata
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo Russia/Moscow > /etc/timezone

RUN apt-get -y update && apt-get install -y postgresql-12
USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER vladimir WITH SUPERUSER PASSWORD 'password';" &&\
    createdb -O vladimir forum &&\
    /etc/init.d/postgresql stop

EXPOSE 5432
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
USER root

WORKDIR /usr/src/app

COPY . .
COPY --from=build /app/main/ .

EXPOSE 5000
ENV PGPASSWORD password
CMD service postgresql start && psql -h localhost -d forum -U vladimir -p 5432 -a -q -f ./db/db.sql && ./main
