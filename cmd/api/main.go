package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/models"
	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/queue"
	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/worker"
)

func main() {

	err := queue.Connect()

	if err != nil {
		fmt.Println("Error connecting to redis")
	}
	fmt.Println(" Connected to Redis")

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker.StartWorker(i, &wg)
	}

	job := &models.Job{
		ID: 3,
		Status: models.StateQueued,
		Retries: 1,
		MaxRetries: 3,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	queue.EnqueueJob(job)
	fmt.Printf("[Job %d] inserted successfully\n",job.ID)

	wg.Wait()

}
