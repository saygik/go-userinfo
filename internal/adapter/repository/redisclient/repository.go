package redisclient

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Repository struct {
	cl *redis.Client
}

var ctx = context.Background()

func New(redisClient *redis.Client) *Repository {
	return &Repository{
		cl: redisClient,
	}
}
