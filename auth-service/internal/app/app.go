package app

import (
	grpcapp "github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/app/grpc"
	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/services/auth"
	redisstore "github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/storage/redisStore"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(port string) *App {
	strg, err := redisstore.NewRedisTokenStorage("localhost:6379", "", 0)
	if err != nil {
		panic(err)
	}
	auth := auth.NewAuth(strg, strg)
	grpcapp := grpcapp.NewApp(port, auth)
	return &App{
		GRPCServer: grpcapp,
	}
}
