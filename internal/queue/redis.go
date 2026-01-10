package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

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

	if job.Status == "" {
		job.Status = models.StateQueued
	}
	if job.MaxRetries == 0 {
		job.MaxRetries = 3
	}
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()

	jsonData, err := json.Marshal(job)

	if err != nil {
		return fmt.Errorf("Failed to marshal job into json\n")
	}

	jobIDStr := strconv.Itoa(job.ID)

	hashKey := "job:" + jobIDStr
	err = RDB.HSet(Ctx, hashKey, map[string]interface{}{
		"data":   jsonData,
		"status": job.Status,
	}).Err()

	if err != nil {
		return fmt.Errorf("Could not save job in Redis Hash!\n")
	}

	return RDB.LPush(Ctx, "job_queue", job.ID).Err()

}
