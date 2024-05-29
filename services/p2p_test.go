package services

import (
	"container-manager/types"
	"encoding/json"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestNewP2PService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobQueue := NewMockQueue(ctrl)

	service, err := NewP2PService(jobQueue, 4041)
	require.NoError(t, err)
	require.NotNil(t, service)
	require.NotNil(t, service.ID())
}

func TestP2PServiceStartAndStop(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobQueue1 := NewMockQueue(ctrl)
	jobQueue2 := NewMockQueue(ctrl)

	service1, err := NewP2PService(jobQueue1, 4041)
	require.NoError(t, err)
	require.NotNil(t, service1)

	service2, err := NewP2PService(jobQueue2, 4042)
	require.NoError(t, err)
	require.NotNil(t, service2)

	go service1.Start()
	go service2.Start()

	time.Sleep(2 * time.Second)
	require.Equal(t, 2, service1.host.Peerstore().Peers().Len())
	require.Equal(t, 2, service2.host.Peerstore().Peers().Len())
	service1.Stop()
	service2.Stop()
}

func TestP2PServiceBroadcast(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	jobQueue1 := NewMockQueue(ctrl)
	jobQueue2 := NewMockQueue(ctrl)

	service1, err := NewP2PService(jobQueue1, 4041)
	require.NoError(t, err)
	require.NotNil(t, service1)

	service2, err := NewP2PService(jobQueue2, 4042)
	require.NoError(t, err)
	require.NotNil(t, service2)

	go service1.Start()
	go service2.Start()

	time.Sleep(2 * time.Second)
	require.Equal(t, 2, service1.host.Peerstore().Peers().Len())
	require.Equal(t, 2, service2.host.Peerstore().Peers().Len())

	container := types.Container{
		Image: "alpine",
		Env:   map[string]string{},
	}
	data, err := json.Marshal(container)
	require.NoError(t, err)

	jobQueue2.EXPECT().GetStatus("job-1").Times(1).Return(types.JobStatusPending, false)
	jobQueue2.EXPECT().Enqueue("job-1", container).Times(1)

	msg := Message{
		JobID: "job-1",
		Type:  types.P2PMessageTypeDeployContainer,
		Data:  data,
	}
	err = service1.Broadcast(msg)

	time.Sleep(2 * time.Second)
	service1.Stop()
	service2.Stop()
}

func TestGetHostIP(t *testing.T) {
	ip, err := getHostIP()
	require.NoError(t, err)
	require.NotEmpty(t, ip)
}
