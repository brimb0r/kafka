package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	ConnectTimeout int      `yaml:"connectTimeout"`
	DbName         string   `yaml:"dbName"`
	Hosts          []string `yaml:"hosts"`
	Password       string   `yaml:"password"`
	PoolSize       int      `yaml:"poolSize"`
	ReadConcern    string   `yaml:"readConcern"`
	ReadPreference string   `yaml:"readPreference"`
	SocketTimeout  int      `yaml:"socketTimeout"`
	Uri            string   `yaml:"uri"`
	Username       string   `yaml:"username"`
	WriteConcern   string   `yaml:"writeConcern"`
}

func (mConfig *Config) MongoDatabase() *mongo.Database {
	clientOptions := mConfig.toClientOptions()
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(fmt.Errorf("mongo failure: %w", err))
	}
	return client.Database(mConfig.DbName)
}

func (mConfig *Config) toClientOptions() *options.ClientOptions {
	clientOptions := options.Client()

	if len(mConfig.Hosts) > 0 {
		clientOptions.SetHosts(mConfig.Hosts)
	}

	if mConfig.PoolSize > 0 {
		clientOptions.SetMaxPoolSize(uint64(mConfig.PoolSize))
	}

	if mConfig.ReadConcern != "" {
		clientOptions.SetReadConcern(readconcern.New(readconcern.Level(mConfig.ReadConcern)))
	}

	if mConfig.ReadPreference != "" {
		switch strings.ToLower(mConfig.ReadPreference) {
		case "primary":
			clientOptions.SetReadPreference(readpref.Primary())
		case "secondary":
			clientOptions.SetReadPreference(readpref.Secondary())
		case "nearest":
			clientOptions.SetReadPreference(readpref.Nearest())
		}
	}

	if mConfig.ConnectTimeout > 0 {
		clientOptions.SetConnectTimeout(time.Second * time.Duration(mConfig.ConnectTimeout))
	}

	if mConfig.SocketTimeout > 0 {
		clientOptions.SetSocketTimeout(time.Second * time.Duration(mConfig.ConnectTimeout))
	}

	if mConfig.WriteConcern != "" {
		if strings.ToLower(mConfig.WriteConcern) == "majority" {
			clientOptions.SetWriteConcern(writeconcern.New(writeconcern.WMajority()))
		} else if number, err := strconv.Atoi(mConfig.WriteConcern); err != nil {
			clientOptions.SetWriteConcern(writeconcern.New(writeconcern.W(number)))
		}
	}

	if mConfig.Username != "" {
		clientOptions.SetAuth(options.Credential{
			Username: mConfig.Username,
			Password: mConfig.Password,
		})
	}

	if mConfig.Uri != "" {
		clientOptions.ApplyURI(mConfig.Uri)
	}

	return clientOptions
}
