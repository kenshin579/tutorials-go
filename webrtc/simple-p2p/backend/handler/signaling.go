package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kenshin579/tutorials-go/webrtc/simple-p2p/backend/room"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type SignalingMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
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

	if !s.rm.Join(roomID, ws) {
		ws.WriteMessage(websocket.TextMessage, []byte(`{"type":"error","payload":"room is full"}`))
		return nil
	}
	defer s.rm.Leave(roomID, ws)

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Printf("[ROOM:%s] read error: %v", roomID, err)
			break
		}

		var sm SignalingMessage
		if err := json.Unmarshal(msg, &sm); err != nil {
			log.Printf("[ROOM:%s] invalid message: %v", roomID, err)
			continue
		}

		log.Printf("[ROOM:%s] relay %s", roomID, sm.Type)

		peer := s.rm.GetPeer(roomID, ws)
		if peer != nil {
			if err := peer.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Printf("[ROOM:%s] relay error: %v", roomID, err)
			}
		}
	}

	return nil
}
