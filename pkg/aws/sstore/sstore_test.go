package sstore

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/stretchr/testify/assert"
)

type mockSecretsManager struct {
	secretsmanageriface.SecretsManagerAPI
	getSecretCounter int
	secretString     string
}

func (sm *mockSecretsManager) GetSecretValue(*secretsmanager.GetSecretValueInput) (*secretsmanager.GetSecretValueOutput, error) {
	sm.getSecretCounter++
	return &secretsmanager.GetSecretValueOutput{SecretString: aws.String(sm.secretString)}, nil
}

func (sm *mockSecretsManager) setSecretString(s string) *mockSecretsManager {
	sm.secretString = s
	return sm
}

func TestRetrieveSecret(t *testing.T) {
	client := &mockSecretsManager{}
	sstore := &secretStore{
		client:      client,
		secretCache: make(map[string]cachedValue),
		environment: "local",
		region:      "us-east-1",
	}
	client.setSecretString("secretValue")
	value, _ := sstore.retrieveSecretFromAws("foo")
	assert.Equal(t, "secretValue", value)
}

func TestSecretManagerCache(t *testing.T) {
	client := &mockSecretsManager{}
	sstore := &secretStore{
		client:      client,
		secretCache: make(map[string]cachedValue),
		environment: "local",
		region:      "us-east-1",
	}

	client.setSecretString("secretString")

	assert.Equal(t, client.getSecretCounter, 0)
	value, err := sstore.Get("key")
	assert.Equal(t, "secretString", value)
	assert.NoError(t, err)

	assert.Equal(t, client.getSecretCounter, 1)
	_, err = sstore.Get("key")
	assert.NoError(t, err)
	assert.Equal(t, client.getSecretCounter, 1)
}
