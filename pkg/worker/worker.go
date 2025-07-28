package worker

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/jrowles447/go-worker-pool/pkg/task"
)

// worker function
func SpawnWorker(ctx context.Context, id int, tasks <-chan *task.Task, retry chan<- *task.Task, wg *sync.WaitGroup, logger *slog.Logger) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			logger.Info(fmt.Sprintf("program shutting down, Worker%d shutting down.", id))
			return
		case t, ok := <-tasks:
			if !ok {
				logger.Info(fmt.Sprintf("Task queue closed, Worker%d shutting down...", id))
				return
			}
			// do work
			taskError := task.SimulateError(t)
			if se, ok := taskError.(*task.SystemError); ok {
				// check retriable
				if se.Retriable {
					// queue task in retry queue
					logger.Info(fmt.Sprintf("Encountered retriable error, placing Task:%s in retry queue", se.Task.Name))
					retry <- se.Task
				} else {
					// drop task
					logger.Info(fmt.Sprintf("%s encountered non-retriable error", se.Task.Name))
					continue
				}
			} else { // simulate work
				time.Sleep(1 * time.Second)
				logger.Info(fmt.Sprintf("Worker%d running %s.", id, t.Name))
			}
		}
	}
}

func SpawnRetryWorker(ctx context.Context, id int, retry <-chan *task.Task, wg *sync.WaitGroup, logger *slog.Logger) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			logger.Info(fmt.Sprintf("program shutting down, RetryWorker%d shutting down.", id))
			return
		case t, ok := <-retry:
			if !ok {
				logger.Info(fmt.Sprintf("Retry queue closed, RetryWorker%d shutting down...", id))
				return
			}
			// do work
			taskError := task.SimulateError(t)
			if se, ok := taskError.(*task.SystemError); ok {
				// check retriable
				if se.Retriable {
					if se.Task.RetryCount < se.Task.MaxRetries {
						logger.Info(fmt.Sprintf("Encountered retriable error, placing Task:%s in retry queue with retry count:%d", se.Task.Name, se.Task.RetryCount))

						// queue task in retry submission queue to be handled by dispatcher
						select {
						case <-ctx.Done():
							logger.Info(fmt.Sprintf("RetryWorker%d shutting down before re-submitting task", id))
							return
						case <-time.After(1 * time.Second):
							logger.Info(fmt.Sprintf("RetryWorker timeout on re-submitting %s, dropping.", se.Task.Name))
						}
					} else {
						logger.Info(fmt.Sprintf("Max retry attempts reached, dropping Task:%s", se.Task.Name))
						continue
					}
				} else {
					// drop task
					logger.Info(fmt.Sprintf("Encountered non-retriable error, cancelling Task:%s", se.Task.Name))
					continue
				}
			} else { // simulate work
				time.Sleep(1 * time.Second)
				logger.Info(fmt.Sprintf("RetryWorker%d running %s.", id, t.Name))
			}
		}
	}
}
