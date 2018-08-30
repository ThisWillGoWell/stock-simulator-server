@echo off
set FILE_SERVE=c:/Users/Willi/Dropbox/workspace/go/src/github.com/stock-simulator-server/debug_frontend
docker-compose up -d
docker-compose logs -f

