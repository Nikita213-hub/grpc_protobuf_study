package auth

import (
	"context"
	"net/mail"

	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/domain/models"
	authV2 "github.com/Nikita213-hub/travel_proto_contracts/pkg/proto/auth/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Auth interface {
	StartLogin(ctx context.Context, userEmail string) error
	VerifyCode(ctx context.Context, userEmail, code string) (*models.Session, error)
	GetSession(ctx context.Context, sessionId string) (*models.Session, error)
}

type Server struct {
	authV2.UnimplementedAuthServiceServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authV2.RegisterAuthServiceServer(gRPC, &Server{auth: auth})
	reflection.Register(gRPC) // TODO: optional
}

func (s *Server) StartLogin(ctx context.Context, req *authV2.LoginRequest) (*authV2.LoginResponse, error) {
	userEmail := req.GetEmail()
	_, err := mail.ParseAddress(userEmail)
	if err != nil {
		return nil, err
	}
	err = s.auth.StartLogin(ctx, userEmail)
	if err != nil {
		return nil, err
	}
	return &authV2.LoginResponse{
		Message: "Code was sent to ur email",
	}, nil
}

func (s *Server) VerifyCode(ctx context.Context, req *authV2.VerifyCodeRequest) (*authV2.VerifyCodeResponse, error) {
	session, err := s.auth.VerifyCode(ctx, req.Email, req.Code)
	if err != nil {
		return nil, err
	}
	return &authV2.VerifyCodeResponse{SessionId: session.ID}, nil
}

func (s *Server) GetSession(ctx context.Context, req *authV2.SessionRequest) (*authV2.Session, error) {
	sessionId := req.GetSessionId()
	session, err := s.auth.GetSession(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	return &authV2.Session{
		Id:        session.ID,
		Email:     session.Email,
		ExpiresAt: session.ExpiresAt,
	}, nil
}
