package main

import (
	"fmt"
	"sync"
	"time"
)

type Request struct {
	Id int
}

var (
	totalProcessed int
	mutex_id       sync.Mutex
)

func Worker(req_ID chan Request, workerId int) {

	for req := range req_ID {

		fmt.Printf("Worker %d received Request %d\n", workerId, req.Id)

		done := make(chan bool)

		// Simulating backend processing
		go func(r Request) {

			// Simulate slow request
			if r.Id%4 == 0 {
				time.Sleep(3 * time.Second)
			} else {
				time.Sleep(1 * time.Second)
			}

			done <- true

		}(req)

		select {

		case <-done:

			mutex_id.Lock()
			totalProcessed++
			mutex_id.Unlock()

			fmt.Printf("Worker %d completed Request %d\n", workerId, req.Id)

		case <-time.After(2 * time.Second):

			fmt.Printf("Worker %d timeout for Request %d\n", workerId, req.Id)
		}
	}
}

func implement() {

	request := make(chan Request)

	// Starting worker pool
	for i := 1; i <= 3; i++ {
		go Worker(request, i)
	}

	// Sending requests
	for i := 1; i <= 10; i++ {
		request <- Request{Id: i}
	}

	close(request)

	// Temporary wait
	time.Sleep(10 * time.Second)

	mutex_id.Lock()
	fmt.Println("Total Processed:", totalProcessed)
	mutex_id.Unlock()
}
