package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

var wg sync.WaitGroup

type Task func() error

func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	jobs := make(chan Task)
	res := make(chan error, n)
	done := make(chan struct{}, 1)
	defer close(jobs)

	for range n {
		wg.Add(1)
		go worker(jobs, res, done)
	}

	var errCounter, errReceived int
	var jobsSent int

	for range min(n, len(tasks)) {
		jobs <- tasks[jobsSent]
		jobsSent++
	}

	for errReceived < len(tasks) {
		if errCounter == m {
			close(done)
			wg.Wait()
			return ErrErrorsLimitExceeded
		}

		err := <-res
		errReceived++

		if err != nil {
			errCounter++
		}

		if jobsSent < len(tasks) {
			jobs <- tasks[jobsSent]
			jobsSent++
		}
	}
	close(done)
	wg.Wait()
	close(res)
	return nil
}

func worker(jobs <-chan Task, res chan<- error, done chan struct{}) {
	defer wg.Done()
	for {
		select {
		case j := <-jobs:
			res <- j()
		case <-done:
			return
		}
	}
}
