package auth

import (
	"context"
	"errors"
	"net/mail"

	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/domain/models"
	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/validators"
	authV1 "github.com/Nikita213-hub/grpc_protobuf_study/pkg/proto/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Auth interface {
	GenToken(ctx context.Context, userEmail, userId string) (*models.Token, error)
	VerifyToken(ctx context.Context, token string) (*models.UserData, error)
}

type Server struct {
	authV1.UnimplementedAuthServiceServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authV1.RegisterAuthServiceServer(gRPC, &Server{auth: auth})
	reflection.Register(gRPC) // TODO: optional
}

func (s *Server) GenToken(ctx context.Context, req *authV1.GenTokenReq) (*authV1.GenTokenRes, error) {
	userEmail, userId := req.GetUserEmail(), req.GetUserId()
	_, err := mail.ParseAddress(userEmail)
	if err != nil {
		return nil, err
	}
	isValid := validators.IsUserIdValid(userId)
	if !isValid {
		return nil, errors.New("error: invalid user id")
	}
	token, err := s.auth.GenToken(ctx, userEmail, userId)
	if err != nil {
		return nil, err
	}
	return &authV1.GenTokenRes{
		Token: token.Token,
		Exp:   token.Exp,
	}, nil
}

func (s *Server) VerifyToken(ctx context.Context, req *authV1.VerifyTokenReq) (*authV1.VerifyTokenRes, error) {
	toekn := req.GetToken()
	userData, err := s.auth.VerifyToken(ctx, toekn)
	if err != nil {
		return nil, err
	}
	return &authV1.VerifyTokenRes{
		UserId:    userData.UserId,
		UserEmail: userData.UserEmail,
		Exp:       userData.ExpiresAt,
	}, nil
}
