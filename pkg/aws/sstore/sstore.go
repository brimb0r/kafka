package sstore

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

type secretStore struct {
	client      secretsmanageriface.SecretsManagerAPI
	secretCache map[string]cachedValue
	region      string
	environment string
}

type cachedValue struct {
	value string
}

func New(sess *session.Session) *secretStore {
	sstore := &secretStore{
		client:      secretsmanager.New(sess),
		secretCache: make(map[string]cachedValue),
	}
	return sstore
}

func (sstore *secretStore) Get(key string) (interface{}, error) {
	v, ok := sstore.secretCache[key]

	if ok {
		return v.value, nil
	}

	return sstore.retrieveSecretFromAws(key)
}

func (sstore *secretStore) retrieveSecretFromAws(key string) (string, error) {
	resp, err := sstore.client.GetSecretValue(&secretsmanager.GetSecretValueInput{
		SecretId: aws.String(key),
	})

	if err != nil {
		return "", fmt.Errorf("retrieving key %s: [%w]", key, err)
	}

	if resp.SecretString == nil {
		return "", fmt.Errorf("secret value not present for %s", key)
	}

	value := *resp.SecretString
	sstore.secretCache[key] = cachedValue{value: value}

	return value, nil
}
