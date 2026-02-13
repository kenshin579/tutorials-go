# Implementation: Pion SFU 기반 다자간 화상 통화

## 1. 프로젝트 초기 설정

### 1.1 Backend (`webrtc/multi-users-sfu/backend/`)

```bash
cd webrtc/multi-users-sfu/backend
go mod init github.com/kenshin579/tutorials-go/webrtc/multi-users-sfu/backend
```

**주요 의존성:**
- `github.com/labstack/echo/v4` - HTTP 프레임워크
- `github.com/gorilla/websocket` - WebSocket
- `github.com/pion/webrtc/v4` - SFU 핵심 (RTP 수신/포워딩)

### 1.2 Frontend (`webrtc/multi-users-sfu/frontend/`)

```bash
cd webrtc/multi-users-sfu/frontend
npm create vite@latest . -- --template react-ts
```

**주요 의존성:**
- React 18 + TypeScript + Vite (기존 `webrtc/simple-p2p/frontend`과 동일 스택)

---

## 2. Backend 구현

### 2.1 디렉터리 구조

```
webrtc/multi-users-sfu/backend/
├── main.go
├── handler/
│   └── signaling.go
├── room/
│   └── manager.go
└── sfu/
    ├── peer.go
    └── router.go
```

### 2.2 `main.go` - 서버 부트스트랩

- Echo 서버 생성, CORS/Logger/Recover 미들웨어 적용
- `room.Manager` 및 `sfu.Router` 생성
- `handler.Signaling`에 의존성 주입
- `GET /ws?roomId={roomId}` 라우트 등록
- 포트 `:8080`에서 시작

> 기존 `webrtc/simple-p2p/backend/main.go`와 동일 패턴, `sfu.Router` 의존성 추가

### 2.3 `room/manager.go` - Room 및 참가자 관리

기존 P2P 버전과의 차이점:
- `rooms` 맵 값: `[]*websocket.Conn` → `map[string]*sfu.Peer` (peerID → Peer 매핑)
- 최대 인원: 2명 → 6명
- `Join()` 반환: 생성된 `peerID` (UUID)
- `Leave()` 시 해당 Peer의 트랙을 다른 Peer에서 제거
- `Broadcast()` 메서드 추가: 특정 Peer 제외 전체에게 메시지 전송

```go
type Manager struct {
    mu    sync.RWMutex
    rooms map[string]*Room
}

type Room struct {
    peers map[string]*sfu.Peer  // peerID → Peer
}

func (m *Manager) Join(roomID string, conn *websocket.Conn) (peerID string, ok bool)
func (m *Manager) Leave(roomID, peerID string)
func (m *Manager) Broadcast(roomID, excludePeerID string, msg []byte)
func (m *Manager) GetRoom(roomID string) *Room
```

### 2.4 `sfu/peer.go` - 개별 Peer 관리

각 클라이언트에 대응하는 Peer 구조체. PeerConnection과 WebSocket을 함께 관리한다.

```go
type Peer struct {
    ID     string
    Conn   *websocket.Conn
    PC     *webrtc.PeerConnection
    mu     sync.Mutex
}

func NewPeer(id string, conn *websocket.Conn) (*Peer, error)
func (p *Peer) Close()
func (p *Peer) SendJSON(v interface{}) error
```

- `NewPeer()`에서 `webrtc.NewPeerConnection()` 호출 (ICE 서버: Google STUN)
- `PC.OnICECandidate` 콜백: ICE candidate를 WebSocket으로 클라이언트에 전송
- `PC.OnConnectionStateChange` 콜백: 연결 상태 로깅
- `Close()`에서 PeerConnection 정리

### 2.5 `sfu/router.go` - RTP 포워딩 핵심 로직

SFU의 핵심. 한 Peer에서 수신한 트랙의 RTP 패킷을 Room 내 다른 Peer에게 포워딩한다.

```go
type Router struct {
    mu     sync.RWMutex
    tracks map[string]*webrtc.TrackLocalStaticRTP  // trackID → localTrack
}

func NewRouter() *Router
func (r *Router) SetupPeer(peer *Peer, room *Room)
func (r *Router) RemovePeer(peerID string, room *Room)
```

**`SetupPeer()` 핵심 흐름:**

1. `peer.PC.OnTrack()` 콜백 등록
2. OnTrack 발생 시:
   - `webrtc.NewTrackLocalStaticRTP()` 생성 (같은 코덱 정보)
   - Room 내 다른 Peer의 PC에 `AddTrack()` 호출
   - Renegotiation 필요 시 다른 Peer에게 새 Offer 전송
   - goroutine으로 `track.ReadRTP()` → `localTrack.WriteRTP()` 루프 실행

**Renegotiation 흐름:**
- 새 트랙이 추가/제거되면 다른 Peer와 SDP 재협상 필요
- `PC.OnNegotiationNeeded` 콜백에서 새 Offer 생성 → WebSocket으로 전송

### 2.6 `handler/signaling.go` - WebSocket Signaling 핸들러

기존 P2P 버전과의 차이점:
- 메시지에 `senderId` 필드 추가
- `join`/`leave` 메시지 타입 처리
- `offer`/`answer` 처리: 클라이언트 ↔ SFU 서버 간 직접 처리 (relay가 아님)
- `chat` 메시지: `Broadcast()`로 Room 전체 전달

```go
type SignalingMessage struct {
    Type     string          `json:"type"`
    SenderID string          `json:"senderId,omitempty"`
    Payload  json.RawMessage `json:"payload"`
}
```

**메시지 처리 흐름:**

| type | 처리 |
|------|------|
| `offer` | SFU가 `SetRemoteDescription` → `CreateAnswer` → 클라이언트에 answer 전송 |
| `answer` | SFU가 `SetRemoteDescription` (renegotiation 응답) |
| `ice` | SFU의 PeerConnection에 `AddICECandidate` |
| `chat` | Room 내 다른 Peer에게 브로드캐스트 |

---

## 3. Frontend 구현

### 3.1 디렉터리 구조

```
webrtc/multi-users-sfu/frontend/src/
├── App.tsx
├── main.tsx
├── components/
│   ├── VideoGrid.tsx
│   └── ChatPanel.tsx
└── hooks/
    ├── useWebRTC.ts
    └── useSignaling.ts
```

### 3.2 `hooks/useSignaling.ts`

기존 P2P 버전과의 차이점:
- 메시지 타입 확장: `offer | answer | ice | join | leave | chat`
- `senderId` 필드 추가
- 채팅 메시지도 WebSocket으로 전송/수신 (DataChannel 대신)

```typescript
export interface SignalingMessage {
  type: 'offer' | 'answer' | 'ice' | 'join' | 'leave' | 'chat';
  senderId?: string;
  senderName?: string;
  message?: string;
  payload?: unknown;
}
```

### 3.3 `hooks/useWebRTC.ts`

기존 P2P 버전과의 핵심 차이점:
- PeerConnection 1개로 SFU 서버와 연결 (P2P에서는 상대 클라이언트와 직접 연결)
- DataChannel 제거 (채팅은 WebSocket으로)
- `ontrack` 이벤트에서 Remote Stream을 Map으로 관리 (`Map<streamId, MediaStream>`)
- SFU에서 renegotiation offer 수신 시 자동으로 answer 응답

```typescript
// Remote streams를 Map으로 관리 (다자간)
const [remoteStreams, setRemoteStreams] = useState<Map<string, MediaStream>>(new Map());

// ontrack 이벤트
pc.ontrack = (event) => {
  const stream = event.streams[0];
  setRemoteStreams(prev => new Map(prev).set(stream.id, stream));
};
```

**연결 흐름:**
1. WebSocket 연결 후 `getUserMedia()`로 로컬 미디어 획득
2. PeerConnection 생성 → 로컬 트랙 추가
3. Offer 생성 → WebSocket으로 SFU에 전송
4. SFU의 Answer 수신 → `setRemoteDescription`
5. SFU에서 다른 참가자의 트랙이 추가되면 renegotiation offer 수신
6. 자동으로 Answer 생성 → 전송

### 3.4 `components/VideoGrid.tsx`

기존 `VideoPanel` 대체. 다자간 비디오 그리드를 렌더링한다.

- Local Video + Remote Streams를 격자 배치
- 참가자 수에 따라 CSS Grid 동적 조정:
  - 1명: 1x1
  - 2명: 1x2
  - 3~4명: 2x2
  - 5~6명: 2x3

```typescript
interface VideoGridProps {
  localVideoRef: RefObject<HTMLVideoElement>;
  remoteStreams: Map<string, MediaStream>;
}
```

### 3.5 `components/ChatPanel.tsx`

기존 ChatPanel 확장:
- `disabled` 조건: DataChannel → WebSocket 연결 여부로 변경
- `senderName` 표시 추가 (다자간이므로 누가 보낸 메시지인지 구분)
- 헤더: "DataChannel Chat" → "Chat"

```typescript
interface ChatMessage {
  text: string;
  mine: boolean;
  senderName?: string;
}
```

### 3.6 `App.tsx`

기존 App 구조 유지하되:
- `VideoPanel` → `VideoGrid` 교체
- DataChannel 관련 상태(`dcOpen`) 제거
- 채팅은 WebSocket 기반으로 연결
- 입장 시 자동으로 SFU와 PeerConnection 수립 (별도 "통화 시작" 버튼 불필요)

---

## 4. 메시지 프로토콜

### 4.1 Client → Server

| type | 용도 | payload |
|------|------|---------|
| `offer` | SDP Offer 전송 | `RTCSessionDescription` |
| `answer` | SDP Answer 응답 | `RTCSessionDescription` |
| `ice` | ICE Candidate 전송 | `RTCIceCandidate` |
| `chat` | 채팅 메시지 | `{ senderName, message }` |

### 4.2 Server → Client

| type | 용도 | payload |
|------|------|---------|
| `offer` | Renegotiation Offer | `RTCSessionDescription` |
| `answer` | SDP Answer 응답 | `RTCSessionDescription` |
| `ice` | ICE Candidate 전송 | `RTCIceCandidate` |
| `join` | 새 참가자 입장 알림 | `{ peerId }` |
| `leave` | 참가자 퇴장 알림 | `{ peerId }` |
| `chat` | 채팅 메시지 브로드캐스트 | `{ senderId, senderName, message }` |

---

## 5. 리소스 정리

### Peer 퇴장 시 정리 순서

1. WebSocket 연결 종료 감지
2. Room에서 Peer 제거
3. 해당 Peer의 LocalTrack을 다른 Peer의 PC에서 `RemoveTrack()`
4. 다른 Peer에게 renegotiation 트리거
5. 다른 Peer에게 `leave` 메시지 브로드캐스트
6. PeerConnection `Close()`
7. Room이 비면 Room 삭제
