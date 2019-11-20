DevHost=ec2-34-211-231-2.us-west-2.compute.amazonaws.com
DevDatabase=mockstarket-dev.c6ejpamhqiq5.us-west-2.rds.amazonaws.com
DevFrontendBucket=dev.mockstarket.com

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
#					Database Connection
#################################################################

dev_database:
	etc/database.sh mockstarket/dev/database

dev_database_master:
	etc/database.sh mockstarket/dev/database-master

prod_database: DatabaseHost=${ProdDatabase}

prod_database: connect_database


connect_database:
	etc/secrets.sh prod && ssh -i mockstarket.pem ec2-user@$\${ServerHost} "  psql \"host=${DatabaseHost} port=5432 user=postgres password=$$RDS_PASSWORD dbname=postgres\""


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

# The port tells the frontend what to connect too
frontend_local:
	cd front_end && python ../etc/server.py 8080

frontend_dev:
	cd front_end && python ../etc/server.py 8081

frontend_prod:
	cd front_end && python ../etc/server.py 8082



download_key:
	aws --profile mockstarket s3 cp s3://mockstarket-keys/mockstarket.pem mockstarket.pem




