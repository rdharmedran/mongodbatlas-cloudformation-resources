package util

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type DeploymentSecret struct {
	PublicKey  string              `json:"PublicKey"`
	PrivateKey string              `json:"PrivateKey"`
	ResourceID *ResourceIdentifier `json:"ResourceID"`
	Properties *map[string]string  `json:"Properties"`
}

func CreateDeploymentSecret(req *handler.Request, cfnID *ResourceIdentifier, publicKey, privateKey string, properties map[string]string) (*string, error) {
	deploySecret := &DeploymentSecret{
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		ResourceID: cfnID,
		Properties: &properties,
	}
	log.Printf("deploySecret: %v", deploySecret)
	deploySecretString, _ := json.Marshal(deploySecret)
	log.Printf("deploySecretString: %s", deploySecretString)

	log.Println("===============================================")
	log.Printf("%+v", os.Environ())
	log.Println("===============================================")

	// sess := credentials.SessionFromCredentialsProvider(creds)
	// create a new secret from this struct with the json string

	// Create service client value configured for credentials
	// from assumed role.
	svc := secretsmanager.New(req.Session)
	input := &secretsmanager.CreateSecretInput{
		Description:  aws.String("MongoDB Atlas Quickstart Deployment Secret"),
		Name:         aws.String(cfnID.String()),
		SecretString: aws.String(string(deploySecretString)),
	}

	result, err := svc.CreateSecret(input)
	if err != nil {
		// Print the error, cast err to awserr. Error to get the Code and
		// Message from an error.
		log.Printf("error create secret: %+v", err.Error())
		return nil, err
	}
	log.Printf("Created secret result:%+v", result)
	return result.Name, nil
}

func GetAPIKeyFromDeploymentSecret(req *handler.Request, secretName string) (DeploymentSecret, error) {
	fmt.Printf("secretName=%s\n", secretName)
	sm := secretsmanager.New(req.Session)
	output, err := sm.GetSecretValue(&secretsmanager.GetSecretValueInput{SecretId: &secretName})
	if err != nil {
		log.Printf("Error --- %v", err.Error())
		return DeploymentSecret{}, err
	}
	fmt.Println(*output.SecretString)
	var key DeploymentSecret
	err = json.Unmarshal([]byte(*output.SecretString), &key)
	if err != nil {
		log.Printf("Error --- %v", err.Error())
		return key, err
	}
	fmt.Printf("%v", key)
	return key, nil
}
