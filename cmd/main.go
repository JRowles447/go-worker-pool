package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jrowles447/go-worker-pool/internal/logger"
	"github.com/jrowles447/go-worker-pool/pkg/task"
	"github.com/jrowles447/go-worker-pool/pkg/worker"
)

func main() {
	logger := logger.NewLogger()

	logger.Info("Logger set up.")

	// define number of workers
	// define task counts
	numWorkers := 3
	numTasks := 10
	numRetryWorkers := 1

	tasks := make(chan *task.Task, 20)
	retryQueue := make(chan *task.Task, 20)

	workerPool := sync.WaitGroup{}
	retryPool := sync.WaitGroup{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(15*time.Second))
	defer cancel()

	// gracefully handle shutdowns (Ctrl+C, SIGTERM)
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		sig := <-sigCh
		logger.Info("Shutdown signal received", "signal", sig)
		cancel()
	}()

	// spawn workers
	for i := range numWorkers {
		workerPool.Add(1)
		go worker.SpawnWorker(ctx, i, tasks, retryQueue, &workerPool, logger)
	}
	logger.Info("Workers started")

	// spawn retry workers
	for i := range numRetryWorkers {
		retryPool.Add(1)
		go worker.SpawnRetryWorker(ctx, i, retryQueue, &retryPool, logger)
	}
	logger.Info("RetryWorkers started")

	// queue tasks
	for i := range numTasks {
		tasks <- &task.Task{Name: fmt.Sprintf("Task%d", i), MaxRetries: 3}
	}
	close(tasks)

	logger.Info("Tasks queued")
	workerPool.Wait()
	logger.Info("Primary workers completed")

	retryPool.Wait()
	logger.Info("RetryWorkers completed")

	logger.Info("All tasks completed, shutting down!")
}
