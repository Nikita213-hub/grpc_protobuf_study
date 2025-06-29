package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AuthServiceConfig struct {
	Env              string `yaml:"env" env-default:"local"`
	StrgUrl          string `yaml:"storage_url" env-requird:"ture"`
	GRPCServerConfig `yaml:"grpc"`
}

type GRPCServerConfig struct {
	Port    string `yaml:"port" env-required:"true"`
	TimeOut int    `yaml:"timeout" env-default:"30sec"`
}

func NewAuthServiceCfg() *AuthServiceConfig {
	return &AuthServiceConfig{}
}

func (asc *AuthServiceConfig) Configure() error {
	//TODO: reduce hardcoded cfg name, use flags instead
	err := getYamlCfg(asc, "local.yaml")
	if err != nil {
		return err
	}
	return nil
}

func getYamlCfg(src *AuthServiceConfig, cfgFileName string) error {
	data, err := os.ReadFile(cfgFileName)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(data, &src)
	if err != nil {
		return err
	}
	return nil
}
