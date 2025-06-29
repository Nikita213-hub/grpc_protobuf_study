package grpcapp

import (
	"fmt"
	"net"

	authgrpc "github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/grpc/auth"
	"google.golang.org/grpc"
)

type App struct {
	gRPCserver *grpc.Server
	port       string
}

func NewApp(port string) *App {
	server := grpc.NewServer()
	authgrpc.Register(server)
	return &App{
		gRPCserver: server,
		port:       port,
	}
}

func (app *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", app.port))
	if err != nil {
		return err
	}
	err = app.gRPCserver.Serve(l)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) Stop() {
	app.gRPCserver.GracefulStop()
}
