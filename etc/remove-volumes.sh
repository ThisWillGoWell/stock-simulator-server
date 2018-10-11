#!/usr/bin/env bash
docker-compose stop db server
docker-compose rm -f db
docker volume rm -f stocksimulatorserver_db_volume stocksimulatorserver_ts_volume
docker-compose up -d db
