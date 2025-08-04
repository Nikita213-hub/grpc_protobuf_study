package redisstore

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/domain/models"
	"github.com/go-redis/redis"
)

type RedisStorage struct {
	client *redis.Client
}

func NewRedisTokenStorage(addr, password string, db int) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &RedisStorage{
		client: client,
	}, nil
}

func (r *RedisStorage) SaveVCode(ctx context.Context, email, code string) error {
	return r.client.Set(email, code, 5*time.Minute).Err()
}

func (r *RedisStorage) GetVCode(ctx context.Context, email string) (string, error) {
	return r.client.Get(email).Result()
}

func (r *RedisStorage) RemoveVCode(ctx context.Context, email string) error {
	return r.client.Del(email).Err()
}

func (r *RedisStorage) AddSession(ctx context.Context, session *models.Session) error {
	if _, err := r.client.Pipelined(func(rdb redis.Pipeliner) error {
		rdb.HSet(session.ID, "user_email", session.Email)
		rdb.HSet(session.ID, "expires_at", session.ExpiresAt)
		rdb.Pipeline().Expire(session.ID, time.Duration(session.ExpiresAt-time.Now().Unix())*time.Second)
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (r *RedisStorage) GetSession(ctx context.Context, sessionID string) (*models.Session, error) {
	result, err := r.client.HGetAll(sessionID).Result()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("key not found")
	}
	return &models.Session{
		ID:        sessionID,
		Email:     result["user_email"],
		ExpiresAt: parseInt(result["expires_at"]),
	}, nil
}

func (r *RedisStorage) RemoveSession(ctx context.Context, sessionId string) error {
	return r.client.Del(sessionId).Err()
}

func parseInt(s string) int64 {
	if s == "" {
		return time.Now().Add(1 * time.Hour).Unix()
	}

	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		if f, ferr := strconv.ParseFloat(s, 64); ferr == nil {
			return int64(f)
		}

		return time.Now().Add(1 * time.Hour).Unix()
	}

	return n
}
