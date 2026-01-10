package queue

import (
	"context"

	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/models"
	"github.com/redis/go-redis/v9"
)

var (
	RDB *redis.Client
	Ctx = context.Background()
)

func Connect() error {

	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 2,
	})

	return RDB.Ping(Ctx).Err()
}

func EnqueueJob(job *models.Job) error {

	return RDB.LPush(Ctx, "job_queue", job.ID).Err()

}
