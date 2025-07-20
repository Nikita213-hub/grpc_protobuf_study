package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/config"
	kafkaContracts "github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/infra/kafka"
	grpcserver "github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/internal/grpc"
	"github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/internal/repository"
	"github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/internal/usecase"

	_ "github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
)

type App struct {
	GRPCServer *grpc.Server
}

func NewApp(cfg *config.ContractsServiceCfg) *App {
	grpcServer := grpc.NewServer()
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DbCfg.DbUser, cfg.DbCfg.DbPassword, cfg.DbCfg.DbHost, cfg.DbCfg.DbPort, cfg.DbCfg.DbName)
	fmt.Println(connStr)
	db, err := sql.Open("pgx", connStr)
	if err != nil { //TODO: handle error
		panic(err)
	}
	if err := db.PingContext(context.Background()); err != nil {
		panic(err) //TODO: handle error
	}
	fmt.Println("database is reachable")
	publisher := kafkaContracts.NewKafkaEventPublisher([]string{cfg.KafkaCfg.Brokers}, cfg.KafkaCfg.Topic)
	contractsRepo := repository.NewPostgreRepository(db)
	contractsUsecase := usecase.NewContractsUsecase(contractsRepo, publisher)
	grpcserver.Register(grpcServer, contractsUsecase)
	return &App{
		GRPCServer: grpcServer,
	}
}
