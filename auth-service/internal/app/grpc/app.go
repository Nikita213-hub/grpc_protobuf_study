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

func NewApp(port string, auth authgrpc.Auth) *App {
	server := grpc.NewServer()
	authgrpc.Register(server, auth)
	return &App{
		gRPCserver: server,
		port:       port,
	}
}

func (app *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", app.port))
	if err != nil {
		return err
	}
	fmt.Println("Running grpc auth-service server")
	err = app.gRPCserver.Serve(l)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) Stop() {
	app.gRPCserver.GracefulStop()
}
