package app

import (
	"fmt"

	grpcapp "github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/app/grpc"
	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/config"
	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/services/auth"
	redisstore "github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/storage/redisStore"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(port string, redisCfg *config.RedisConfig) *App {
	strg, err := redisstore.NewRedisTokenStorage(redisCfg.Address, redisCfg.UserPassword, redisCfg.DB)
	if err != nil {
		panic(err)
	}
	fmt.Println("Redis was connected successfully")
	auth := auth.NewAuth(strg, strg)
	grpcapp := grpcapp.NewApp(port, auth)
	return &App{
		GRPCServer: grpcapp,
	}
}

func (app *App) Stop() {
	fmt.Println("Treminating service...")
	app.GRPCServer.Stop()
}
