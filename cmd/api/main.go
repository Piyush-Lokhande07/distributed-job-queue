package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/api"
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

	http.HandleFunc("/jobs", api.HandleCreateJob)
	http.HandleFunc("/metrics", api.GetMetrics)
	http.HandleFunc("/status",api.GetJobStatus)

	fmt.Println("Server running on port:[8080]")

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Server failed %v\n", err)
	}

	wg.Wait()

}
