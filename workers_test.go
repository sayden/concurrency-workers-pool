package main

import (
	"testing"
	"time"
)

//worker receives a channel where the jobs will come from. He has to return the
//channel to available workers after use
func TestWorker(t *testing.T) {
	idleWorkers := make(chan *chan int, 1)
	jobCh := make(chan int)

	go worker(&jobCh, &idleWorkers)

	jobCh <- 33

	time.Sleep((TIMEOUT_TIME_MILLISECONDS + 200) * time.Millisecond)

	if len(idleWorkers) != 1 {
		t.Error("Worker didn't add itself to queue after job")
	}

}

//createWorkers will create N goroutines to process jobs concurrently
func TestCreateWorkers(t *testing.T) {
	idleWorkers := make(chan *chan int, WORKERS_SIZE)

	createWorkers(&idleWorkers)

	if len(idleWorkers) != WORKERS_SIZE {
		t.Errorf("Size of idle workers after a creating the workers must be " +
		"%d but was %d\n", WORKERS_SIZE, len(idleWorkers))
	}

	worker := <-idleWorkers

	*worker <- 33

	//Wait half of the timeout time and check that our idle number is WORKERS_SIZE - 10
	time.Sleep(time.Millisecond * (WORKERS_SIZE + TIMEOUT_TIME_MILLISECONDS -
	(TIMEOUT_TIME_MILLISECONDS / 2)))

	if len(idleWorkers) != 9 {
		t.Errorf("The size of the pool just after sending a job must shrink by " +
		"one so we must have %d but we had %d", WORKERS_SIZE-1, len(idleWorkers))
	}

	//Wait the other half now so we give time to finish the job
	time.Sleep(time.Millisecond * (WORKERS_SIZE + TIMEOUT_TIME_MILLISECONDS -
	(TIMEOUT_TIME_MILLISECONDS / 2)))

	if len(idleWorkers) != 10 {
		t.Errorf("After sending a job to the worker and waiting a bit more than " +
		"%d the number of idle workers must be again %d\n", TIMEOUT_TIME_MILLISECONDS, WORKERS_SIZE)
	}
}

//dispatcher will receive messages and use the availableWorkers to get idle workers
func TestDispatcher(t *testing.T){
	idleWorkers := make(chan *chan int, WORKERS_SIZE)
	jobCh := make(chan int)

	go dispatcher(&jobCh, &idleWorkers)
	createWorkers(&idleWorkers)

	if len(idleWorkers) != WORKERS_SIZE {
		t.Errorf("Size of idle workers after a creating the workers must be " +
		"%d but was %d\n", WORKERS_SIZE, len(idleWorkers))
	}

	//Quickly send WORKERS_SIZE + 20 jobs so that queue size must be zero
	for i:=0; i<WORKERS_SIZE + 20; i++ {
		jobCh <- i
	}

	if len(idleWorkers) != 0 {
		t.Errorf("Size of idle workers must be 0 after sending 12 fast requests" +
		"but was %d\n", len(idleWorkers))
	}

	//If I wait at least a second, size of idle workers must be 10 again
	time.Sleep(time.Millisecond * (TIMEOUT_TIME_MILLISECONDS + 200))
	if len(idleWorkers) != 10 {
		t.Errorf("Size of idle workers must be 10 after waiting %d ms but " +
		"was %d\n", TIMEOUT_TIME_MILLISECONDS, len(idleWorkers))
	}
}