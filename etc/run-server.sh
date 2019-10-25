#!/bin/bash
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"
cd $GOPATH/src/github.com/ThisWillGoWell/stock-simulator-server
echo "building binary"

echo "building docker container"
docker build . -t mockstarket-$1  --no-cache
echo "save docker container"
docker save mockstarket-$1:latest > mockstarket-$1.tgz

echo "uploading container"

DB_URI=""
SERVER_HOST=""
if [ "$1" == "prod" ]; then
  SERVER_HOST="ec2-34-221-86-219.us-west-2.compute.amazonaws.com"
  DB_URI="mockstarket-prod.c6ejpamhqiq5.us-west-2.rds.amazonaws.com"
elif [ "$1" == "dev" ]; then
  SERVER_HOST="ec2-35-164-117-217.us-west-2.compute.amazonaws.com"
  DB_URI="mockstarket-dev.c6ejpamhqiq5.us-west-2.rds.amazonaws.com"
fi;

scp -i $DIR/../mockstarket.pem mockstarket-$1.tgz ec2-user@$SERVER_HOST:

rm -f mockstarket-$1.tgz
rm stock-simulator-server

echo "getting secrets"
. $DIR/secrets.sh "prod"
echo "running container "
ssh -i $DIR/../mockstarket.pem  ec2-user@$SERVER_HOST "
docker load -i mockstarket-$1.tgz
echo \"stopping any containers\"
docker stop mockstarket-$1
docker rm mockstarket-$1
echo \"strting new container\"
docker run -p 8000:8000 --name mockstarket-$1 \\
-e CONFIG_FOLDER=opt/server/config/ \\
-e SEED_JSON=seed_prod.json \\
-e ITEMS_JSON=items.json \\
-e LEVELS_JSON=levels.json \\
-e DB_URI=\"host=$DB_URI port=5432 user=postgres password=$RDS_PASSWORD dbname=postgres\" \\
-e DISCORD_TOKEN=$DISCORD_TOKEN \\
mockstarket-$1:latest"
