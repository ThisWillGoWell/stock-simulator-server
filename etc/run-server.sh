#!/usr/bin/env bash
echo hello world
pwd=$(pwd)
export FILE_SERVE=$(pwd)/front_end
#if [ "$HOSTNAME" = high-in-the-clouds ]; then

#else
#
#fi
echo $FILE_SERVE
docker-compose build --no-cache server
docker-compose up -d
docker-compose logs -f
