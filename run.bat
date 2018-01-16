@echo off
cd docker
set STATIC_FOLDER=/c/Users/Willi/OneDrive/workspace/go/src/github.com/stock-simulator-server/static
docker-compose build
docker-compose run  -d server
cd ..
docker container logs -f starket_server
