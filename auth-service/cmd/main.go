package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/app"
	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/config"
	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/logger"
)

func main() {
	bootstrapLog := logger.New("development")

	bootstrapLog.Info("auth_service_starting", "version", "1.0.0")

	cfg := config.NewAuthServiceCfg()
	err := cfg.Configure()
	if err != nil {
		bootstrapLog.Fatal("config_load_failed", "error", err.Error())
	}

	env := cfg.Env
	log := logger.New(env)

	log.Info("logger_configured", "environment", env)
	log.Debug("config_loaded", "grpc_port", cfg.GrpcCfg.Port, "redis_addr", cfg.RedisCfg.Address)
	app := app.New(cfg.GrpcCfg.Port, &cfg.RedisCfg, log)

	errCh := make(chan error, 1)
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		errCh <- app.GRPCServer.Run()
	}()

	select {
	case err := <-errCh:
		log.Fatal("grpc_server_startup_failed", "error", err.Error())
	case <-stopCh:
		log.Info("shutdown_signal_received", "signal", "SIGTERM/SIGINT")
		app.Stop()
		log.Info("auth_service_stopped")
	}
}
