package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type configurator struct {
	Env         string
	Region      string
	ConfigsPath string
	accessors   map[string]*Accessor
}

func New() *configurator {
	cfg := &configurator{}
	cfg.initEnv()
	cfg.initRegion()
	cfg.initConfigPath()
	return cfg
}

func (config *configurator) Unmarshal(configObj interface{}) error {
	path, err := filepath.Abs(config.ConfigsPath)
	if err != nil {
		return fmt.Errorf("file path %s: %w", config.ConfigsPath, err)
	}
	configurationFileName := fmt.Sprintf("%s/%s-%s.yml", path, config.Env, config.Region)
	yamlConfig, err := os.ReadFile(configurationFileName)
	if err != nil {
		return fmt.Errorf("reading configuration file %s %w", configurationFileName, err)
	}

	err = config.yamlConfig(yamlConfig, configObj)
	if err != nil {
		return fmt.Errorf("configuration initialization [env:%s] [region:%s] %w", config.Env, config.Region, err)
	}
	return nil
}

func (config *configurator) SetAccessor(id string, accessor Accessor) {
	config.setAccessor(id, accessor)
}

func (config *configurator) initEnv() {
	env := os.Getenv("ENVIRONMENT")
	if env != "" {
		config.Env = env
	} else if config.Env == "" {
		config.Env = "local"
	}
}

func (config *configurator) initRegion() {
	region := os.Getenv("AWS_REGION")
	if region != "" {
		config.Region = region
	} else if config.Region == "" {
		config.Region = "us-east-1" // Default Region
	}
}

func (config *configurator) initConfigPath() {
	config.ConfigsPath = os.Getenv("CONFIG_PATH")
}
