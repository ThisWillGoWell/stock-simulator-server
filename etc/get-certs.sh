#!/usr/bin/env bash
docker-compose run --rm letsencrypt \
  letsencrypt certonly --webroot \
  --email contact@mockstarket.com --agree-tos \
  -w /var/www/letsencrypt -d mockstarket.com
docker-compose kill -s SIGHUP proxy
