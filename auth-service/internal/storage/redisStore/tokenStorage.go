package redisstore

import (
	"context"
	"strconv"

	"github.com/Nikita213-hub/grpc_protobuf_study/auth-service/internal/domain/models"
	"github.com/go-redis/redis"
)

type RedisTokenStoarage struct {
	client *redis.Client
}

func NewRedisTokenStorage(addr, password string, db int) (*RedisTokenStoarage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}

	return &RedisTokenStoarage{
		client: client,
	}, nil
}

func (rts *RedisTokenStoarage) AddToken(ctx context.Context, token string, userData *models.UserData) error {
	key := token
	if _, err := rts.client.Pipelined(func(rdb redis.Pipeliner) error {
		rdb.HSet(key, "user_id", userData.UserId)
		rdb.HSet(key, "user_email", userData.UserEmail)
		rdb.HSet(key, "expires_at", userData.ExpiresAt)
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (rts *RedisTokenStoarage) GetToken(ctx context.Context, token string) (*models.UserData, error) {
	userData, err := rts.client.HGetAll(token).Result()
	if err != nil {
		return nil, err
	}
	exp, err := strconv.Atoi(userData["expires_at"])
	if err != nil {
		return nil, err
	}
	return &models.UserData{
		UserId:    userData["user_id"],
		UserEmail: userData["user_email"],
		ExpiresAt: int64(exp),
	}, nil
}
