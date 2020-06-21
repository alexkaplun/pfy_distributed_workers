package workers

import (
	"context"
	"errors"

	"github.com/alexkaplun/pfy_distributed_workers/model"
	"golang.org/x/sync/semaphore"
)

type URLTask func(url *model.URL) error

type WorkerPool struct {
	limit int
	sem   *semaphore.Weighted
}

func NewWorkerPool(limit int) (*WorkerPool, error) {
	if limit <= 0 {
		return nil, errors.New("limit of worker pool must be >0")
	}

	return &WorkerPool{
		limit: limit,
		sem:   semaphore.NewWeighted(int64(limit)),
	}, nil
}

func (wp *WorkerPool) Run(ctx context.Context, task URLTask, data *model.URL) <-chan error {
	errchan := make(chan error, 1)

	err := wp.sem.Acquire(ctx, 1)
	if err != nil {
		errchan <- err
		close(errchan)
		return errchan
	}

	go func() {
		defer wp.sem.Release(1)
		defer close(errchan)

		err = task(data)
		if err != nil {
			errchan <- err
		}
	}()

	return errchan
}

func (wp *WorkerPool) Wait() error {
	// acquire all available slots in semaphore
	for i := 0; i < wp.limit; i++ {
		err := wp.sem.Acquire(context.Background(), 1)
		if err != nil {
			return err
		}
	}

	// all tasks have completed; release the semaphore
	wp.sem.Release(int64(wp.limit))

	return nil
}
