package job_test

import (
	"container-manager/job"
	"container-manager/types"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobQueueImpl(t *testing.T) {
	queueSize := 10
	workerCount := 3

	// Create a new job queue
	jobQueue, err := job.NewQueue(queueSize)
	require.NoError(t, err)

	// Enqueue some jobs
	jobCount := 5
	for i := 0; i < jobCount; i++ {
		jobID := fmt.Sprintf("job-%d", i)
		err := jobQueue.Enqueue(jobID, types.Container{})
		require.NoError(t, err)
	}

	// Run the job queue with workerCount workers
	go jobQueue.Run(workerCount)

	// Give some time for jobs to be processed
	time.Sleep(2 * time.Second)

	// Stop the job queue
	jobQueue.Stop()

	// Check the job statuses
	for i := 0; i < jobCount; i++ {
		jobID := fmt.Sprintf("job-%d", i)
		status, exists := jobQueue.GetStatus(jobID)
		assert.True(t, exists)
		assert.Equal(t, types.JobStatusComplete, status)
	}
}

func TestJobQueueImplConcurrentEnqueue(t *testing.T) {
	queueSize := 100
	workerCount := 5

	// Create a new job queue
	jobQueue, err := job.NewQueue(queueSize)
	require.NoError(t, err)

	// Enqueue jobs concurrently
	jobCount := 50
	var wg sync.WaitGroup
	for i := 0; i < jobCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			jobID := fmt.Sprintf("job-%d", i)
			err := jobQueue.Enqueue(jobID, types.Container{})
			require.NoError(t, err)
		}(i)
	}

	// Wait for all jobs to be enqueued
	wg.Wait()

	// Run the job queue with workerCount workers
	go jobQueue.Run(workerCount)

	// Give some time for jobs to be processed
	time.Sleep(2 * time.Second)

	// Stop the job queue
	jobQueue.Stop()

	// Check the job statuses
	for i := 0; i < jobCount; i++ {
		jobID := fmt.Sprintf("job-%d", i)
		status, exists := jobQueue.GetStatus(jobID)
		assert.True(t, exists)
		assert.Equal(t, types.JobStatusComplete, status)
	}
}
