#!/usr/bin/env bash
echo hello world
pwd=$(pwd)
export FILE_SERVE=$(pwd)/debug_frontend

#if [ "$HOSTNAME" = high-in-the-clouds ]; then

#else
#
#fi
echo $FILE_SERVE
docker-compose build --no-caches
docker-compose up -d
docker-compose logs -f
