DevHost=ec2-34-211-231-2.us-west-2.compute.amazonaws.com
DevDatabase=mockstarket-dev.c6ejpamhqiq5.us-west-2.rds.amazonaws.com
DevFrontendBucket=mockstarket-frontend-dev

ProdHost=ec2-34-221-86-219.us-west-2.compute.amazonaws.com
ProdDatabase=mockstarket-prod.c6ejpamhqiq5.us-west-2.rds.amazonaws.com
ProdFrontendBucket=mockstarket-frontend
ProdDockerTag=mockstarket-prod

DevEc2Instance=i-047084609e35c4df6
DevRdsInstance=mockstarket-dev
DevDockerTag=mockstarket-dev


#################################################################
#					Connect
#################################################################

convert_lines:
	find . -name "*" | xargs dos2unix -q


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

deploy_prod: DockerTag=${ProdDockerTag}
deploy_prod: ServerHost=${ProdHost}
deploy_prod: DatabaseHost=${ProdDatabase}
deploy_prod: | build_linux build_container save_container upload_container run_container

deploy_dev: ServerHost=${DevHost}
deploy_dev: DatabaseHost=${DevDatabase}
deploy_dev: DockerTag=${DevDockerTag}
deploy_dev:  | build_linux build_container save_container upload_container run_container

build_linux:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build

build_container:
	docker build . -t ${DockerTag} --no-cache

save_container:
	docker save ${DockerTag} > ${DockerTag}.tgz

upload_container:
	scp -i mockstarket.pem ${DockerTag}.tgz ec2-user@${ServerHost}:

run_container:
	. etc/secrets.sh prod && ssh -i mockstarket.pem  ec2-user@${ServerHost} " \
		sudo service docker start; \
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

deploy_prod_frontend: FrontendBucket=${ProdFrontendBucket}
deploy_prod_frontend: deploy_frontend

deploy_dev_frontend: FrontendBucket=${DevFrontendBucket}
deploy_dev_frontend: deploy_frontend

deploy_frontend:
	cd front_end
	find . -name '.DS_Store' -type f -delete
	aws --profile mockstarket s3 sync ./front_end s3://${FrontendBucket}/ --delete


#################################################################
#					Control Dev State
#################################################################

stop_dev: DbId=${DevRdsInstance}
stop_dev: Ec2Id=${DevEc2Instance}
stop_dev: action=stop
stop_dev: stop_or_start

start_dev: DbId=${DevRdsInstance}
start_dev: Ec2Id=${DevEc2Instance}
start_dev: DockerTag=${DevDockerTag}
start_dev: action=start
start_dev: ServerHost=${DevHost}
start_dev: DatabaseHost=${DevDatabase}
start_dev: stop_or_start wait_for_running run_container

stop_or_start:
	-aws --profile mockstarket --region us-west-2 rds ${action}-db-instance --db-instance-identifier ${DbId}
	-aws --profile mockstarket --region us-west-2 ec2 ${action}-instances --instance-id ${Ec2Id}

wait_for_running:
	aws --profile mockstarket --region us-west-2 ec2 wait instance-running --instance-ids ${Ec2Id}
	aws --profile mockstarket --region us-west-2 rds wait db-instance-available --db-instance-identifier ${DbId}



#################################################################
#					Other things
#################################################################

local_serve:
	cd front_end && python -m SimpleHTTPServer


download_key:
	aws --profile mockstarket s3 cp s3://mockstarket-keys/mockstarket.pem mockstarket.pem




