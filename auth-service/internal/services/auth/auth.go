package auth

import (
	"context"
	"time"

	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/domain/models"
)

type Auth struct {
	tokenGetter TokenGetter
	tokenSaver  TokenSaver
}

func NewAuth(ts TokenSaver, tg TokenGetter) *Auth {
	return &Auth{
		tokenGetter: tg,
		tokenSaver:  ts,
	}
}

type TokenSaver interface {
	AddToken(ctx context.Context, token string, userData *models.UserData) error
}

type TokenGetter interface {
	GetToken(ctx context.Context, token string) (*models.UserData, error)
}

// GenToken(ctx context.Context, userEmail, userId string) (string, int64, error)
//
//	VerifyToken(ctx context.Context, token string, exp int64) (bool, error)
//
// TODO: implement
func (a *Auth) GenToken(ctx context.Context, userEmail, userId string) (*models.Token, error) {
	token := userEmail + userId
	exp := time.Now().Add(24 * time.Hour).Unix()

	err := a.tokenSaver.AddToken(ctx, token, &models.UserData{
		UserId:    userId,
		UserEmail: userEmail,
		ExpiresAt: exp,
	})
	if err != nil {
		return nil, err
	}
	return &models.Token{
		Token: token,
		Exp:   exp,
	}, nil
}

func (a *Auth) VerifyToken(ctx context.Context, token string) (*models.UserData, error) {
	userData, err := a.tokenGetter.GetToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return userData, nil
}
