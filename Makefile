DevHost=ec2-35-164-117-217.us-west-2.compute.amazonaws.com
DevDatabase=mockstarket-dev.c6ejpamhqiq5.us-west-2.rds.amazonaws.com

ProdHost=ec2-34-221-86-219.us-west-2.compute.amazonaws.com
ProdDatabase=mockstarket-prod.c6ejpamhqiq5.us-west-2.rds.amazonaws.com


#################################################################
#					Connect
#################################################################


connect_prod: ServerHost=${ProdHost}
connect_prod: connect

connect_dev: ServerHost=${DevHost}
connect_dev: connect

connect:
	ssh -i mockstarket.pem ec2-user@${ServerHost}



#################################################################
#					Build and deploy Backend
#################################################################


AwsProfile := mockstarket

RootPath=${GOPATH}/src/github.com/ThisWillGoWell/stock-simulator-server

deploy_prod: DockerTag=mockstarket-prod
deploy_prod: ServerHost=${ProdHost}
deploy_prod: DatabaseHost=${ProdDatabase}
deploy_prod: | build_linux build_container save_container upload_container run_container

deploy_dev: ServerHost=${DevHost}
deploy_dev: DatabaseHost=${DevDatabase}
deploy_dev: DockerTag=mockstarket-dev
deploy_dev:  | build_linux build_container save_container upload_container run_container

build_linux:
	GOARCH=amd64 GOOS=linux go build
build_container:
	docker build . -t ${DockerTag} --no-cache

save_container:
	docker save ${DockerTag} > ${DockerTag}.tgz

upload_container:
	scp -i mockstarket.pem ${DockerTag}.tgz ec2-user@${ServerHost}:

run_container:
	. etc/secrets.sh prod && ssh -i mockstarket.pem  ec2-user@${ServerHost} " \
		docker load -i ${DockerTag}.tgz; \
		docker stop ${DockerTag}; \
		docker rm ${DockerTag}; \
		docker run -p 8000:8000 --name ${DockerTag} \
		-e CONFIG_FOLDER=opt/server/config/ \
		-e SEED_JSON=seed_prod.json \
		-e ITEMS_JSON=items.json \
		-e LEVELS_JSON=levels.json \
		-e DB_URI=\"host=${DatabaseHost} port=5432 user=postgres password=$$RDS_PASSWORD dbname=postgres\" \
		-e DISCORD_TOKEN=$$DISCORD_API_KEY \
		${DockerTag}:latest"

#################################################################
#					Database Connection
#################################################################

dev_database: DatabaseHost=${DevDatabase}
dev_database: ServerHost=${DevHost}
dev_database: connect_database

prod_database: DatabaseHost=${ProdDatabase}
prod_database: ServerHost=${ProdHost}
prod_database: connect_database

connect_database:
	. etc/secrets.sh prod && ssh -i mockstarket.pem ec2-user@$\${ServerHost} "  psql \"host=${DatabaseHost} port=5432 user=postgres password=$$RDS_PASSWORD dbname=postgres\""


#################################################################
#					Build and deploy Frontend
#################################################################



#################################################################
#					Other things
#################################################################

download_key:
	@$(AWS_PROFILE=${AWS_PROFILE} aws s3 cp s3://mockstarket-keys/mockstarket.pem mockstarket.pem)



