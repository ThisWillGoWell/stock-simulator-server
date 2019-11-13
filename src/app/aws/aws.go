package aws

import (
	"encoding/json"
	"fmt"
	"github.com/ThisWillGoWell/stock-simulator-server/src/app/log"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/sirupsen/logrus"
)

type ProgramSecrets struct {
	DiscordToken string `json:"discord_api"`
	RdsPassword  string `json:"rds_password"`
}

type DatabaseSecret struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host string `json:"host"`
	Port int `json:"port"`
}

func GetDatabaseSecret(env string)(DatabaseSecret, error){
	awsSecret, err := getSecret(fmt.Sprintf("mockstarket/%s/database", env))
	var secrets DatabaseSecret
	if err != nil {
		return secrets, fmt.Errorf("failed to get database secret [%v]", err)
	}

	if err := json.Unmarshal([]byte(*awsSecret.SecretString), &secrets); err != nil {
		return secrets, fmt.Errorf("failed to marshal secret [%v]", err)
	}
	return secrets, nil
}

func Secrets(env string) ProgramSecrets {
	awsSecret, err := getSecret(env)
	if err != nil {
		logrus.Errorf("failed to get secret %v", err)
	}
	var secrets ProgramSecrets
	if err := json.Unmarshal(awsSecret.SecretBinary, &secrets); err != nil {
		logrus.Errorf("failed to get marshal secret %v", err)
	}
	return secrets
}

func getSecret(name string) (*secretsmanager.GetSecretValueOutput, error) {

	// This example assumes that you're connecting to ap-southeast-1 region
	// For a full list of endpoints, you can refer to this site -> https://godoc.org/github.com/aws/aws-sdk-go/aws/endpoints
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(endpoints.UsWest2RegionID),
	}))

	svc := secretsmanager.New(sess)
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(name),
	}

	result, err := svc.GetSecretValue(input)
	if err != nil {

			if aerr, ok := err.(awserr.Error); ok {
				switch aerr.Code() {
				case secretsmanager.ErrCodeResourceNotFoundException:
					log.Log.Println(secretsmanager.ErrCodeResourceNotFoundException, aerr.Error())
				case secretsmanager.ErrCodeInvalidParameterException:
					log.Log.Println(secretsmanager.ErrCodeInvalidParameterException, aerr.Error())
				case secretsmanager.ErrCodeInvalidRequestException:
					log.Log.Println(secretsmanager.ErrCodeInvalidRequestException, aerr.Error())
				case secretsmanager.ErrCodeDecryptionFailure:
					log.Log.Println(secretsmanager.ErrCodeDecryptionFailure, aerr.Error())
				case secretsmanager.ErrCodeInternalServiceError:
					log.Log.Println(secretsmanager.ErrCodeInternalServiceError, aerr.Error())
				default:
					log.Log.Println(aerr.Error())
				}
			} else {
				log.Log.Println(err.Error())
			}
			return nil, err

	}

	return result, nil
}
