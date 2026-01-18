package device

import (
	"encoding/json"
	"math/rand"
	"sync"
	"time"
)

// State represents the device state
type State struct {
	DeviceID    string  `json:"deviceId"`
	Status      string  `json:"status"`
	Temperature float64 `json:"temperature"`
	Timestamp   int64   `json:"timestamp"`
}

// Simulator represents a virtual device
type Simulator struct {
	mu     sync.RWMutex
	status string
}

// NewSimulator creates a new device simulator
func NewSimulator() *Simulator {
	return &Simulator{status: "idle"}
}

// GetState returns the current device state
func (s *Simulator) GetState() State {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return State{
		DeviceID:    "1",
		Status:      s.status,
		Temperature: 35.0 + rand.Float64()*5.0,
		Timestamp:   time.Now().Unix(),
	}
}

// HandleCommand handles incoming commands
func (s *Simulator) HandleCommand(action string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	switch action {
	case "start":
		s.status = "running"
	case "stop":
		s.status = "idle"
	}
}

// IsRunning returns true if the device is running
func (s *Simulator) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status == "running"
}

// ToJSON returns the state as JSON
func (s *Simulator) ToJSON() ([]byte, error) {
	return json.Marshal(s.GetState())
}
