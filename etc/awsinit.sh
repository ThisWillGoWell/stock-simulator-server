#!/bin/bash
if [ ! -f "$HOME/.aws/credentials" ]; then
  echo "maing file"
  mkdir -p $HOME/.aws
  touch $HOME/.aws/credentials
fi

echo "[mockstarket]" >> ~/.aws/credentials
echo "aws_access_key_id = $1" >> ~/.aws/credentials
echo "aws_secret_access_key = $2" >> ~/.aws/credentials
