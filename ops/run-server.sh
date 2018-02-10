#!/usr/bin/env bash
cd docker
export STATIC_FOLDER=/root/starket/mockstarket-front-end
docker-compose build
docker-compose up  -d server