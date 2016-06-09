package main

import (
	"time"
	"fmt"
	"math/rand"
)

const (
	WORKERS_SIZE = 10
	TIMEOUT_TIME_MILLISECONDS = 1000
)

//createWorkers will create N goroutines to process jobs concurrently
func createWorkers(idleWorkers *chan *chan int) {
	for i := 0; i < WORKERS_SIZE; i++ {
		ch := make(chan int)
		go worker(&ch, idleWorkers)
		*idleWorkers <- &ch
	}
}

//Lets create an infinite loop of requests every 20 milliseconds so we force the
//discard of rcal
func main() {
	//Concurrency party starts here with initialization. We could use the "init"
	//function but we won't be able to unit test it correctly

	//idleWorkers is an queue of goroutines workers
	idleWorkers := make(chan *chan int, WORKERS_SIZE)

	//A channel to communicate to the job dispatcher
	jobsCh := make(chan int)

	//Create the SPF dispatcher in a different process
	go dispatcher(&jobsCh, &idleWorkers)

	createWorkers(&idleWorkers)


	iter:=0
	for {
		rand.Seed(int64(iter))
		wait := rand.Int31n(200)
		jobsCh <- iter
		iter++
		time.Sleep(time.Duration(wait) * time.Millisecond)
	}
}

//dispatcher will receive messages and use the availableWorkers to get idle workers
func dispatcher(c *chan int, idleWorkers *chan *chan int) {
	for {
		value := <-*c

		//Check if we don't have idle works so we discard the message
		if(len(*idleWorkers) == 0) {
			dropCall(value)
			continue
		}

		worker := <-*idleWorkers
		*worker <- value
	}
}

//worker receives a channel where the jobs will come from. He has to return the
//channel to available workers after use
func worker(c *chan int, idleWorkers *chan *chan int) {
	for {
		value := <-*c
		doActuallyCall(value)
		time.Sleep(time.Millisecond * TIMEOUT_TIME_MILLISECONDS)
		*idleWorkers <- c
	}
}

func doActuallyCall(i int) {
	fmt.Printf("%04d: Called\n", i)
}

func dropCall(i int) {
	fmt.Printf("%04d: Dropped\n", i)
}
