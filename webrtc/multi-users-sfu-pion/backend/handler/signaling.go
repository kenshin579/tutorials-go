package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/kenshin579/tutorials-go/webrtc/multi-users-sfu-pion/backend/room"
	"github.com/kenshin579/tutorials-go/webrtc/multi-users-sfu-pion/backend/sfu"
	"github.com/labstack/echo/v4"
	"github.com/pion/webrtc/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type SignalingMessage struct {
	Type       string          `json:"type"`
	SenderID   string          `json:"senderId,omitempty"`
	SenderName string          `json:"senderName,omitempty"`
	Message    string          `json:"message,omitempty"`
	Payload    json.RawMessage `json:"payload,omitempty"`
}

type Signaling struct {
	rm *room.Manager
}

func NewSignaling(rm *room.Manager) *Signaling {
	return &Signaling{rm: rm}
}

func (s *Signaling) HandleWebSocket(c echo.Context) error {
	roomID := c.QueryParam("roomId")
	if roomID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "roomId is required")
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	peerID := uuid.New().String()[:8]
	peer, err := sfu.NewPeer(peerID, ws)
	if err != nil {
		log.Printf("[ROOM:%s] failed to create peer: %v", roomID, err)
		return nil
	}
	defer peer.Close()

	if !s.rm.Join(roomID, peer) {
		peer.SendJSON(SignalingMessage{Type: "error", Message: "room is full"})
		return nil
	}
	defer s.cleanupPeer(roomID, peer)

	s.setupTrackForwarding(roomID, peer)

	// Notify existing peers about the new peer
	s.rm.Broadcast(roomID, peerID, SignalingMessage{
		Type:     "join",
		SenderID: peerID,
	})

	// Message loop
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("[ROOM:%s] peer %s read error: %v", roomID, peerID, err)
			break
		}

		var sm SignalingMessage
		if err := json.Unmarshal(msg, &sm); err != nil {
			log.Printf("[ROOM:%s] invalid message: %v", roomID, err)
			continue
		}

		sm.SenderID = peerID
		log.Printf("[ROOM:%s] peer %s -> %s", roomID, peerID, sm.Type)

		switch sm.Type {
		case "offer":
			s.handleOffer(roomID, peer, sm)
		case "answer":
			s.handleAnswer(peer, sm)
		case "ice":
			s.handleICE(peer, sm)
		case "chat":
			s.rm.Broadcast(roomID, peerID, sm)
		}
	}

	return nil
}

func (s *Signaling) handleOffer(roomID string, peer *sfu.Peer, msg SignalingMessage) {
	var offer webrtc.SessionDescription
	if err := json.Unmarshal(msg.Payload, &offer); err != nil {
		log.Printf("[PEER:%s] invalid offer: %v", peer.ID, err)
		return
	}

	if err := peer.PC.SetRemoteDescription(offer); err != nil {
		log.Printf("[PEER:%s] SetRemoteDescription error: %v", peer.ID, err)
		return
	}

	answer, err := peer.PC.CreateAnswer(nil)
	if err != nil {
		log.Printf("[PEER:%s] CreateAnswer error: %v", peer.ID, err)
		return
	}

	if err := peer.PC.SetLocalDescription(answer); err != nil {
		log.Printf("[PEER:%s] SetLocalDescription error: %v", peer.ID, err)
		return
	}

	peer.SendJSON(SignalingMessage{
		Type:    "answer",
		Payload: mustMarshal(answer),
	})

	// After initial offer/answer, add existing tracks from other peers
	s.addExistingTracks(roomID, peer)
}

func (s *Signaling) handleAnswer(peer *sfu.Peer, msg SignalingMessage) {
	var answer webrtc.SessionDescription
	if err := json.Unmarshal(msg.Payload, &answer); err != nil {
		log.Printf("[PEER:%s] invalid answer: %v", peer.ID, err)
		return
	}

	if err := peer.PC.SetRemoteDescription(answer); err != nil {
		log.Printf("[PEER:%s] SetRemoteDescription error: %v", peer.ID, err)
	}
}

func (s *Signaling) handleICE(peer *sfu.Peer, msg SignalingMessage) {
	var candidate webrtc.ICECandidateInit
	if err := json.Unmarshal(msg.Payload, &candidate); err != nil {
		log.Printf("[PEER:%s] invalid ICE candidate: %v", peer.ID, err)
		return
	}

	if err := peer.PC.AddICECandidate(candidate); err != nil {
		log.Printf("[PEER:%s] AddICECandidate error: %v", peer.ID, err)
	}
}

// setupTrackForwarding registers OnTrack to forward RTP packets from this peer to others.
func (s *Signaling) setupTrackForwarding(roomID string, peer *sfu.Peer) {
	peer.PC.OnTrack(func(remoteTrack *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		log.Printf("[PEER:%s] received track: %s (kind: %s)", peer.ID, remoteTrack.Codec().MimeType, remoteTrack.Kind())

		localTrack, err := webrtc.NewTrackLocalStaticRTP(
			remoteTrack.Codec().RTPCodecCapability,
			remoteTrack.ID(),
			peer.ID,
		)
		if err != nil {
			log.Printf("[PEER:%s] error creating local track: %v", peer.ID, err)
			return
		}

		peer.AddLocalTrack(localTrack)

		// Add this track to all other peers in the room
		r := s.rm.GetRoom(roomID)
		if r != nil {
			for _, otherPeer := range r.GetOtherPeers(peer.ID) {
				if err := otherPeer.AddRemoteTrack(localTrack); err != nil {
					log.Printf("[PEER:%s] error adding track to %s: %v", peer.ID, otherPeer.ID, err)
					continue
				}
				s.renegotiate(otherPeer)
			}
		}

		// Forward RTP packets
		buf := make([]byte, 1500)
		for {
			n, _, readErr := remoteTrack.Read(buf)
			if readErr != nil {
				log.Printf("[PEER:%s] track read ended: %v", peer.ID, readErr)
				return
			}
			if _, writeErr := localTrack.Write(buf[:n]); writeErr != nil {
				return
			}
		}
	})
}

// addExistingTracks adds tracks from other peers that are already streaming.
func (s *Signaling) addExistingTracks(roomID string, peer *sfu.Peer) {
	r := s.rm.GetRoom(roomID)
	if r == nil {
		return
	}

	tracksAdded := false
	for _, otherPeer := range r.GetOtherPeers(peer.ID) {
		for _, track := range otherPeer.GetLocalTracks() {
			if err := peer.AddRemoteTrack(track); err != nil {
				log.Printf("[PEER:%s] error adding existing track: %v", peer.ID, err)
				continue
			}
			tracksAdded = true
		}
	}

	if tracksAdded {
		s.renegotiate(peer)
	}
}

// renegotiate sends a new offer to the peer when tracks have been added/removed.
func (s *Signaling) renegotiate(peer *sfu.Peer) {
	if peer.PC.SignalingState() != webrtc.SignalingStateStable {
		log.Printf("[PEER:%s] skipping renegotiation, signaling state: %s", peer.ID, peer.PC.SignalingState())
		return
	}

	offer, err := peer.PC.CreateOffer(nil)
	if err != nil {
		log.Printf("[PEER:%s] CreateOffer error: %v", peer.ID, err)
		return
	}

	if err := peer.PC.SetLocalDescription(offer); err != nil {
		log.Printf("[PEER:%s] SetLocalDescription error: %v", peer.ID, err)
		return
	}

	peer.SendJSON(SignalingMessage{
		Type:    "offer",
		Payload: mustMarshal(offer),
	})
	log.Printf("[PEER:%s] renegotiation offer sent", peer.ID)
}

func (s *Signaling) cleanupPeer(roomID string, peer *sfu.Peer) {
	// Remove this peer's tracks from other peers
	r := s.rm.GetRoom(roomID)
	if r != nil {
		localTracks := peer.GetLocalTracks()
		for _, otherPeer := range r.GetOtherPeers(peer.ID) {
			for _, track := range localTracks {
				if err := otherPeer.RemoveTrack(track.ID()); err != nil {
					log.Printf("[PEER:%s] error removing track from %s: %v", peer.ID, otherPeer.ID, err)
				}
			}
			s.renegotiate(otherPeer)
		}
	}

	s.rm.Leave(roomID, peer.ID)

	// Notify remaining peers
	s.rm.Broadcast(roomID, peer.ID, SignalingMessage{
		Type:     "leave",
		SenderID: peer.ID,
	})
}

func mustMarshal(v any) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}
