package configuration

import (
	"awesomeProject/pkg/aws"
	"awesomeProject/pkg/aws/sstore"
	"awesomeProject/pkg/config"
	"awesomeProject/pkg/mongo"
)

type Configuration struct {
	Environment string
	Region      string `yaml:"awsRegion"`
	Schedule    string
	ServiceName string `yaml:"serviceName"`
	Mongo       mongo.Config
}

func Configure() *Configuration {
	configure := config.New()
	configure.SetAccessor("SECRET", sstore.New(aws.Session()))
	configuration := &Configuration{}
	configure.Unmarshal(configuration)
	return configuration
}
