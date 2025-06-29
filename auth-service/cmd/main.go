package main

import (
	"fmt"

	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/app"
	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/config"
)

func main() {
	cfg := config.NewAuthServiceCfg()
	err := cfg.Configure()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
	app := app.New(cfg.Port)
	err = app.GRPCServer.Run()
	if err != nil {
		panic(err)
	}
}
