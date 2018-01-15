@echo off
copy stock-simulator-server docker\
xcopy static docker\static /s /e /h /Y
cd docker
docker-compose build
docker-compose up -d server
cd ..
docker container logs -f starket_server
