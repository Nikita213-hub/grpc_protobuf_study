package app

import (
	grpcapp "github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/app/grpc"
	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/config"
	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/logger"
	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/services/auth"
	redisstore "github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/storage/redisStore"
)

type App struct {
	GRPCServer *grpcapp.App
	logger     *logger.Logger
}

func New(port string, redisCfg *config.RedisConfig, log *logger.Logger) *App {
	strg, err := redisstore.NewRedisTokenStorage(redisCfg.Address, redisCfg.UserPassword, redisCfg.DB)
	if err != nil {
		log.Fatal("redis_connection_failed", "error", err.Error())
	}
	log.Info("redis_connected", "addr", redisCfg.Address, "db", redisCfg.DB)

	auth := auth.NewAuth(strg, strg)
	grpcapp := grpcapp.NewApp(port, auth, log)
	return &App{
		GRPCServer: grpcapp,
		logger:     log,
	}
}

func (app *App) Stop() {
	app.logger.Info("service_stopping")
	app.GRPCServer.Stop()
}
