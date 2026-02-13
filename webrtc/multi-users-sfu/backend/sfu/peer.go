package sfu

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v4"
)

var iceServers = []webrtc.ICEServer{
	{URLs: []string{"stun:stun.l.google.com:19302"}},
}

// Peer wraps a WebSocket connection and a Pion PeerConnection for a single client.
type Peer struct {
	ID   string
	Conn *websocket.Conn
	PC   *webrtc.PeerConnection

	mu          sync.Mutex
	localTracks []*webrtc.TrackLocalStaticRTP // tracks published BY this peer (forwarded to others)
	senders     map[string]*webrtc.RTPSender  // tracks added TO this peer's PC from other peers
}

func NewPeer(id string, conn *websocket.Conn) (*Peer, error) {
	pc, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: iceServers,
	})
	if err != nil {
		return nil, err
	}

	p := &Peer{
		ID:      id,
		Conn:    conn,
		PC:      pc,
		senders: make(map[string]*webrtc.RTPSender),
	}

	pc.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}
		log.Printf("[PEER:%s] ICE candidate generated", id)
		p.SendJSON(map[string]interface{}{
			"type":    "ice",
			"payload": c.ToJSON(),
		})
	})

	pc.OnConnectionStateChange(func(state webrtc.PeerConnectionState) {
		log.Printf("[PEER:%s] connection state: %s", id, state.String())
	})

	return p, nil
}

func (p *Peer) SendJSON(v interface{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.Conn.WriteJSON(v)
}

// AddRemoteTrack adds another peer's forwarded track to this peer's PeerConnection.
func (p *Peer) AddRemoteTrack(track *webrtc.TrackLocalStaticRTP) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	sender, err := p.PC.AddTrack(track)
	if err != nil {
		return err
	}
	p.senders[track.ID()] = sender

	// Drain RTCP from the sender
	go func() {
		buf := make([]byte, 1500)
		for {
			if _, _, err := sender.Read(buf); err != nil {
				return
			}
		}
	}()

	return nil
}

// RemoveTrack removes a forwarded track from this peer's PeerConnection.
func (p *Peer) RemoveTrack(trackID string) error {
	p.mu.Lock()
	sender, ok := p.senders[trackID]
	if !ok {
		p.mu.Unlock()
		return nil
	}
	delete(p.senders, trackID)
	p.mu.Unlock()

	return p.PC.RemoveTrack(sender)
}

// AddLocalTrack stores a local track that this peer is publishing.
func (p *Peer) AddLocalTrack(track *webrtc.TrackLocalStaticRTP) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.localTracks = append(p.localTracks, track)
}

// GetLocalTracks returns a copy of the local tracks this peer is publishing.
func (p *Peer) GetLocalTracks() []*webrtc.TrackLocalStaticRTP {
	p.mu.Lock()
	defer p.mu.Unlock()
	tracks := make([]*webrtc.TrackLocalStaticRTP, len(p.localTracks))
	copy(tracks, p.localTracks)
	return tracks
}

func (p *Peer) Close() {
	if p.PC != nil {
		p.PC.Close()
	}
}
