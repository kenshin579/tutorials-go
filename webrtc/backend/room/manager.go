package room

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Manager struct {
	mu    sync.RWMutex
	rooms map[string][]*websocket.Conn
}

func NewManager() *Manager {
	return &Manager{
		rooms: make(map[string][]*websocket.Conn),
	}
}

func (m *Manager) Join(roomID string, conn *websocket.Conn) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	clients := m.rooms[roomID]
	if len(clients) >= 2 {
		log.Printf("[ROOM:%s] room is full, rejecting client", roomID)
		return false
	}

	m.rooms[roomID] = append(clients, conn)
	log.Printf("[ROOM:%s] client joined (%d/2)", roomID, len(m.rooms[roomID]))
	return true
}

func (m *Manager) Leave(roomID string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	clients := m.rooms[roomID]
	for i, c := range clients {
		if c == conn {
			m.rooms[roomID] = append(clients[:i], clients[i+1:]...)
			break
		}
	}

	if len(m.rooms[roomID]) == 0 {
		delete(m.rooms, roomID)
	}

	log.Printf("[ROOM:%s] client left", roomID)
}

func (m *Manager) GetPeer(roomID string, conn *websocket.Conn) *websocket.Conn {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, c := range m.rooms[roomID] {
		if c != conn {
			return c
		}
	}
	return nil
}
