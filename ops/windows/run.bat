@echo off
cd docker
set STATIC_FOLDER=/c/Users/Willi/OneDrive/workspace/go/src/github.com/stock-simulator-server/mockstarket-front-end
set DB_FOLDER=/c/Users/Willi/OneDrive/workspace/go/src/github.com/stock-simulator-server/db
docker-compose build
docker-compose run  -d server
cd ..
docker container logs -f docker_server_run_1

