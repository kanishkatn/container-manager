package p2p

import (
	"container-manager/types"
	"context"
	"encoding/json"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/sirupsen/logrus"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

const (
	// ProtocolID is the protocol ID for the container manager p2p service
	ProtocolID = "/container-manager/1.0.0"
	// ServiceName is the service name for mDNS discovery
	ServiceName = "container-manager"
	// CleanupInterval is the interval at which the message hash map is cleaned up
	CleanupInterval = 10 * time.Minute
	// MessageTTL is the time-to-live for each message hash
	MessageTTL = 30 * time.Minute
)

// Message is a P2P message sent between peers
// Type is the message type
// Hash is the message hash
// Data is the message data
type Message struct {
	Type string          `json:"type"`
	Hash string          `json:"hash"`
	Data json.RawMessage `json:"data"`
}

// peerNotifee is a notifee for peer discovery
type peerNotifee struct {
	handler *Service
}

// HandlePeerFound is called when a new peer is discovered
func (pn *peerNotifee) HandlePeerFound(pi peer.AddrInfo) {
	logrus.WithField("peer", pi.ID).Trace("Peer discovered")
	pn.handler.peersLock.Lock()
	pn.handler.peers[pi.ID] = pi
	pn.handler.peersLock.Unlock()
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
	Stop()
	Peers() map[peer.ID]peer.AddrInfo
	Messages() map[string]time.Time
}

// Service is a P2P service
// host is the libp2p host
// peers is a map of discovered peers
// peersLock is a mutex for peers
// ctx is the service context
// cancel is the cancel function for the service context
// messages is a map of received messages
// messagesLock is a mutex for messages
type Service struct {
	host         host.Host
	peers        map[peer.ID]peer.AddrInfo
	peersLock    sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
	messages     map[string]time.Time
	messagesLock sync.RWMutex
}

// NewP2PService creates a new P2P service
func NewP2PService() (*Service, error) {
	ctx, cancel := context.WithCancel(context.Background())

	p2pHost, err := libp2p.New()
	if err != nil {
		cancel()
		return nil, err
	}
	service := &Service{
		host:     p2pHost,
		peers:    make(map[peer.ID]peer.AddrInfo),
		ctx:      ctx,
		cancel:   cancel,
		messages: make(map[string]time.Time),
	}
	// Start the cleanup goroutine
	go service.cleanupMessages()

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

	if s.isDuplicateMessage(msg.Hash) {
		logrus.WithField("hash", msg.Hash).Trace("Duplicate message received")
		return
	}

	s.messagesLock.Lock()
	s.messages[msg.Hash] = time.Now()
	s.messagesLock.Unlock()

	switch msg.Type {
	case types.P2PMessageTypeDeployContainer.String():
	// TODO: Handle the deployment request
	default:
		logrus.Warnf("unknown message type: %s", msg.Type)
	}
}

// cleanupMessages periodically removes old message hashes
func (s *Service) cleanupMessages() {
	ticker := time.NewTicker(CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.messagesLock.Lock()
			for hash, timestamp := range s.messages {
				if time.Since(timestamp) > MessageTTL {
					delete(s.messages, hash)
				}
			}
			s.messagesLock.Unlock()
		case <-s.ctx.Done():
			return
		}
	}
}

// Stop stops the P2P service
func (s *Service) Stop() {
	logrus.Trace("Stopping P2P Service")
	s.cancel()
	s.host.Close()
}

// Peers returns the discovered peers
func (s *Service) Peers() map[peer.ID]peer.AddrInfo {
	s.peersLock.RLock()
	defer s.peersLock.RUnlock()
	return s.peers
}

// Messages returns the received messages
func (s *Service) Messages() map[string]time.Time {
	s.messagesLock.RLock()
	defer s.messagesLock.RUnlock()
	return s.messages
}

// isDuplicateMessage checks if a message hash is a duplicate
func (s *Service) isDuplicateMessage(hash string) bool {
	s.messagesLock.RLock()
	_, ok := s.messages[hash]
	s.messagesLock.RUnlock()
	return ok
}
