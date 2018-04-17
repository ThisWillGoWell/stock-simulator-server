#!/usr/bin/env bash

cd docker
export DB_FOLDER=$(pwd)
export TS_FOLDER=$(pwd)
pwd=$(pwd)
cd ${pwd:0:$(echo `expr "$pwd" :  '.*stock-simulator-server'`)}

export FILE_SERVE=$(pwd)/debug_frontend

#if [ "$HOSTNAME" = high-in-the-clouds ]; then

#else
#
#fi

cd docker
echo $FILE_SERVE
docker-compose rm -f server
docker-compose build
docker-compose up -d
docker-compose logs -f
