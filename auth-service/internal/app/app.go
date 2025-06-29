package app

import grpcapp "github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/app/grpc"

type App struct {
	GRPCServer *grpcapp.App
}

func New(port string) *App {
	grpcapp := grpcapp.NewApp(port)
	return &App{
		GRPCServer: grpcapp,
	}
}
