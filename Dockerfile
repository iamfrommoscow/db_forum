FROM golang as build

WORKDIR /app
COPY . .
RUN go get -d & go build -v

FROM ubuntu:18.04
ENV DEBIAN_FRONTEND=noninteractive 
RUN apt-get update && apt-get install -y postgresql-10

USER postgres

RUN service postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    service postgresql stop

RUN echo "listen_addresses = '*'" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "synchronous_commit = off" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "fsync = off" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "autovacuum = off" >> /etc/postgresql/10/main/postgresql.conf
RUN echo "unix_socket_directories = '/var/run/postgresql'" >> /etc/postgresql/10/main/postgresql.conf

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

COPY createDB.sql .
COPY --from=build /app/app .

CMD service postgresql start && ./app






