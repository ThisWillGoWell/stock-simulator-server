#!/usr/bin/env bash
echo git commit -a -m '"'$1'"'
git commit -a -m "message"
ssh root@159.89.154.221 "cd starket/stock-simulator-server && git status && ~/starket/stock-simulator-server/etc/run-server.sh"
