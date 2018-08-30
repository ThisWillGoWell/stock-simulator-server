#!/usr/bin/env bash
message=$1
echo git commit -a -m "$message"
git commit -a -m "message"
command="cd starket/stock-simulator-server && git pull && /root/starket/stock-simulator-server/etc/run-server.sh"
ssh root@159.89.154.221 $command
