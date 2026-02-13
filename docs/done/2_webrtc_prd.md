# PRD: Pion SFU 기반 다자간 화상 통화 확장 프로젝트

## 1. 개요 (Overview)

본 프로젝트는 기존 1:1 P2P WebRTC 화상 통화 프로젝트(`webrtc/simple-p2p/`)와 별도로,
**Pion WebRTC 라이브러리 기반 SFU(Selective Forwarding Unit) 서버**를 구현하고
다자간 화상 통화를 지원하는 애플리케이션을 `webrtc/multi-users-sfu/` 폴더에 독립적으로 개발한다.

> 블로그 게시용 프로젝트이므로 기존 P2P(`webrtc/simple-p2p/`)와 SFU(`webrtc/multi-users-sfu/`)를 별도 폴더로 분리하여,
> 각각 독립적인 예제로 실행 가능하도록 구성한다.

**목적:**
- SFU 아키텍처의 원리와 P2P 대비 장점을 실습으로 이해
- Pion WebRTC 라이브러리를 활용한 서버 측 미디어 라우팅 구현
- 기존 코드 구조(Echo + React Hooks)를 참고하되, 독립 실행 가능한 프로젝트로 구성

---

## 2. 목표 (Goals)

### 기능적 목표
- Pion SFU 서버를 통한 다자간(N:N) 화상 통화 구현
- WebSocket 브로드캐스트 기반 다자간 텍스트 채팅 구현
- Room 단위 참가자 관리 및 동적 입퇴장 처리

### 학습 목표
- SFU 아키텍처에서의 미디어 흐름(Publish/Subscribe) 이해
- Pion WebRTC의 RTP 패킷 포워딩 메커니즘 이해
- 다자간 통화에서의 PeerConnection 관리 전략 이해
- P2P DataChannel vs WebSocket 브로드캐스트의 트레이드오프 이해

---

## 3. 범위 (Scope)

### 포함 (In Scope)
- `webrtc/multi-users-sfu/` 폴더에 독립 프로젝트로 구성 (기존 `webrtc/simple-p2p/`와 분리)
- 동일 Room 내 최대 4~6명 다자간 통화
- Pion WebRTC 기반 SFU 서버 구현
- STUN 서버 사용 (Google Public STUN)
- WebSocket 기반 Signaling 및 채팅
- 로컬 개발 환경(macOS) 기준
- 인증 없이 roomId 기반 연결
- 참가자 입퇴장 시 동적 UI 업데이트

### 제외 (Out of Scope)
- TURN 서버 구성
- Simulcast / SVC 적용
- 녹화 및 화면 공유
- 사용자 인증 / 권한 관리
- 모바일 대응
- 배포 환경 구성

---

## 4. 전체 아키텍처

### P2P vs SFU 비교

```
[ P2P (기존 1:1) ]               [ SFU (다자간) ]

Client A ←──P2P──→ Client B      Client A ──┐
                                  Client B ──┤── Pion SFU Server
                                  Client C ──┘
                                  (서버가 RTP 패킷을 선택적으로 포워딩)
```

### 프로젝트 폴더 구조

```
tutorials-go/
└── webrtc/
    ├── simple-p2p/            # 기존 1:1 P2P 화상 통화 (블로그 1편)
    │   ├── backend/
    │   └── frontend/
    └── multi-users-sfu/       # SFU 다자간 화상 통화 (블로그 2편) ← 이번 프로젝트
        ├── backend/
        └── frontend/
```

> 두 프로젝트는 독립적으로 실행 가능하며, 각각 별도의 `go.mod`와 `package.json`을 가진다.

### 구성 요소
- **Frontend**: React 18 + Vite + TypeScript (WebRTC Peer 역할)
- **Backend**: Go + Echo (WebSocket Signaling) + Pion WebRTC (SFU 미디어 라우팅)
- **Media**: 클라이언트 → SFU 서버 → 다른 클라이언트들로 RTP 포워딩

### 아키텍처 다이어그램

```
                     ┌─────────────────────┐
                     │   Go Echo Server    │
                     │                     │
                     │  ┌───────────────┐  │
React Client A ──WS──┤  │  Signaling    │  │
                     │  │  (WebSocket)  │  │
React Client B ──WS──┤  └───────────────┘  │
                     │                     │
React Client C ──WS──┤  ┌───────────────┐  │
                     │  │  Pion SFU     │  │
                     │  │  (RTP Relay)  │  │
                     │  └───────────────┘  │
                     └─────────────────────┘

Media Flow:
  Client A ──(Upload Track)──→ SFU ──(Forward RTP)──→ Client B, C
  Client B ──(Upload Track)──→ SFU ──(Forward RTP)──→ Client A, C
  Client C ──(Upload Track)──→ SFU ──(Forward RTP)──→ Client A, B
```

---

## 5. 기능 요구사항 (Functional Requirements)

### 5.1 Room 연결
- 사용자는 roomId를 기준으로 방에 입장한다
- 동일 Room에 최대 4~6명까지 입장 가능하다
- 새로운 참가자 입장 시 기존 참가자들에게 알림을 전송한다
- 참가자 퇴장 시 해당 트랙을 정리하고 UI를 업데이트한다

### 5.2 WebRTC Signaling
- Signaling은 WebSocket을 통해 수행한다
- SFU 구조에서 각 클라이언트는 SFU 서버와 1개의 PeerConnection을 맺는다
- 지원 메시지 타입: `offer`, `answer`, `ice`, `join`, `leave`

**메시지 형식:**
```json
{
  "type": "offer | answer | ice | join | leave",
  "senderId": "client-uuid",
  "payload": {}
}
```

### 5.3 SFU 미디어 라우팅
- 클라이언트는 자신의 미디어 트랙을 SFU 서버로 전송한다 (Publish)
- SFU 서버는 수신한 RTP 패킷을 Room 내 다른 클라이언트들에게 포워딩한다 (Subscribe)
- 미디어 디코딩/인코딩 없이 RTP 패킷을 그대로 전달한다 (낮은 서버 부하)

### 5.4 채팅
- 다자간 채팅은 WebSocket 브로드캐스트 방식으로 구현한다
- DataChannel 대신 WebSocket을 사용하는 이유: SFU 구조에서는 서버가 중계하므로 WebSocket이 자연스러움
- 채팅 메시지는 Room 내 모든 참가자에게 실시간 전달된다

**채팅 메시지 형식:**
```json
{
  "type": "chat",
  "senderId": "client-uuid",
  "senderName": "User1",
  "message": "안녕하세요!"
}
```

---

## 6. Frontend 요구사항 (React)

### 기술 스택
- React 18
- TypeScript
- Vite
- Web API (WebRTC, MediaDevices)

### 구조 요구사항
- 기존 `webrtc/simple-p2p/frontend`의 Hook 패턴을 참고하여 새로 작성한다
  - `useWebRTC`: SFU 서버와의 PeerConnection 관리 (1개의 PC로 Publish/Subscribe)
  - `useSignaling`: WebSocket Signaling + 채팅 메시지 관리
- `RTCPeerConnection` 객체는 `useRef`로 관리한다
- 새로운 참가자의 Remote Track 수신 시 동적으로 비디오 요소를 추가한다

### 주요 UI
- Local Video 영역 (자신의 카메라)
- Remote Video 그리드 (다른 참가자들의 영상을 격자로 배치)
- 참가자 수에 따라 그리드 레이아웃 동적 조정 (2명: 1x2, 3~4명: 2x2 등)
- 채팅 입력창 및 메시지 목록 (발신자 이름 표시)
- 참가자 목록 표시

### UI 와이어프레임

#### 화면 1: Room 입장 화면

```
┌──────────────────────────────────────────────┐
│                                              │
│                                              │
│          WebRTC 다자간 화상 통화              │
│                                              │
│       Room ID를 입력하고 입장하세요           │
│                                              │
│       ┌──────────────────┐  ┌──────┐         │
│       │ room-123         │  │ 입장 │         │
│       └──────────────────┘  └──────┘         │
│                                              │
│                                              │
└──────────────────────────────────────────────┘
```

#### 화면 2: 통화 화면 (참가자 4명 기준)

```
┌──────────────────────────────────────────────────────────────┐
│ WebRTC 다자간 화상 통화    Room: room-123    ● Connected     │
├──────────────────────────────────────────────────┬───────────┤
│                                                  │  Chat     │
│  ┌─────────────────────┐ ┌─────────────────────┐ ├───────────┤
│  │                     │ │                     │ │           │
│  │                     │ │                     │ │ User2:    │
│  │     나 (Local)      │ │     User2           │ │ 안녕!     │
│  │                     │ │                     │ │           │
│  │                     │ │                     │ │ User3:    │
│  └─────────────────────┘ └─────────────────────┘ │ 반갑습니다│
│                                                  │           │
│  ┌─────────────────────┐ ┌─────────────────────┐ │ 나:       │
│  │                     │ │                     │ │ ㅎㅇ      │
│  │                     │ │                     │ │           │
│  │     User3           │ │     User4           │ │           │
│  │                     │ │                     │ │           │
│  │                     │ │                     │ │           │
│  └─────────────────────┘ └─────────────────────┘ ├───────────┤
│                                                  │┌─────┐┌──┐│
│                                                  ││메시지││▶ ││
│                                                  │└─────┘└──┘│
└──────────────────────────────────────────────────┴───────────┘
```

#### 그리드 레이아웃 규칙

```
참가자 1명 (나만)        참가자 2명              참가자 3~4명
┌───────────────┐    ┌───────┬───────┐    ┌───────┬───────┐
│               │    │       │       │    │       │       │
│    나 (Local) │    │  나   │ User2 │    │  나   │ User2 │
│               │    │       │       │    │       │       │
└───────────────┘    └───────┴───────┘    ├───────┼───────┤
                                          │       │       │
                                          │ User3 │ User4 │
                                          │       │       │
                                          └───────┴───────┘

참가자 5~6명
┌───────┬───────┬───────┐
│       │       │       │
│  나   │ User2 │ User3 │
│       │       │       │
├───────┼───────┼───────┤
│       │       │       │
│ User4 │ User5 │ User6 │
│       │       │       │
└───────┴───────┴───────┘
```

---

## 7. Backend 요구사항 (Go + Echo + Pion)

### 기술 스택
- Go 1.21+
- Echo framework
- Gorilla WebSocket
- **Pion WebRTC** (`github.com/pion/webrtc/v4`)

### 역할
- WebSocket 연결 관리 및 Signaling 메시지 중계
- roomId 기준 참가자 그룹 관리
- **Pion WebRTC를 통한 SFU 미디어 라우팅**
  - 각 클라이언트의 업로드 트랙(Track) 수신
  - Room 내 다른 클라이언트들에게 RTP 패킷 포워딩
- 참가자 입퇴장 이벤트 처리
- 채팅 메시지 브로드캐스트

### 주요 컴포넌트

```
webrtc/multi-users-sfu/
├── backend/
│   ├── main.go                # 서버 부트스트랩
│   ├── handler/
│   │   └── signaling.go       # WebSocket Signaling 핸들러
│   ├── room/
│   │   └── manager.go         # Room 및 참가자 관리
│   └── sfu/
│       ├── peer.go            # 개별 Peer 관리 (PeerConnection + Tracks)
│       └── router.go          # RTP 패킷 포워딩 로직
└── frontend/
    ├── src/
    │   ├── App.tsx
    │   ├── components/
    │   │   ├── VideoGrid.tsx   # 다자간 비디오 그리드
    │   │   └── ChatPanel.tsx   # 채팅 패널
    │   └── hooks/
    │       ├── useWebRTC.ts    # SFU PeerConnection 관리
    │       └── useSignaling.ts # WebSocket Signaling + 채팅
    └── package.json
```

### API
| Method | Endpoint | 설명 |
|--------|----------|------|
| GET | `/ws?roomId={roomId}` | WebSocket 업그레이드, Signaling + 채팅 통신 채널 |

---

## 8. SFU 핵심 로직 (Pion WebRTC)

### 8.1 PeerConnection 생명주기

```
1. 클라이언트 입장 → WebSocket 연결
2. SFU 서버가 PeerConnection 생성
3. 클라이언트가 Offer 전송 → SFU가 Answer 응답
4. ICE Candidate 교환
5. 클라이언트가 미디어 트랙 전송 (addTrack)
6. SFU가 OnTrack 콜백에서 RTP 패킷 수신
7. 수신한 RTP 패킷을 Room 내 다른 PeerConnection에 WriteRTP로 전달
8. 클라이언트 퇴장 → PeerConnection 정리 및 트랙 제거
```

### 8.2 RTP 포워딩 핵심 로직 (의사코드)

```go
// 새 트랙이 수신되면
peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
    // Room 내 다른 모든 Peer에게 LocalTrack 생성
    localTrack, _ := webrtc.NewTrackLocalStaticRTP(...)

    // 다른 Peer의 PeerConnection에 트랙 추가
    for _, otherPeer := range room.Peers {
        otherPeer.PC.AddTrack(localTrack)
    }

    // RTP 패킷을 지속적으로 읽어서 포워딩
    for {
        rtpPacket, _ := track.ReadRTP()
        localTrack.WriteRTP(rtpPacket)
    }
})
```

---

## 9. 비기능 요구사항 (Non-Functional)

- 코드 가독성과 학습 목적을 우선한다
- 기존 `webrtc/simple-p2p/` 프로젝트의 코드 구조와 패턴을 참고하되 독립 프로젝트로 구성한다
- `webrtc/multi-users-sfu/` 폴더만으로 단독 실행이 가능해야 한다 (별도 `go.mod`, `package.json`)
- 최소한의 외부 의존성만 사용한다 (Pion WebRTC 추가)
- 로컬 환경에서 바로 실행 가능해야 한다
- 로그를 통해 Signaling, ICE, RTP 포워딩 흐름을 확인할 수 있어야 한다
- 참가자 입퇴장 시 리소스(PeerConnection, Track)가 올바르게 정리되어야 한다

---

## 10. 개발 환경

| 항목 | 버전 |
|------|------|
| OS | macOS |
| Browser | Chrome |
| Node.js | LTS |
| Go | 1.21+ |
| Pion WebRTC | v4 |

---

## 11. 기존 프로젝트와의 비교

> 두 프로젝트는 `webrtc/` 하위에 별도 폴더로 독립적으로 존재한다.

| 항목 | `webrtc/simple-p2p/` (1:1 P2P) | `webrtc/multi-users-sfu/` (SFU 다자간) |
|------|--------------------------------|----------------------------------------|
| 프로젝트 폴더 | `webrtc/simple-p2p/` | `webrtc/multi-users-sfu/` |
| 미디어 전송 | 클라이언트 간 직접 P2P | 클라이언트 → SFU → 클라이언트 |
| PeerConnection | 클라이언트 간 1개 | 클라이언트-SFU 간 1개 |
| 최대 참가자 | 2명 | 4~6명 |
| 채팅 방식 | DataChannel (P2P) | WebSocket 브로드캐스트 |
| 서버 역할 | Signaling만 | Signaling + 미디어 라우팅 |
| 프론트엔드 비디오 | Local + Remote 1개 | Local + Remote 그리드 |
| 백엔드 의존성 | Echo + Gorilla WS | Echo + Gorilla WS + Pion WebRTC |
| 독립 실행 | 단독 실행 가능 | 단독 실행 가능 |

---

## 12. 성공 기준 (Success Metrics)

- [ ] 3개 이상의 브라우저 탭에서 동시에 영상/음성이 정상 송출된다
- [ ] 참가자 입장 시 기존 참가자들에게 새 비디오 스트림이 나타난다
- [ ] 참가자 퇴장 시 해당 비디오가 UI에서 제거된다
- [ ] WebSocket 채팅이 Room 내 모든 참가자에게 전달된다
- [ ] Pion SFU의 RTP 포워딩이 로그로 확인된다
- [ ] ICE Candidate 교환이 로그로 확인된다
- [ ] React Hook 구조가 기존 `webrtc/simple-p2p/` 패턴을 참고하되 독립적으로 동작한다
- [ ] `webrtc/multi-users-sfu/` 폴더만으로 단독 빌드 및 실행이 가능하다

---

## 13. 확장 아이디어 (Post-Scope)

- Simulcast 적용 (네트워크 상태에 따른 품질 자동 조절)
- 화면 공유 기능
- 참가자별 음소거 / 비디오 끄기 토글
- 녹화 기능 (서버 측 RTP → 파일 저장)
- TURN 서버 추가
- JWT 기반 Room 접근 제어
