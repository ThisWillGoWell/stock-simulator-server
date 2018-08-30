#!/usr/bin/env bash
git commit -am '$@'
ssh root@159.89.154.221 "bash -c 'cd starket/stock-simulator-server && git status && ~/starket/stock-simulator-server/etc/run-server.sh'"
