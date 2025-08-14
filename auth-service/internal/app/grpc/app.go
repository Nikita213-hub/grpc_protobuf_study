package grpcapp

import (
	"fmt"
	"net"

	authgrpc "github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/grpc/auth"
	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type App struct {
	gRPCserver *grpc.Server
	port       string
	logger     *logger.Logger
}

func NewApp(port string, auth authgrpc.Auth, log *logger.Logger) *App {
	server := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	authgrpc.Register(server, auth, log)
	return &App{
		gRPCserver: server,
		port:       port,
		logger:     log,
	}
}

func (app *App) Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", app.port))
	if err != nil {
		return err
	}
	app.logger.Info("grpc_server_starting", "port", app.port)
	err = app.gRPCserver.Serve(l)
	if err != nil {
		return err
	}
	return nil
}

func (app *App) Stop() {
	app.gRPCserver.GracefulStop()
}
