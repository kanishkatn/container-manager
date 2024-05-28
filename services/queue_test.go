package services

import (
	"container-manager/types"
	"fmt"
	"go.uber.org/mock/gomock"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJobQueueImpl(t *testing.T) {
	queueSize := 10
	workerCount := 3
	jobCount := 5

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock Docker service
	mockDockerService := NewMockDockerService(ctrl)
	mockDockerService.EXPECT().DeployContainer(types.Container{}).Times(jobCount).Return("container-id", nil)
	mockDockerService.EXPECT().GetContainerStatus("container-id").Times(jobCount).Return("running", nil)

	// Create a new job queue
	jobQueue := NewQueue(queueSize, mockDockerService)

	// Enqueue some jobs
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
	jobCount := 50

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create a mock Docker service
	mockDockerService := NewMockDockerService(ctrl)
	mockDockerService.EXPECT().DeployContainer(types.Container{}).Times(jobCount).Return("container-id", nil)
	mockDockerService.EXPECT().GetContainerStatus("container-id").Times(jobCount).Return("running", nil)

	// Create a new job queue
	jobQueue := NewQueue(queueSize, mockDockerService)

	// Enqueue jobs concurrently
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
