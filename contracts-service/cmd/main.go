package main

import (
	"fmt"
	"net"

	"github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/config"
	"github.com/Nikita213-hub/grpc_protobuf_study/contracts-service/internal/app"
)

func main() {
	cfg := config.NewContractsServiceCfg()
	cfg.Configure()
	fmt.Println(cfg)
	app := app.NewApp(cfg)
	list, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%s", cfg.GrpcSrvCfg.Port))
	if err != nil {
		fmt.Println(err)
	}
	app.GRPCServer.Serve(list)
}
