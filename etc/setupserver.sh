#!/bin/bash
yum -y update
sudo yum update -y
sudo yum install -y docker
sudo yum install postgresql postgresql-server postgresql-devel postgresql-contrib
sudo usermod -aG docker ec2-user
sudo service docker start
sudo yum install jq
# install aws cli
sudo yum install -y ruby
cd /home/ec2-user
curl -O https://aws-codedeploy-us-west-2.s3.amazonaws.com/latest/install
chmod +x ./install
sudo ./install auto







