package api

import (
	"encoding/json"
	"net/http"

	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/models"
	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/queue"
)

type CreateJobRequest struct {
	ID int `json:"id"`
}

func HandleCreateJob(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateJobRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	job := models.Job{
		ID: req.ID,
	}
	err := queue.EnqueueJob(&job)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Job enqueued successfully",
		"job_id":  string(rune(req.ID)),
	})

}

func GetMetrics(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	queueDepth, err := queue.RDB.LLen(queue.Ctx, "job_queue").Result()

	if err != nil {
		queueDepth = 0
	}

	var processed, failed, retried, inProgress int64

	iter := queue.RDB.Scan(queue.Ctx, 0, "job:*", 0).Iterator()

	for iter.Next(queue.Ctx) {

		fields, err := queue.RDB.HMGet(queue.Ctx, iter.Val(), "status", "data").Result()

		if err != nil || len(fields) < 2 {
			continue
		}
		status, _ := fields[0].(string)
		var retriesCount int

		if dataJson, ok := fields[1].(string); ok {
			var job models.Job
			if err := json.Unmarshal([]byte(dataJson), &job); err == nil {
				retriesCount = job.Retries
			}
		}

		switch status {
		case models.StateCompleted:
			processed++
		case models.StateFailed:
			failed++
		case models.StateInProgress:
			inProgress++
		case models.StateQueued:

		}

		if retriesCount > 0 {
			retried++
		}

	}

	metrics := map[string]int64{
		"processed":   processed,
		"failed":      failed,
		"in_progress": inProgress,
		"queued":      queueDepth,
		"retried":     retried,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)

}

func GetJobStatus(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {

		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")

	if id == "" {
		http.Error(w, "Missing Job ID", http.StatusBadRequest)
		return
	}

	val, err := queue.RDB.HGet(queue.Ctx, "job:"+id, "data").Result()
	if err != nil {
		http.Error(w, "Job not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(val))
}
