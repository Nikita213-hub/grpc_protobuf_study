package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/app"
	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/config"
)

func main() {
	fmt.Println("HELLO YOOO")
	cfg := config.NewAuthServiceCfg()
	err := cfg.Configure("cmd/local.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
	app := app.New(cfg.GrpcCfg.Port, &cfg.RedisCfg)
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGTERM, syscall.SIGINT)
	go app.GRPCServer.Run()
	<-stopCh
	app.Stop()
}
