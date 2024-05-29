package services

import (
	"container-manager/types"
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

var (
	// ErrQueueFull is the error returned when the queue is full
	ErrQueueFull = fmt.Errorf("job queue is full")
)

// job is the object that represents a job to be run.
// id: The ID of the job
// container: The container to run
type job struct {
	id        string
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
	dockerService DockerService
}

// NewQueue creates a new job queue.
func NewQueue(size int, ds DockerService) *QueueHandler {
	return &QueueHandler{
		jobs:          make(chan job, size),
		jobStatus:     make(map[string]types.JobStatus),
		quit:          make(chan bool),
		dockerService: ds,
	}
}

// Enqueue enqueues a job to be run.
func (q *QueueHandler) Enqueue(jobID string, container types.Container) error {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	logrus.WithField("job_id", jobID).Debug("enqueuing job")

	if len(q.jobs) == cap(q.jobs) {
		return ErrQueueFull
	}

	jobToQueue := job{
		id:        jobID,
		container: container,
	}
	q.jobs <- jobToQueue
	q.jobStatus[jobToQueue.id] = types.JobStatusPending
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
			logrus.WithField("job_id", newJob.id).Info("running job")
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
	// TODO: retry mechanism
	containerID, err := q.dockerService.DeployContainer(job.container)
	if err != nil {
		logrus.WithField("job_id", job.id).Errorf("failed to deploy container: %v", err)
		q.updateJobStatus(job.id, types.JobStatusFailed)
	} else {
		status, err := q.dockerService.GetContainerStatus(containerID)
		switch {
		case err != nil:
			q.updateJobStatus(job.id, types.JobStatusFailed)
		case status == "running":
			q.updateJobStatus(job.id, types.JobStatusComplete)
		default:
			q.updateJobStatus(job.id, types.JobStatusFailed)
		}
		logrus.WithField("job_id", job.id).Infof("container deployed successfully")
	}
}

// updateJobStatus updates the status of a job.
func (q *QueueHandler) updateJobStatus(jobID string, status types.JobStatus) {
	q.mutex.Lock()
	defer q.mutex.Unlock()

	q.jobStatus[jobID] = status
}
