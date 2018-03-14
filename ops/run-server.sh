#!/usr/bin/env bash

cd docker

if [ "$HOSTNAME" = high-in-the-clouds ]; then
    export FILE_SERVE=/root/starket/stock-simulator-server/mockstarket-front-end
else
    export FILE_SERVE=$(cd .. && pwd cd docker)/debug_frontend
fi
echo $FILE_SERVE
docker-compose rm -f server
docker-compose build
docker-compose up -d server