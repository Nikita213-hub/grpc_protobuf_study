package auth

import (
	"context"
	"errors"
	"time"

	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/domain/models"
	"github.com/google/uuid"
)

var (
	ErrIncorrectPassword = errors.New("incorrect password")
)

type Auth struct {
	sessionRepo SessionsRepository
	vcodeRepo   VCodeRepository
}

func NewAuth(sr SessionsRepository, vcr VCodeRepository) *Auth {
	return &Auth{
		sessionRepo: sr,
		vcodeRepo:   vcr,
	}
}

type SessionsRepository interface {
	AddSession(ctx context.Context, session *models.Session) error
	GetSession(ctx context.Context, sessionId string) (*models.Session, error)
	RemoveSession(ctx context.Context, sessionId string) error
}

type VCodeRepository interface {
	SaveVCode(ctx context.Context, email string, code string) error
	GetVCode(ctx context.Context, email string) (string, error)
	RemoveVCode(ctx context.Context, email string) error
}

// GenToken(ctx context.Context, userEmail, userId string) (string, int64, error)
//
//	VerifyToken(ctx context.Context, token string, exp int64) (bool, error)
//
// TODO: implement

func (a *Auth) StartLogin(ctx context.Context, userEmail string) error {
	code := "1111" //TODO: add random code impl
	err := a.vcodeRepo.SaveVCode(ctx, userEmail, code)
	if err != nil {
		return err
	}
	return nil
}

func (a *Auth) VerifyCode(ctx context.Context, userEmail, code string) (*models.Session, error) {
	storedCode, err := a.vcodeRepo.GetVCode(ctx, userEmail)
	if err != nil {
		return nil, err
	}
	if code != storedCode {
		return nil, errors.New("incorrect code")
	}
	session := &models.Session{
		ID:        generateSessionID(),
		Email:     userEmail,
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour).Unix(),
	}
	err = a.sessionRepo.AddSession(ctx, session)
	if err != nil {
		return nil, err
	}
	err = a.vcodeRepo.RemoveVCode(ctx, userEmail)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (a *Auth) GetSession(ctx context.Context, sessionId string) (*models.Session, error) {
	session, err := a.sessionRepo.GetSession(ctx, sessionId)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func generateSessionID() string {
	return uuid.New().String()
}
