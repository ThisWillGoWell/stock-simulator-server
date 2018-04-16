#!/bin/bash
go build -o stock-server-linux
cp -f stock-server-linux docker/
cp -f stock-server-linux bin/linux/

GOARCH="amd64" GOOS=darwin go build -o stock-server-osx
cp -f stock-server-osx bin/osx/
rm stock-server-linux
rm stock-server-osx

