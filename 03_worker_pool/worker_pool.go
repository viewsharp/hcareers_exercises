package main

import (
	"errors"
	"sync"
	"sync/atomic"
)

const MaxGoroutinesCount = 10000

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")
var ErrIncorrectGoroutinesCount = errors.New("incorrect n value")
var ErrGoroutinesCountLimitExceeded = errors.New("incorrect n value")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if n < 1 {
		return ErrIncorrectGoroutinesCount
	}
	if n > MaxGoroutinesCount {
		return ErrGoroutinesCountLimitExceeded
	}
	if m < 0 {
		return ErrErrorsLimitExceeded
	}

	tasksChan := make(chan Task)
	var errorsCount atomic.Int64

	var wg sync.WaitGroup
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			for task := range tasksChan {
				err := task()
				if err != nil {
					errorsCount.Add(1)
				}
			}
		}()
	}

	for _, task := range tasks {
		tasksChan <- task
		if errorsCount.Load() > int64(m) {
			break
		}
	}

	close(tasksChan)
	wg.Wait()

	if errorsCount.Load() > int64(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
