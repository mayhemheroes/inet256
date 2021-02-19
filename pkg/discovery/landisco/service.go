package landisco

import (
	"context"
	"encoding/json"
	"net"
	"sync"
	"time"

	"github.com/brendoncarroll/go-p2p"
	"github.com/sirupsen/logrus"
)

const multicastAddr = "[ff02::1]:25600"

type service struct {
	conn *net.UDPConn
	cf   context.CancelFunc

	mu         sync.RWMutex
	lookingFor map[p2p.PeerID][]string
}

func New() (p2p.DiscoveryService, error) {
	gaddr, err := net.ResolveUDPAddr("udp6", multicastAddr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenMulticastUDP("udp6", nil, gaddr)
	if err != nil {
		return nil, err
	}
	ctx, cf := context.WithCancel(context.Background())
	s := &service{
		lookingFor: make(map[p2p.PeerID][]string),
		cf:         cf,
		conn:       conn,
	}
	go s.run(ctx)
	return s, nil
}

func (s *service) run(ctx context.Context) error {
	buf := make([]byte, 1<<16)
	for {
		n, err := s.conn.Read(buf[:])
		if err != nil {
			return err
		}
		if err := s.handleMessage(buf[:n]); err != nil {
			logrus.Error(err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
	}
}

func (s *service) handleMessage(buf []byte) error {
	msg, err := ParseMessage(buf)
	if err != nil {
		return err
	}
	s.mu.RLock()
	var peers []p2p.PeerID
	for id := range s.lookingFor {
		peers = append(peers, id)
	}
	s.mu.RUnlock()
	i, data, err := UnpackMessage(msg, peers)
	if err != nil {
		return err
	}
	peer := peers[i]
	adv := Advertisement{}
	if err := json.Unmarshal(data, &adv); err != nil {
		return err
	}
	s.mu.Lock()
	if _, exists := s.lookingFor[peer]; exists {
		s.lookingFor[peer] = adv.Transports
	}
	s.mu.Unlock()
	return nil
}

func (s *service) Find(ctx context.Context, id p2p.PeerID) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	addrs, exists := s.lookingFor[id]
	if exists {
		return addrs, nil
	}
	s.lookingFor[id] = []string{}
	return nil, nil
}

func (s *service) Announce(ctx context.Context, id p2p.PeerID, xs []string, ttl time.Duration) error {
	adv := Advertisement{
		Transports: xs,
	}
	data, err := json.Marshal(adv)
	if err != nil {
		return err
	}
	msg := NewMessage(id, data)
	_, err = s.conn.Write(msg)
	return err
}

func (s *service) Close() error {
	s.cf()
	return s.conn.Close()
}
