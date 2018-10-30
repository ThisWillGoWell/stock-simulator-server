@echo off
set FILE_SERVE=%CD%\front_end
set CONFIG_FILE=%CD%\config\config_prod.json
docker-compose build server
docker-compose up -d
docker-compose logs -f

