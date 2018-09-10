#!/usr/bin/env bash
git pull
./etc/remove-volumes.sh
./etc/run-server.sh
