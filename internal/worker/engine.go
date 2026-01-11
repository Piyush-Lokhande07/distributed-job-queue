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

		err = PerformWork(ID,idInt)

		if err != nil {

			currentRetries, _ := queue.HandleFailure(idInt)
			fmt.Printf("[Worker %d] Job %d failed. Attempting retry...[%d]\n", ID, idInt,currentRetries)

			

			if currentRetries <= 3 {
				delay := time.Duration(1<<uint(currentRetries)) * time.Second

				fmt.Printf("[Worker %d] Retrying Job: %d in %v..\n",ID,idInt,delay)

				go func(id string,d time.Duration){
					time.Sleep(d)
					err:=queue.RDB.LPush(queue.Ctx, "job_queue",id).Err()

					if err!=nil{
						fmt.Printf("CRITICAL: Failed to re-queue Job:%s: %v\n", jobID,err)
					}

				}(jobID, delay)

			}else{
				fmt.Printf("[Worker %d] Permanent Job failed for Job:%d\n",ID,idInt)
			}

			continue

			
		}

		

		
	}

}

func PerformWork(wId int, jId int) error {

	if time.Now().UnixNano()%2 == 0 {
		return fmt.Errorf("Simulated error")
	}

	fmt.Printf("[Worker %d] Grabbed the Job:%d\n", wId, jId)
	err := queue.UpdateStatus(jId, models.StateInProgress)

	if err != nil {
		fmt.Printf("[Worker %d] Error updating status for Job:%d %v\n",wId,jId, err)
	}

	fmt.Printf("[Worker %d] Processing the Job:%d\n",wId, jId)
	time.Sleep(2 * time.Second)

	err = queue.UpdateStatus(jId, models.StateCompleted)

	if err != nil {
		fmt.Printf("[Worker %d] Error updating status for Job:%d %v\n",wId,jId, err)
	}

	fmt.Printf("[Worker %d] Finished Job:%d \n",wId,jId)
	time.Sleep(2 * time.Second)
	return nil
}
