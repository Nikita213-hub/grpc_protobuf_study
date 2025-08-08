package auth

import (
	"context"
	"net/mail"

	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/domain/models"
	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/logger"
	authV2 "github.com/Nikita213-hub/travel_proto_contracts/pkg/proto/auth/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Auth interface {
	StartLogin(ctx context.Context, userEmail string) error
	VerifyCode(ctx context.Context, userEmail, code string) (*models.Session, error)
	GetSession(ctx context.Context, sessionId string) (*models.Session, error)
}

type Server struct {
	authV2.UnimplementedAuthServiceServer
	auth   Auth
	logger *logger.Logger
}

func Register(gRPC *grpc.Server, auth Auth, log *logger.Logger) {
	authV2.RegisterAuthServiceServer(gRPC, &Server{auth: auth, logger: log})
	reflection.Register(gRPC) // TODO: optional for production
}

func (s *Server) StartLogin(ctx context.Context, req *authV2.LoginRequest) (*authV2.LoginResponse, error) {
	userEmail := req.GetEmail()
	_, err := mail.ParseAddress(userEmail)
	if err != nil {
		s.logger.Warn("start_login_invalid_input")
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}
	err = s.auth.StartLogin(ctx, userEmail)
	if err != nil {
		s.logger.Error("start_login_failed", "error", err.Error())
		return nil, convertServiceError(err)
	}
	s.logger.Info("start_login_success", "user_email", userEmail)
	return &authV2.LoginResponse{
		Message: "Code was sent to ur email",
	}, nil
}

func (s *Server) VerifyCode(ctx context.Context, req *authV2.VerifyCodeRequest) (*authV2.VerifyCodeResponse, error) {
	email := req.GetEmail()
	code := req.GetCode()

	if err := validateEmail(email); err != nil {
		s.logger.Warn("verify_code_invalid_input")
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}

	if code == "" || len(code) != 6 {
		s.logger.Warn("verify_code_invalid_input")
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}

	session, err := s.auth.VerifyCode(ctx, email, code)
	if err != nil {
		s.logger.Warn("verify_code_failed", "error", err.Error())
		return nil, convertServiceError(err)
	}
	s.logger.Info("verify_code_success", "user_email", email, "session_id", session.ID)
	return &authV2.VerifyCodeResponse{SessionId: session.ID}, nil
}

func (s *Server) GetSession(ctx context.Context, req *authV2.SessionRequest) (*authV2.Session, error) {
	sessionId := req.GetSessionId()
	session, err := s.auth.GetSession(ctx, sessionId)
	if err != nil {
		s.logger.Warn("get_session_failed", "error", err.Error())
		return nil, convertServiceError(err)
	}
	s.logger.Info("get_session_success", "user_email", session.Email)
	return &authV2.Session{
		Id:        session.ID,
		Email:     session.Email,
		ExpiresAt: session.ExpiresAt,
	}, nil
}

func validateEmail(email string) error {
	if email == "" {
		return ErrEmptyEmail
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return ErrInvalidEmail
	}
	return nil
}
