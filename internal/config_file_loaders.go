package internal

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrNoDefaultConfigFile                   = errors.New("hearbeat.yaml file was not found")
	DefaultConfigFileReader configFileReader = &yamlFileLoader{}
)

type configFileReader interface {
	Read(*Config) error
}

type yamlFileLoader struct{}

func (y *yamlFileLoader) Read(config *Config) error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		if configFile == DefaultConfigFile {
			return ErrNoDefaultConfigFile
		}
		return err
	}
	if err := yaml.Unmarshal(data, config); err != nil {
		return err
	}
	return nil
}
