package auth

import (
	"context"

	authV1 "github.com/Nikita213-hub/grpc_protobuf_study/pkg/proto/auth/v1"
	"google.golang.org/grpc"
)

type Server struct {
	authV1.UnimplementedAuthServiceServer
}

func Register(gRPC *grpc.Server) {
	authV1.RegisterAuthServiceServer(gRPC, &Server{})
}

func (s *Server) GenToken(ctx context.Context, req *authV1.GenTokenReq) (*authV1.GenTokenRes, error) {
	panic("unimplemented")
}

func (s *Server) VerifyToken(ctx context.Context, req *authV1.VerifyTokenReq) (*authV1.VerifyTokenRes, error) {
	panic("unimplemented")
}
