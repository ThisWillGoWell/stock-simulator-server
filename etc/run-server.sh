#!/usr/bin/env bash

cd docker

if [ "$HOSTNAME" = high-in-the-clouds ]; then
    export FILE_SERVE=/root/starket/stock-simulator-server/debug_frontend
else
    export FILE_SERVE=$(pwd cd docker)/debug_frontend
fi
echo $FILE_SERVE
docker-compose rm -f server
docker-compose build
docker-compose up -d server
docker-compose logs -f
