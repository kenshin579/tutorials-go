# Implementation: Echo + WebRTC + React 1:1 화상 통화

## 프로젝트 구조

```
webrtc/
├── backend/
│   ├── main.go                  # Echo 서버 엔트리포인트
│   ├── go.mod
│   ├── handler/
│   │   └── signaling.go         # WebSocket 핸들러
│   └── room/
│       └── manager.go           # Room 관리 (roomId → 클라이언트 매핑)
└── frontend/
    ├── package.json
    ├── vite.config.ts
    ├── tsconfig.json
    ├── index.html
    └── src/
        ├── main.tsx
        ├── App.tsx
        ├── hooks/
        │   ├── useSignaling.ts  # WebSocket 연결 및 메시지 송수신
        │   └── useWebRTC.ts     # RTCPeerConnection 관리
        └── components/
            ├── VideoPanel.tsx    # Local/Remote 비디오 표시
            └── ChatPanel.tsx    # DataChannel 채팅 UI
```

---

## Backend 구현

### 1. Echo 서버 설정 (main.go)

- Echo v4 인스턴스 생성
- CORS 미들웨어 설정 (localhost:5173 허용)
- Logger 미들웨어 설정
- `GET /ws?roomId={roomId}` 라우트 등록
- 포트: `:8080`

### 2. Room Manager (room/manager.go)

```go
type RoomManager struct {
    mu    sync.RWMutex
    rooms map[string][]*websocket.Conn // roomId → connections (최대 2개)
}
```

- `Join(roomId, conn)`: 방에 클라이언트 추가, 2명 초과 시 거부
- `Leave(roomId, conn)`: 방에서 클라이언트 제거
- `GetPeer(roomId, conn)`: 같은 방의 상대방 연결 반환
- thread-safe 처리 (sync.RWMutex)

### 3. WebSocket Signaling Handler (handler/signaling.go)

- Gorilla WebSocket upgrader 설정
- 연결 시 roomId 파라미터 읽기 → RoomManager.Join 호출
- 메시지 수신 루프:
  - JSON 파싱: `{"type": "offer|answer|ice", "payload": {}}`
  - GetPeer로 상대방 조회 → 메시지 그대로 relay
- 연결 종료 시 RoomManager.Leave 호출
- 모든 Signaling 이벤트 로그 출력

### Signaling 메시지 흐름

```
Client A                  Server                  Client B
   │                        │                        │
   ├── join(roomId) ───────>│                        │
   │                        │<────── join(roomId) ───┤
   │                        │                        │
   ├── offer ──────────────>│── relay offer ────────>│
   │                        │                        │
   │<────── relay answer ───│<────── answer ─────────┤
   │                        │                        │
   ├── ice candidate ──────>│── relay ice ──────────>│
   │<────── relay ice ──────│<────── ice candidate ──┤
   │                        │                        │
   │<═══════ WebRTC P2P Media/Data ════════════════>│
```

---

## Frontend 구현

### 1. 프로젝트 설정

- Vite + React 18 + TypeScript
- 외부 라이브러리 없음 (WebRTC/WebSocket은 브라우저 네이티브 API)
- `vite.config.ts`: dev server 포트 5173

### 2. useSignaling Hook

```typescript
interface SignalingMessage {
  type: 'offer' | 'answer' | 'ice';
  payload: any;
}
```

- `useSignaling(roomId: string)`
- WebSocket 연결 관리 (`ws://localhost:8080/ws?roomId=xxx`)
- `send(message: SignalingMessage)`: 메시지 전송
- `onMessage` 콜백 등록: offer/answer/ice 수신 처리
- 연결 상태 관리 (connected/disconnected)
- cleanup: 컴포넌트 unmount 시 WebSocket close

### 3. useWebRTC Hook

- `useWebRTC(signaling)`: signaling hook과 연동
- `RTCPeerConnection`을 `useRef`로 관리 (re-render 방지)
- STUN 서버 설정: `stun:stun.l.google.com:19302`
- 핵심 로직:
  - `getUserMedia()` → localStream 획득 → localVideoRef에 연결
  - `addTrack()` → localStream 트랙을 PeerConnection에 추가
  - `ontrack` → remoteStream 수신 → remoteVideoRef에 연결
  - `onicecandidate` → ICE candidate를 signaling으로 전송
  - `createOffer()` / `createAnswer()` / `setRemoteDescription()` 처리
- DataChannel:
  - Caller: `createDataChannel('chat')` 생성
  - Callee: `ondatachannel` 이벤트로 수신
  - `onmessage` → 채팅 메시지 수신
  - `send()` → 채팅 메시지 전송

### 4. App.tsx

- roomId 입력 UI (간단한 input + 버튼)
- roomId 입력 후 → 통화 화면 전환
- VideoPanel + ChatPanel 렌더링

### 5. VideoPanel 컴포넌트

- `<video>` 태그 2개 (local, remote)
- local video: muted, autoplay
- remote video: autoplay
- 간단한 레이아웃 (나란히 배치)

### 6. ChatPanel 컴포넌트

- 메시지 목록 (스크롤 영역)
- 입력창 + 전송 버튼
- DataChannel `send()` / `onmessage` 연동

---

## 로깅 요구사항

### Backend 로그
- WebSocket 연결/해제: `[ROOM:{roomId}] client joined/left`
- 메시지 relay: `[ROOM:{roomId}] relay {type}`

### Frontend 로그 (console)
- ICE connection state 변경
- Signaling 메시지 송수신
- DataChannel open/close

---

## 실행 방법

```bash
# Backend
cd webrtc/backend
go mod tidy
go run main.go
# → http://localhost:8080

# Frontend
cd webrtc/frontend
npm install
npm run dev
# → http://localhost:5173

# 테스트
# Chrome 탭 2개에서 동일한 roomId로 접속
```
