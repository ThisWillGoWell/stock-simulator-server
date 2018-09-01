#!/usr/bin/env bash
docker-compose stop
docker-compose rm -f
docker volume rm -f stocksimulatorserver_db_volume stocksimulatorserver_ts_volume
docker-compose up -d