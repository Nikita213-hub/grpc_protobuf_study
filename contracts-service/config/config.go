package config

import (
	"os"
	"strconv"
)

var DEFAULT_GRPC_TIME_OUT = 30

type ContractsServiceCfg struct {
	GrpcSrvCfg GrpcServerCfg
	DbCfg      DBCfg
	KafkaCfg   KafkaCfg
}

type GrpcServerCfg struct {
	Port    string
	TimeOut int
}

type DBCfg struct {
	DbHost     string
	DbPort     string
	DbUser     string
	DbName     string
	DbPassword string
}

type KafkaCfg struct {
	Brokers string
	Topic   string
}

func NewContractsServiceCfg() *ContractsServiceCfg {
	return &ContractsServiceCfg{}
}

func (csc *ContractsServiceCfg) Configure() {
	mustLoadCfgFromEnv(csc)
}

func mustLoadCfgFromEnv(cfg *ContractsServiceCfg) {

	cfg.DbCfg.DbHost = os.Getenv("DB_HOST")
	cfg.DbCfg.DbPort = os.Getenv("DB_PORT")
	cfg.DbCfg.DbUser = os.Getenv("DB_USER")
	cfg.DbCfg.DbName = os.Getenv("DB_NAME")
	cfg.DbCfg.DbPassword = os.Getenv("DB_PASSWORD")

	cfg.GrpcSrvCfg.Port = os.Getenv("GRPC_CONTRACTS_SERVICE_PORT")
	grpcTimeOut, err := strconv.Atoi(os.Getenv("GRPC_CONTRACTS_SERVIE_TIME_OUT"))
	if err != nil {
		grpcTimeOut = DEFAULT_GRPC_TIME_OUT
	}
	cfg.GrpcSrvCfg.TimeOut = grpcTimeOut

	cfg.KafkaCfg.Brokers = os.Getenv("KAFKA_BROKERS")
	cfg.KafkaCfg.Topic = os.Getenv("KAFKA_CONTRACS_TOPIC")
}
