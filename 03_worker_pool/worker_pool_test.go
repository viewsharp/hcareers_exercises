package main

import (
	"errors"
	"math/rand"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_MLessZero(t *testing.T) {
	err := Run([]Task{}, 1, -1)
	require.ErrorIs(t, err, ErrErrorsLimitExceeded)
}

func Test_NoErrors(t *testing.T) {
	// если задачи работают без ошибок, то выполнятся len(tasks) задач, т.е. все задачи;

	pendingTasks := make(map[int]struct{})
	var pendingTasksMutex sync.Mutex

	var tasks []Task
	for i := 0; i < 1000; i++ {
		pendingTasks[i] = struct{}{}
		tasks = append(tasks, func() error {
			pendingTasksMutex.Lock()
			defer pendingTasksMutex.Unlock()

			delete(pendingTasks, i)
			return nil
		})
	}

	err := Run(tasks, runtime.NumCPU()+1, 0)
	require.NoError(t, err)

	require.Len(t, pendingTasks, 0)
}

func Test_Errors(t *testing.T) {
	// если в первых выполненных M задачах (или вообще всех) происходят ошибки, то всего выполнится не более N+M задач.
	n := 10 + rand.Intn(90)
	m := 10 + rand.Intn(90)

	var completedTasks atomic.Int64

	var tasks []Task
	for i := 0; i < (m+n)*2; i++ {
		tasks = append(tasks, func() error {
			completedTasks.Store(1)
			return errors.New("some error")
		})
	}

	err := Run(tasks, n, m)
	require.ErrorIs(t, err, ErrErrorsLimitExceeded)

	require.LessOrEqual(t, completedTasks.Load(), int64(n+m))
}
