#!/usr/bin/env bash
message=$1
echo git commit -a -m "$message"
git commit -a -m "message"
echo 'ssh'
command="cd starket/stock-simulator-server && git pull"
run="cd /root/starket/stock-simulator-server && echo here && ./etc/build.sh"
ssh root@159.89.154.221 $command
ssh root@159.89.154.221 $run
