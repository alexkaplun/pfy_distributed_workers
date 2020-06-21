package workers

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexkaplun/pfy_distributed_workers/model"
)

func TestWorkerPool(t *testing.T) {
	// create a test task, that returns error if url.Id%3==0
	var task URLTask = func(url *model.URL) error {
		time.Sleep(time.Second * 1)
		//
		if url.Id%3 == 0 {
			return errors.New(fmt.Sprintf("%v: error", url.Id))
		}
		return nil
	}

	// max tasks in a worker pool
	limit := 5

	// test negative limit for worker pool
	_, err := NewWorkerPool(-5)
	require.Error(t, err)

	// create a worker pool
	wp, err := NewWorkerPool(limit)
	assert.NoError(t, err)

	// failed tasks counter
	failedCnt := 0
	totalTests := 17

	// run tests, 1 through totlTests
	for i := 1; i <= totalTests; i++ {
		errch := wp.Run(context.Background(), task, &model.URL{Id: i})

		// read from errch
		// must be async to ensure we don't block the rest of tasks
		go func() {
			err := <-errch
			if err != nil {
				fmt.Println(err)
				failedCnt++
			}
		}()
	}

	// wait until all workers done end stop the worker pool
	errWait := wp.Wait()
	assert.NoError(t, errWait)

	// assert amount of failed tests
	assert.Equal(t, totalTests/3, failedCnt)
}
