package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AuthServiceConfig struct {
	Env      string
	StrgUrl  string
	GrpcCfg  GRPCServerConfig
	RedisCfg RedisConfig
}

type GRPCServerConfig struct {
	Port    string
	TimeOut int
}

type RedisConfig struct {
	Address      string
	UserPassword string
	DB           int
}

func NewAuthServiceCfg() *AuthServiceConfig {
	return &AuthServiceConfig{}
}

func (cfg *AuthServiceConfig) Configure() error {
	if os.Getenv("ENV") == "local" {
		if err := godotenv.Load(".env"); err != nil {
			fmt.Printf(".env file not found")
		}
	}

	grpcPort := os.Getenv("GRPC_AUTH_SERVICE_PORT")
	if grpcPort == "" {
		return missingEnvError("GRPC_AUTH_SERVICE_PORT")
	}

	timeout, err := parseTimeout(os.Getenv("GRPC_AUTH_SERVICE_TIMEOUT"))
	if err != nil {
		return err
	}
	cfg.GrpcCfg = GRPCServerConfig{
		Port:    grpcPort,
		TimeOut: timeout,
	}

	redisPort := getEnv("REDIS_PORT", "6379")
	redisHost := getEnv("REDIS_HOST", "redis")
	redisPassword := os.Getenv("REDIS_USER_PASSWORD")
	if redisPassword == "" {
		return missingEnvError("REDIS_USER_PASSWORD")
	}

	db, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return invalidEnvError("REDIS_DB", "must be integer")
	}

	cfg.RedisCfg = RedisConfig{
		Address:      redisHost + ":" + redisPort,
		UserPassword: redisPassword,
		DB:           db,
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseTimeout(timeoutStr string) (int, error) {
	if timeoutStr == "" {
		return 30, nil
	}
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return 0, invalidEnvError("GRPC_AUTH_SERVICE_TIMEOUT", "must be integer")
	}
	return timeout, nil
}

func missingEnvError(name string) error {
	return &EnvError{Field: name, Msg: "is required"}
}

func invalidEnvError(name, reason string) error {
	return &EnvError{Field: name, Msg: reason}
}

type EnvError struct {
	Field string
	Msg   string
}

func (e *EnvError) Error() string {
	return "env config error: " + e.Field + " " + e.Msg
}
