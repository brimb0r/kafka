package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"os"
)

var cachedSession *session.Session

func Session() *session.Session {
	var c aws.Config

	if localStack() {
		c = aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
			Credentials:                   credentials.NewStaticCredentials("not", "empty", ""),
			DisableSSL:                    aws.Bool(true),
			Endpoint:                      aws.String("http://localstack:4566"),
			S3ForcePathStyle:              aws.Bool(true),
		}
	} else {
		c = aws.Config{
			CredentialsChainVerboseErrors: aws.Bool(true),
		}
	}

	if cachedSession == nil {
		var err error
		cachedSession, err = session.NewSession(&c)
		if err != nil {
			log.Fatalf("no aws session %s", err.Error())
		}
	}

	return cachedSession.Copy(&c)
}

func localStack() bool {
	localstack := os.Getenv("LOCALSTACK")
	if localstack == "true" {
		return true
	}
	return false
}
