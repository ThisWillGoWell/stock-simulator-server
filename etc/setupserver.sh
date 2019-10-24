sudo yum update -y
sudo yum install -y docker
sudo yum install postgresql postgresql-server postgresql-devel postgresql-contrib
sudo usermod -aG docker ec2-user
sudo service docker start



yum install jq



