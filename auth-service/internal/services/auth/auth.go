package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/domain/models"
	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/storage"
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
	op := "auth.start_login"

	code := generateVCode()
	fmt.Println("code: ", code) //TODO: remove and implement mail sending
	err := a.vcodeRepo.SaveVCode(ctx, userEmail, code)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *Auth) VerifyCode(ctx context.Context, userEmail, code string) (*models.Session, error) {
	op := "auth.verify_code"
	storedCode, err := a.vcodeRepo.GetVCode(ctx, userEmail)
	if err != nil {
		if errors.Is(err, storage.ErrKeyNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrCodeExpired)
		}
		return nil, fmt.Errorf("%s: code lookup failed: %w", op, err)
	}
	if code != storedCode {
		return nil, fmt.Errorf("%s: %w", op, ErrIncorrectCode)
	}
	session := &models.Session{
		ID:        generateSessionID(),
		Email:     userEmail,
		ExpiresAt: time.Now().Add(2 * time.Hour).Unix(),
	}
	err = a.sessionRepo.AddSession(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, ErrSessionCreateFailed)
	}
	_ = a.vcodeRepo.RemoveVCode(ctx, userEmail)

	return session, nil
}

func (a *Auth) GetSession(ctx context.Context, sessionId string) (*models.Session, error) {
	op := "auth.get_session"
	session, err := a.sessionRepo.GetSession(ctx, sessionId)
	if err != nil {
		if errors.Is(err, storage.ErrKeyNotFound) {
			return nil, fmt.Errorf("%s: %w", op, ErrSessionNotFound)
		}
		return nil, fmt.Errorf("%s: session lookup failed: %w", op, err)
	}

	if time.Now().Unix() > session.ExpiresAt {
		_ = a.sessionRepo.RemoveSession(ctx, sessionId)
		return nil, fmt.Errorf("%s: %w", op, ErrSessionExpired)
	}

	return session, nil
}

func generateSessionID() string {
	return uuid.New().String()
}

func generateVCode() string {
	max := big.NewInt(1000000)
	n, _ := rand.Int(rand.Reader, max)
	return fmt.Sprintf("%06d", n.Int64())
}
