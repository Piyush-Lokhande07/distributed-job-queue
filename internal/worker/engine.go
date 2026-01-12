package worker

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/models"
	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/queue"
	"github.com/redis/go-redis/v9"
)

func StartWorker(ctx context.Context, ID int, wg *sync.WaitGroup) {
	defer wg.Done()

	workerLog := slog.With("worker_id", ID)
	workerLog.Info("Worker Started and waiting for jobs")

	for {

		select {
		case <-ctx.Done():
			workerLog.Info("Worker shutting down")
			return
		default:
			result, err := queue.RDB.BRPop(queue.Ctx, time.Second, "job_queue").Result()
			if err != nil {
				if err == redis.Nil {
					continue
				}
				workerLog.Error("Error while popping job", "error", err)
				continue
			}
			jobID := result[1]

			idInt, err := strconv.Atoi(jobID)
			if err != nil {
				workerLog.Error("Error converting ID", "error", err)
				continue
			}

			err = PerformWork(ID, idInt, workerLog)

			if err != nil {

				currentRetries, _ := queue.HandleFailure(idInt)
				jobLog := workerLog.With("job_log", idInt)

				if currentRetries <= 3 {
					delay := time.Duration(1<<uint(currentRetries)) * time.Second
					jobLog.Warn("Job failed, scheduling retry", "delay", delay.String())

					go func(id string, d time.Duration) {
						time.Sleep(d)
						err := queue.RDB.LPush(queue.Ctx, "job_queue", id).Err()

						if err != nil {
							jobLog.Error("CRITICAL: Failed to re-queue Job", "error", err)
						}

					}(jobID, delay)

				} else {
					jobLog.Info("Permanent job failure")
				}

				continue

			}

		}
	}

}

func PerformWork(wId int, jId int, baseLog *slog.Logger) error {

	jobLog := baseLog.With("job_id", jId)

	if time.Now().UnixNano()%2 == 0 {
		jobLog.Debug("Simulated error triggered")
		return fmt.Errorf("Simulated error")
	}

	err := queue.UpdateStatus(jId, models.StateInProgress)

	if err != nil {
		jobLog.Error("Error updating status")
	}

	jobLog.Info("Processing the Job")
	time.Sleep(2 * time.Second)

	err = queue.UpdateStatus(jId, models.StateCompleted)

	if err != nil {
		jobLog.Error("Error updating status for Job")
	}

	jobLog.Info("Finished Job!")
	time.Sleep(2 * time.Second)
	return nil
}
