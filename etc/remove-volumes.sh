#!/usr/bin/env bash
docker-compose stop ts db server
docker-compose rm -f ts db
docker volume rm -f stocksimulatorserver_db_volume stocksimulatorserver_ts_volume
docker-compose up -d ts db
