sudo yum update -y
sudo yum install -y docker
sudo usermod -aG docker ec2-user
sudo service docker start
sudo curl -L "https://github.com/docker/compose/releases/download/1.24.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose
git clone https://github.com/ThisWillGoWell/stock-simulator-server.git
chmod +x stock-simulator-server/ect/run-server.sh
yum install jq



