package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/app"
	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/config"
	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/logger"
	obs "github.com/Nikita213-hub/grpc_protobuf_study/shared/observability"
)

func main() {
	ctx := context.Background()
	bootstrapLog := logger.New("development", nil)

	bootstrapLog.Info("auth_service_starting", "version", "1.0.0")

	cfg := config.NewAuthServiceCfg()
	err := cfg.Configure()
	if err != nil {
		bootstrapLog.ErrorContext(ctx, "config_load_failed", "error", err.Error())
	}

	env := cfg.Env
	telemetry, err := obs.New(ctx, obs.Config{ //aadd error handling
		ServiceName:    "auth-service",
		ServiceVersion: "v0.1.0",
		Environment:    env,
		OtlpEndpoint:   "0.0.0.0:4317",
		SampleRatio:    cfg.OtelCfg.SampleRatio,
	})
	if err != nil {
		bootstrapLog.ErrorContext(ctx, "telemetry_init_failed", "error", err.Error())
	}

	log := logger.New(env, telemetry.LoggerHandler())

	log.InfoContext(ctx, "logger_configured", "environment", env)
	log.DebugContext(ctx, "config_loaded", "grpc_port", cfg.GrpcCfg.Port, "redis_addr", cfg.RedisCfg.Address)

	defer telemetry.Shutdown(ctx)
	app := app.New(cfg.GrpcCfg.Port, &cfg.RedisCfg, log)
	errCh := make(chan error, 1)
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		errCh <- app.GRPCServer.Run()
	}()

	select {
	case err := <-errCh:
		log.ErrorContext(ctx, "grpc_server_startup_failed", "error", err.Error())
	case <-stopCh:
		log.InfoContext(ctx, "shutdown_signal_received", "signal", "SIGTERM/SIGINT")
		app.Stop()
		log.InfoContext(ctx, "auth_service_stopped")
	}

}
