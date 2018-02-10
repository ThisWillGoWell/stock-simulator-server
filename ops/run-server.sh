#!/usr/bin/env bash
cd docker
export STATIC_FOLDER=/root/starket/stock-simulator-server/mockstarket-front-end
docker-compose build
docker-compose up  -d server