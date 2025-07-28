# go-worker-pool

A simple Go worker pool example implementing concurrency patterns with additional features to increase robustness. 

## Features 
Includes the following featurese to improve the robustness of the task processing pipeline: 
- Separate retry workers 
- Retry backoff
- Context based shutdown
- Bounded queues for backpressure (drops tasks to be queued if the retry queue is full)


## Running 
To run the project, execute the following from the project root directory: 

```bash
go run cmd/main.go
```
