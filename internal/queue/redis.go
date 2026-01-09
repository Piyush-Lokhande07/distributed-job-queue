package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	RDB * redis.Client
	Ctx = context.Background()
)

func Connect() error {

		RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password
		DB:       0,  // use default DB
		Protocol: 2,
	})

	return RDB.Ping(Ctx).Err()
}
