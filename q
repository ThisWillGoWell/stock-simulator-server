AWS()                                                                    AWS()



[1mNAME[0m
       aws -

[1mDESCRIPTION[0m
       The  AWS  Command  Line  Interface is a unified tool to manage your AWS
       services.

[1mSYNOPSIS[0m
          aws [options] <command> <subcommand> [parameters]

       Use [4maws[24m [4mcommand[24m [4mhelp[24m for information on a  specific  command.  Use  [4maws[0m
       [4mhelp[24m  [4mtopics[24m  to view a list of available help topics. The synopsis for
       each command shows its parameters and their usage. Optional  parameters
       are shown in square brackets.

[1mOPTIONS[0m
       [1m--debug [22m(boolean)

       Turn on debug logging.

       [1m--endpoint-url [22m(string)

       Override command's default URL with the given URL.

       [1m--no-verify-ssl [22m(boolean)

       By  default, the AWS CLI uses SSL when communicating with AWS services.
       For each SSL connection, the AWS CLI will verify SSL certificates. This
       option overrides the default behavior of verifying SSL certificates.

       [1m--no-paginate [22m(boolean)

       Disable automatic pagination.

       [1m--output [22m(string)

       The formatting style for command output.

       +o json

       +o text

       +o table

       [1m--query [22m(string)

       A JMESPath query to use in filtering the response data.

       [1m--profile [22m(string)

       Use a specific profile from your credential file.

       [1m--region [22m(string)

       The region to use. Overrides config/env settings.

       [1m--version [22m(string)

       Display the version of this tool.

       [1m--color [22m(string)

       Turn on/off color output.

       +o on

       +o off

       +o auto

       [1m--no-sign-request [22m(boolean)

       Do  not  sign requests. Credentials will not be loaded if this argument
       is provided.

       [1m--ca-bundle [22m(string)

       The CA certificate bundle to use when verifying SSL certificates. Over-
       rides config/env settings.

       [1m--cli-read-timeout [22m(int)

       The  maximum socket read time in seconds. If the value is set to 0, the
       socket read will be blocking and not timeout.

       [1m--cli-connect-timeout [22m(int)

       The maximum socket connect time in seconds. If the value is set  to  0,
       the socket connect will be blocking and not timeout.

[1mAVAILABLE SERVICES[0m
       +o acm

       +o acm-pca

       +o alexaforbusiness

       +o amplify

       +o apigateway

       +o apigatewaymanagementapi

       +o apigatewayv2

       +o application-autoscaling

       +o application-insights

       +o appmesh

       +o appstream

       +o appsync

       +o athena

       +o autoscaling

       +o autoscaling-plans

       +o backup

       +o batch

       +o budgets

       +o ce

       +o chime

       +o cloud9

       +o clouddirectory

       +o cloudformation

       +o cloudfront

       +o cloudhsm

       +o cloudhsmv2

       +o cloudsearch

       +o cloudsearchdomain

       +o cloudtrail

       +o cloudwatch

       +o codebuild

       +o codecommit

       +o codepipeline

       +o codestar

       +o cognito-identity

       +o cognito-idp

       +o cognito-sync

       +o comprehend

       +o comprehendmedical

       +o configservice

       +o configure

       +o connect

       +o cur

       +o datapipeline

       +o datasync

       +o dax

       +o deploy

       +o devicefarm

       +o directconnect

       +o discovery

       +o dlm

       +o dms

       +o docdb

       +o ds

       +o dynamodb

       +o dynamodbstreams

       +o ec2

       +o ec2-instance-connect

       +o ecr

       +o ecs

       +o efs

       +o eks

       +o elasticache

       +o elasticbeanstalk

       +o elastictranscoder

       +o elb

       +o elbv2

       +o emr

       +o es

       +o events

       +o firehose

       +o fms

       +o forecast

       +o forecastquery

       +o fsx

       +o gamelift

       +o glacier

       +o globalaccelerator

       +o glue

       +o greengrass

       +o groundstation

       +o guardduty

       +o health

       +o help

       +o history

       +o iam

       +o importexport

       +o inspector

       +o iot

       +o iot-data

       +o iot-jobs-data

       +o iot1click-devices

       +o iot1click-projects

       +o iotanalytics

       +o iotevents

       +o iotevents-data

       +o iotthingsgraph

       +o kafka

       +o kinesis

       +o kinesis-video-archived-media

       +o kinesis-video-media

       +o kinesisanalytics

       +o kinesisanalyticsv2

       +o kinesisvideo

       +o kms

       +o lakeformation

       +o lambda

       +o lex-models

       +o lex-runtime

       +o license-manager

       +o lightsail

       +o logs

       +o machinelearning

       +o macie

       +o managedblockchain

       +o marketplace-entitlement

       +o marketplacecommerceanalytics

       +o mediaconnect

       +o mediaconvert

       +o medialive

       +o mediapackage

       +o mediapackage-vod

       +o mediastore

       +o mediastore-data

       +o mediatailor

       +o meteringmarketplace

       +o mgh

       +o mobile

       +o mq

       +o mturk

       +o neptune

       +o opsworks

       +o opsworks-cm

       +o organizations

       +o personalize

       +o personalize-events

       +o personalize-runtime

       +o pi

       +o pinpoint

       +o pinpoint-email

       +o pinpoint-sms-voice

       +o polly

       +o pricing

       +o qldb

       +o qldb-session

       +o quicksight

       +o ram

       +o rds

       +o rds-data

       +o redshift

       +o rekognition

       +o resource-groups

       +o resourcegroupstaggingapi

       +o robomaker

       +o route53

       +o route53domains

       +o route53resolver

       +o s3

       +o s3api

       +o s3control

       +o sagemaker

       +o sagemaker-runtime

       +o sdb

       +o secretsmanager

       +o securityhub

       +o serverlessrepo

       +o service-quotas

       +o servicecatalog

       +o servicediscovery

       +o ses

       +o shield

       +o signer

       +o sms

       +o snowball

       +o sns

       +o sqs

       +o ssm

       +o stepfunctions

       +o storagegateway

       +o sts

       +o support

       +o swf

       +o textract

       +o transcribe

       +o transfer

       +o translate

       +o waf

       +o waf-regional

       +o workdocs

       +o worklink

       +o workmail

       +o workmailmessageflow

       +o workspaces

       +o xray

[1mSEE ALSO[0m
       +o aws help topics



                                                                         AWS()
