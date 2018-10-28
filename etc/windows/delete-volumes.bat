docker-compose stop
docker-compose rm -f
docker volume rm -f stock-simulator-server_db_volume
docker-compose up -d db
