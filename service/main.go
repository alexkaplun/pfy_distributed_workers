package service

import (
	"context"
	"log"
	"os"

	"github.com/alexkaplun/pfy_distributed_workers/storage"
	workers "github.com/alexkaplun/pfy_distributed_workers/workerpool"
)

var logger = log.New(os.Stdout, "service: ", log.Ldate|log.Ltime|log.Lshortfile)

type Service struct {
	storage    *storage.Storage
	maxWorkers int
}

func New(storage *storage.Storage, maxWorkers int) *Service {
	return &Service{
		storage:    storage,
		maxWorkers: maxWorkers,
	}
}

func (s *Service) Run() {
	// create the worker pool
	wp, err := workers.NewWorkerPool(s.maxWorkers)
	if err != nil {
		panic(err)
	}

	// get the urls from the database
	// TODO: should we re-run the PROCESSING?
	urls, err := s.storage.GetNewURLList()
	if err != nil {
		logger.Println(err)
		return
	}

	logger.Printf("total %v records to process", len(urls))

	// run the tasks
	for _, url := range urls {
		errch := wp.Run(context.Background(), s.processUrl, url)

		go func() {
			err := <-errch
			if err != nil {
				logger.Println(err)
			}
		}()
	}

	// wait until all the tasks done
	err = wp.Wait()
	if err != nil {
		logger.Println("Failed to complete the tasks: ", err)
		return
	}
}
