package queue

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/models"
)

func UpdateStatus(id int, newStatus string) error {

	jobIDStr := strconv.Itoa(id)

	hashKey := "job:" + jobIDStr

	jsonData, err := RDB.HGet(Ctx, hashKey, "data").Result()

	if err != nil {
		return fmt.Errorf("Error getting json data\n")
	}

	var job models.Job
	if err := json.Unmarshal([]byte(jsonData), &job); err != nil {
		return fmt.Errorf("Error while unmarshal the json data")
	}

	job.Status = newStatus
	job.UpdatedAt = time.Now()

	updatedData, err := json.Marshal(job)

	if err != nil {
		return fmt.Errorf("Failed to marshal updated data")
	}

	return RDB.HSet(Ctx, hashKey, map[string]interface{}{
		"data":   string(updatedData),
		"status": newStatus,
	}).Err()

}

func HandleFailure(id int) (int, error) {

	hashKey := "job:" + strconv.Itoa(id)

	jsonData, err := RDB.HGet(Ctx, hashKey, "data").Result()

	if err != nil {
		return 0, err
	}
	var job models.Job

	json.Unmarshal([]byte(jsonData), &job)

	job.Retries++
	job.UpdatedAt = time.Now()

	if job.Retries > job.MaxRetries {

		job.Status = models.StateFailed
	} else {
		job.Status = models.StateQueued
	}

	updatedJSON, _ := json.Marshal(job)
	RDB.HSet(Ctx, hashKey, map[string]interface{}{
		"data":   string(updatedJSON),
		"status": job.Status,
	}).Err()

	return job.Retries, nil
}
