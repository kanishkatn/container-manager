package p2p

import (
	"container-manager/job"
	"container-manager/types"
	"context"
	"encoding/json"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/multiformats/go-multiaddr"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"strings"
)

const (
	// ProtocolID is the protocol ID for the container manager p2p service
	ProtocolID = "/container-manager/1.0.0"
	// ServiceName is the service name for mDNS discovery
	ServiceName = "container-manager"
)

// Message is a P2P message sent between peers
// Type is the message type
// JobID is the identifier of the job
// Data is the message data
type Message struct {
	Type  types.P2PMessageType `json:"type"`
	JobID string               `json:"job_id"`
	Data  json.RawMessage      `json:"data"`
}

// peerNotifee is a notifee for peer discovery
type peerNotifee struct {
	handler *Service
}

// HandlePeerFound is called when a new peer is discovered
func (pn *peerNotifee) HandlePeerFound(pi peer.AddrInfo) {
	logrus.WithField("peer", pi.ID).Info("Peer discovered")
	pn.handler.host.Peerstore().AddAddrs(pi.ID, pi.Addrs, peerstore.PermanentAddrTTL)
}

// newPeerNotifee creates a new peer notifee
func newPeerNotifee(handler *Service) *peerNotifee {
	return &peerNotifee{handler: handler}
}

// P2PService is the interface for a P2P service
// Start starts the service
// Stop stops the service
type P2PService interface {
	ID() string
	Start()
	Broadcast(msg Message) error
	Stop()
}

// Service is a P2P service
// host is the libp2p host
// peers is a map of discovered peers
// peersLock is a mutex for peers
// ctx is the service context
// cancel is the cancel function for the service context
type Service struct {
	host     host.Host
	ctx      context.Context
	cancel   context.CancelFunc
	jobQueue job.Queue
}

// NewP2PService creates a new P2P service
func NewP2PService(jobQueue job.Queue) (*Service, error) {
	ctx, cancel := context.WithCancel(context.Background())

	containerIP, err := getHostIP()
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to get container IP: %w", err)
	}

	var listenAddrs []multiaddr.Multiaddr
	for _, port := range []int{4001, 4002} { // Replace with the actual ports you want to use
		addr, err := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", containerIP, port))
		if err != nil {
			cancel()
			return nil, fmt.Errorf("failed to create multiaddr: %w", err)
		}
		listenAddrs = append(listenAddrs, addr)
	}

	p2pHost, err := libp2p.New(
		libp2p.ListenAddrs(listenAddrs...),
	)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("failed to create libp2p host: %w", err)
	}
	service := &Service{
		host:     p2pHost,
		ctx:      ctx,
		cancel:   cancel,
		jobQueue: jobQueue,
	}

	return service, nil
}

// ID returns the ID of the P2P service
func (s *Service) ID() string {
	return s.host.ID().String()
}

// Start starts the P2P service
func (s *Service) Start() {
	logrus.Trace("Starting P2P Service")
	// Set up mDNS for peer discovery
	notifee := newPeerNotifee(s)
	service := mdns.NewMdnsService(s.host, ServiceName, notifee)
	if err := service.Start(); err != nil {
		logrus.Fatalf("failed to start mDNS service: %v", err)
	}

	s.host.SetStreamHandler(ProtocolID, s.handleStream)
	logrus.Infof("P2P Service started with ID: %s", s.host.ID().String())
}

// Broadcast broadcasts a message to all peers
func (s *Service) Broadcast(msg Message) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	for _, pi := range s.host.Peerstore().Peers() {
		if pi == s.host.ID() {
			continue
		}

		stream, err := s.host.NewStream(s.ctx, pi, ProtocolID)
		if err != nil {
			logrus.Errorf("failed to create stream to peer %s: %v", pi, err)
			continue
		}

		_, err = stream.Write(msgBytes)
		if err != nil {
			logrus.Errorf("failed to write to stream: %v", err)
			continue
		}
	}

	return nil
}

// handleStream handles an incoming stream
func (s *Service) handleStream(stream network.Stream) {
	logrus.Trace("Handling stream")
	defer stream.Close()
	buf := make([]byte, 1024)
	n, err := stream.Read(buf)
	if err != nil {
		logrus.Errorf("failed to read from stream: %v", err)
		return
	}

	var msg Message
	err = json.Unmarshal(buf[:n], &msg)
	if err != nil {
		logrus.Errorf("failed to unmarshal message: %v", err)
		return
	}
	logrus.WithField("message", msg).Trace("Received p2p message")

	if _, ok := s.jobQueue.GetStatus(msg.JobID); ok {
		logrus.WithField("job_id", msg.JobID).Trace("Job has already entered the queue")
		return
	}

	switch msg.Type {
	case types.P2PMessageTypeDeployContainer:
		var container types.Container
		if err := json.Unmarshal(msg.Data, &container); err != nil {
			logrus.Errorf("failed to unmarshal container data: %v", err)
			return
		}

		if err := container.Validate(); err != nil {
			logrus.Errorf("invalid container data: %v", err)
			return
		}

		if err := s.jobQueue.Enqueue(msg.JobID, container); err != nil {
			logrus.Errorf("failed to enqueue job: %v", err)
			return
		}

	default:
		logrus.Warnf("unknown message type: %s", msg.Type)
	}
}

// Stop stops the P2P service
func (s *Service) Stop() {
	logrus.Trace("Stopping P2P Service")
	s.cancel()
	s.host.Close()
}

// getHostIP retrieves the host's IP address
func getHostIP() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("failed to get hostname: %w", err)
	}

	addrs, err := net.LookupHost(hostname)
	if err != nil {
		return "", fmt.Errorf("failed to lookup host: %w", err)
	}

	for _, addr := range addrs {
		if strings.Contains(addr, ".") { // IPv4 address
			return addr, nil
		}
	}

	return "", fmt.Errorf("no IPv4 address found for hostname %s", hostname)
}
