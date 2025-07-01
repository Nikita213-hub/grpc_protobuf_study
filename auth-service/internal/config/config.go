package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type AuthServiceConfig struct {
	Env      string           `yaml:"env" env-default:"local"`
	StrgUrl  string           `yaml:"storage_url" env-required:"ture"`
	GrpcCfg  GRPCServerConfig `yaml:"grpc"`
	RedisCfg RedisConfig      `yaml:"redis"`
}

type GRPCServerConfig struct {
	Port    string `yaml:"port" env-required:"true"`
	TimeOut int    `yaml:"timeout" env-default:"30sec"`
}

type RedisConfig struct {
	Port         string `yaml:"port" env-default:"6379"`
	UserPassword string `yaml:"user_password" env-required:"true"`
	DB           int    `yaml:"db" env-default:"0"`
}

func NewAuthServiceCfg() *AuthServiceConfig {
	return &AuthServiceConfig{}
}

func (asc *AuthServiceConfig) Configure(cfgFileName string) error {
	//TODO: reduce hardcoded cfg name, use flags instead
	err := getYamlCfg(asc, cfgFileName)
	if err != nil {
		return err
	}
	err = getEnvCfg(asc)
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

func getEnvCfg(src *AuthServiceConfig) error {
	env := src.Env
	if env == "local" {
		return nil
	}
	if env == "dev" {
		return nil
	}
	if env == "prod" {
		return nil
	}
	return nil
}
