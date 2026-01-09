package worker

import (
	"fmt"
	"sync"
	"time"
)

func StartWorker(ID int, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("[Worker %d] Started and waiting for jobs...\n", ID)

	fmt.Printf("[Worker %d] Processing Job \n", ID)

	time.Sleep(2 * time.Second)

	fmt.Printf("[Worker %d] finished Job \n", ID)

}
