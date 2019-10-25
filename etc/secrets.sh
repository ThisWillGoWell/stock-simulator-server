#!/bin/bash
cliString=""
if [ "$(whoami)" != "ec2-user" ]; then
  cliString="--profile mockstarket --region us-west-2"
fi

export DISCORD_API_KEY=$(aws $cliString secretsmanager get-secret-value --secret-id mockstarket/$1 | jq --raw-output '.SecretString' | jq -r .discord_api)
export RDS_PASSWORD=$(aws $cliString secretsmanager get-secret-value --secret-id mockstarket/$1 | jq --raw-output '.SecretString' | jq -r .rds_password)

