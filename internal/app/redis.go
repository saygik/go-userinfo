package app

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func (a *App) newRedisConnect(host string, password string, db int) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host,
		Password: password,
		DB:       db,
	})

	ctx := context.Background()
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("redis - unable to connect")
	}
	return redisClient, nil
}
