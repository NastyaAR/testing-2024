#!/bin/bash

export POSTGRES_USER=usr POSTGRES_DB=postgres POSTGRES_PASSWORD=pass PGPASSFILE=creds.pgpass PASS=password
echo $PASS | sudo -S rm -rf ../postgres/postgres-data
docker compose up -d postgres postgres-exporter

sleep 5

docker exec -it postgres-db apt update

sleep 10

docker exec -it postgres-db apt-get install -y postgresql-plperl-12 
docker exec -it postgres-db apt-get install -y procps 

psql -h 0.0.0.0 -p 5432 -U usr -d postgres -f init.sql -w
