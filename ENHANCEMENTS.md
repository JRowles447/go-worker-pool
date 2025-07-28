## Enhancements List
The following are ideas for enhancements to the project. 

### Results
- Result Queue 

### Metrics
- Metrics tracking
- Metrics: Success, Retry, Failure count reporting
- Metrics: worker utilization, processing time, queued jobs, cancellations, job drops
- Observability: Prometheus metrics scrape
- Improved logging 

### Scaling
- Timeout per task (Context)
- Pluggable retry policies (jitter, exp backoff, max attempts, delays)
- Dynamic Scaling + Throttling (add/ remove workers based on load)

### General 
- Generic interface for worker pool
- Config structs for: worker count, retry count, retry logic
- Public methods for submitting tasks, shutdown workflow, 
- Expose HTTP control plane for: metrics, health, job status, retreiveing stats
- API 
- Priority queues and assignment of tasks to workers

### Task Source
- Kafka as a task source
- Files

### Result Destination
- Kafka as result destination

### Tasks Examples
- Real-world examples: uploading files, send emails, compute hash 

### Testing
- Testing and benchmarks (cancellation, queue full, ordered processing?)
- Fuzz testing
- Benchmarking with different worker counts 
- Chaos testing

### Polish
- CI/CD using github actions 
- Long term running as a service
- Helm chart + Dockerfile for service