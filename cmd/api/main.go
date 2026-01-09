package main

import (
	"fmt"

	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/queue"
)

func main() {

	err := queue.Connect()

	if err != nil {
		fmt.Println("Error connecting to redis")
	}
	fmt.Println(" Connected to Redis")
}
