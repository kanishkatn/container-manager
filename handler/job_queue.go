package handler

import (
	"container-manager/types"
	"sync"
)

// job is the object that represents a job to be run.
// ID: The ID of the job
// container: The container to run
type job struct {
	ID        string
	container types.Container
}

// JobQueue is the interface that represents a job queue.
// Enqueue: Enqueues a job to be run
// GetStatus: Gets the status of a job
// Run: Runs the job queue
// Stop: Stops the job queue
type JobQueue interface {
	Enqueue(jobID string, container types.Container) error
	GetStatus(jobID string) (string, bool)
	Run(workerCount int)
	Stop()
}

// JobQueueImpl is the object that represents a job queue.
// jobQueue: The queue of jobs to be run
// jobStatus: The status of each job
// mutex: The mutex to protect the job status
// wg: The wait group to wait for all workers to finish
// quit: The channel to signal workers to quit
type JobQueueImpl struct {
	jobQueue  chan job
	jobStatus map[string]types.JobStatus
	mutex     sync.Mutex
	wg        sync.WaitGroup
	quit      chan bool
}

// NewJobQueue creates a new job queue.
func NewJobQueue(size int) *JobQueueImpl {
	return &JobQueueImpl{
		jobQueue:  make(chan job, size),
		jobStatus: make(map[string]types.JobStatus),
		quit:      make(chan bool),
	}
}

// Enqueue enqueues a job to be run.
func (jm *JobQueueImpl) Enqueue(jobID string, container types.Container) error {
	jm.mutex.Lock()
	defer jm.mutex.Unlock()

	jobToQueue := job{
		ID:        jobID,
		container: container,
	}
	jm.jobQueue <- jobToQueue
	jm.jobStatus[jobToQueue.ID] = types.JobStatusPending
	return nil
}

// GetStatus gets the status of a job.
func (jm *JobQueueImpl) GetStatus(jobID string) (types.JobStatus, bool) {
	jm.mutex.Lock()
	defer jm.mutex.Unlock()

	status, exists := jm.jobStatus[jobID]
	return status, exists
}

// worker runs the jobs in the job queue.
func (jm *JobQueueImpl) worker() {
	defer jm.wg.Done()

	for {
		select {
		case job := <-jm.jobQueue:
			// TODO: Run the job
			jm.mutex.Lock()
			jm.jobStatus[job.ID] = types.JobStatusComplete
			jm.mutex.Unlock()
		case <-jm.quit:
			return
		}
	}
}

// Run runs the job queue.
func (jm *JobQueueImpl) Run(workerCount int) {
	for i := 0; i < workerCount; i++ {
		jm.wg.Add(1)
		go jm.worker()
	}
}

// Stop stops the job queue.
func (jm *JobQueueImpl) Stop() {
	close(jm.quit)
	jm.wg.Wait()
}
