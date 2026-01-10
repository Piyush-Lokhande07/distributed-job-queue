package main

import (
	"fmt"
	"sync"

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

	for i:=101;i<=105;i++{
		job := &models.Job{
			ID: i,
		}

		err = queue.EnqueueJob(job)
		if err != nil {
			fmt.Printf("Enqueue Error: %v\n", err)
		} else {
			fmt.Printf("Job %d enqueued successfully!\n",i)
		}
	}

	wg.Wait()

}
