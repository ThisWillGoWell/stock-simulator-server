@echo off
cd docker
docker-compose build
docker-compose up -d server
cd ..
docker container logs -f starket_server
