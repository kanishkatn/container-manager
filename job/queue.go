package job

import (
	"container-manager/types"
	"fmt"
	"sync"
)

var (
	// ErrQueueFull is the error returned when the queue is full
	ErrQueueFull = fmt.Errorf("job queue is full")
)

// job is the object that represents a job to be run.
// ID: The ID of the job
// container: The container to run
type job struct {
	ID        string
	container types.Container
}

// Queue is the interface that represents a job queue.
// Enqueue: Enqueues a job to be run
// GetStatus: Gets the status of a job
// Run: Runs the job queue
// Stop: Stops the job queue
type Queue interface {
	Enqueue(jobID string, container types.Container) error
	GetStatus(jobID string) (types.JobStatus, bool)
	Run(workerCount int)
	Stop()
}

// QueueHandler is the implementation of the job queue interface.
// jobs: The queue of jobs to be run
// jobStatus: The status of each job
// mutex: The mutex to protect the job status
// wg: The wait group to wait for all workers to finish
// quit: The channel to signal workers to quit
type QueueHandler struct {
	jobs          chan job
	jobStatus     map[string]types.JobStatus
	mutex         sync.Mutex
	wg            sync.WaitGroup
	quit          chan bool
	dockerManager DockerManager
}

// NewQueue creates a new job queue.
func NewQueue(size int) *QueueHandler {
	return &QueueHandler{
		jobs:          make(chan job, size),
		jobStatus:     make(map[string]types.JobStatus),
		quit:          make(chan bool),
		dockerManager: newDockerManager(),
	}
}

// Enqueue enqueues a job to be run.
func (q *QueueHandler) Enqueue(jobID string, container types.Container) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	if len(q.jobs) == cap(q.jobs) {
		return ErrQueueFull
	}

	jobToQueue := job{
		ID:        jobID,
		container: container,
	}
	q.jobs <- jobToQueue
	q.jobStatus[jobToQueue.ID] = types.JobStatusPending
	return nil
}

// GetStatus gets the status of a job.
func (q *QueueHandler) GetStatus(jobID string) (types.JobStatus, bool) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	status, exists := q.jobStatus[jobID]
	return status, exists
}

// worker runs the jobs in the job queue.
func (q *QueueHandler) worker() {
	defer q.wg.Done()

	for {
		select {
		case newJob := <-q.jobs:
			q.executeJob(newJob)
		case <-q.quit:
			return
		}
	}
}

// Run runs the job queue.
func (q *QueueHandler) Run(workerCount int) {
	for i := 0; i < workerCount; i++ {
		q.wg.Add(1)
		go q.worker()
	}
}

// Stop stops the job queue.
func (q *QueueHandler) Stop() {
	close(q.quit)
	q.wg.Wait()
}

// executeJob executes a job.
func (q *QueueHandler) executeJob(job job) {
	containerID, err := q.dockerManager.DeployContainer(job.container)
	if err != nil {
		q.updateJobStatus(job.ID, types.JobStatusFailed)
	} else {
		status, err := q.dockerManager.GetContainerStatus(containerID)
		if err != nil {
			q.updateJobStatus(job.ID, types.JobStatusFailed)
		} else if status == "running" {
			q.updateJobStatus(job.ID, types.JobStatusComplete)
		} else {
			q.updateJobStatus(job.ID, types.JobStatusFailed)
		}
	}
}

// updateJobStatus updates the status of a job.
func (q *QueueHandler) updateJobStatus(jobID string, status types.JobStatus) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.jobStatus[jobID] = status
}
