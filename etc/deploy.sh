#!/usr/bin/env bash
message=$1
echo git commit -a -m "$message"
git commit -a -m "message"
echo 'ssh'
command="git -C /root/starket/stock-simulator-server status"
run="cd /root/starket/stock-simulator-server && echo here && ./etc/build.sh"
ssh -t root@159.89.154.221 "$command"
ssh -t root@159.89.154.221 "$run"
