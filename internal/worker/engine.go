package worker

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/models"
	"github.com/Piyush-Lokhande07/distributed-job-queue/internal/queue"
)

func StartWorker(ID int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("[Worker %d] Started and waiting for jobs...\n", ID)

	for {

		result, err := queue.RDB.BRPop(queue.Ctx, 0, "job_queue").Result()
		if err != nil {
			fmt.Printf("[Worker %d] Error while popping job %v\n", ID, err)
		}
		jobID := result[1]

		idInt, err := strconv.Atoi(jobID)
		if err != nil {
			fmt.Printf("Error converting ID: %v\n", err)
			continue
		}

		fmt.Printf("[Worker %d] grabbed the job:%s\n", ID, jobID)

		err = queue.UpdateStatus(idInt, models.StateInProgress)
		if err != nil {

			fmt.Printf("[Worker %d] Error updating status %v\n", ID, err)
		}

		fmt.Printf("[Worker %d] Processing the Job %s\n", ID, jobID)
		time.Sleep(2 * time.Second)

		err = queue.UpdateStatus(idInt, models.StateCompleted)

		if err != nil {

			fmt.Printf("[Worker %d] Error updating status %v\n", ID, err)
		}

		fmt.Printf("[Worker %d] finished Job \n", ID)
	}

}
