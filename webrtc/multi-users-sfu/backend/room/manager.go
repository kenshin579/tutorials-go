package room

import (
	"log"
	"sync"

	"github.com/kenshin579/tutorials-go/webrtc/multi-users-sfu/backend/sfu"
)

const MaxPeersPerRoom = 6

type Room struct {
	mu    sync.RWMutex
	Peers map[string]*sfu.Peer
}

// GetOtherPeers returns a snapshot of peers excluding the given peer ID.
func (r *Room) GetOtherPeers(excludeID string) []*sfu.Peer {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var peers []*sfu.Peer
	for id, p := range r.Peers {
		if id != excludeID {
			peers = append(peers, p)
		}
	}
	return peers
}

type Manager struct {
	mu    sync.RWMutex
	rooms map[string]*Room
}

func NewManager() *Manager {
	return &Manager{
		rooms: make(map[string]*Room),
	}
}

func (m *Manager) getOrCreateRoom(roomID string) *Room {
	r, ok := m.rooms[roomID]
	if !ok {
		r = &Room{Peers: make(map[string]*sfu.Peer)}
		m.rooms[roomID] = r
	}
	return r
}

func (m *Manager) Join(roomID string, peer *sfu.Peer) bool {
	m.mu.Lock()
	r := m.getOrCreateRoom(roomID)
	m.mu.Unlock()

	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.Peers) >= MaxPeersPerRoom {
		log.Printf("[ROOM:%s] room is full, rejecting peer %s", roomID, peer.ID)
		return false
	}

	r.Peers[peer.ID] = peer
	log.Printf("[ROOM:%s] peer %s joined (%d/%d)", roomID, peer.ID, len(r.Peers), MaxPeersPerRoom)
	return true
}

func (m *Manager) Leave(roomID, peerID string) {
	m.mu.Lock()
	r, ok := m.rooms[roomID]
	if !ok {
		m.mu.Unlock()
		return
	}
	m.mu.Unlock()

	r.mu.Lock()
	delete(r.Peers, peerID)
	remaining := len(r.Peers)
	r.mu.Unlock()

	if remaining == 0 {
		m.mu.Lock()
		delete(m.rooms, roomID)
		m.mu.Unlock()
	}

	log.Printf("[ROOM:%s] peer %s left", roomID, peerID)
}

func (m *Manager) Broadcast(roomID, excludePeerID string, msg interface{}) {
	m.mu.RLock()
	r, ok := m.rooms[roomID]
	m.mu.RUnlock()
	if !ok {
		return
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	for id, peer := range r.Peers {
		if id == excludePeerID {
			continue
		}
		if err := peer.SendJSON(msg); err != nil {
			log.Printf("[ROOM:%s] broadcast to %s error: %v", roomID, id, err)
		}
	}
}

func (m *Manager) GetRoom(roomID string) *Room {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.rooms[roomID]
}
