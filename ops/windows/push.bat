cd docker
docker-compose build
docker tag docker_server 159.203.244.103:5000/starket_server
docker push --insecure-registry 159.203.244.103:5000/starket_server
