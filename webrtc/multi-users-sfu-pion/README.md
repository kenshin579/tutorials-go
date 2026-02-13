# WebRTC SFU 다자간 화상 통화

Pion WebRTC 기반 SFU(Selective Forwarding Unit) 서버를 활용한 다자간 영상 통화 및 텍스트 채팅 애플리케이션입니다.

> 1:1 P2P 버전은 [webrtc/simple-p2p](../simple-p2p/)를 참고하세요.

## 기술 스택

| 구분 | 기술 |
|------|------|
| Backend | Go, Echo v4, Gorilla WebSocket, Pion WebRTC v4 |
| Frontend | React 18, TypeScript, Vite |
| 프로토콜 | WebSocket (시그널링/채팅), WebRTC (SFU 미디어) |

## P2P vs SFU 비교

| 항목 | P2P (simple-p2p) | SFU (multi-users-sfu-pion) |
|------|-------------------|------------------------|
| 최대 인원 | 2명 | 6명 |
| 미디어 경로 | 브라우저 ↔ 브라우저 | 브라우저 ↔ SFU 서버 ↔ 브라우저 |
| 채팅 방식 | DataChannel (P2P) | WebSocket Broadcast |
| 서버 역할 | 시그널링 릴레이만 | 시그널링 + RTP 패킷 포워딩 |
| PeerConnection | 클라이언트당 1개 (상대와 직접) | 클라이언트당 1개 (SFU와) |

## 프로젝트 구조

```
webrtc/multi-users-sfu-pion/
├── Makefile
├── backend/
│   ├── main.go                      # Echo 서버 설정 (CORS, WebSocket 엔드포인트)
│   ├── handler/
│   │   └── signaling.go            # WebSocket Signaling + SFU 연동
│   ├── room/
│   │   └── manager.go              # Room 관리 (최대 6명, Join/Leave/Broadcast)
│   └── sfu/
│       └── peer.go                 # Peer 구조체 (PeerConnection + 트랙 관리)
└── frontend/
    └── src/
        ├── App.tsx                  # Room 입장 화면 + 통화 화면
        ├── components/
        │   ├── VideoGrid.tsx       # 다자간 비디오 그리드 (동적 레이아웃)
        │   └── ChatPanel.tsx       # WebSocket 채팅 (발신자 이름 표시)
        └── hooks/
            ├── useSignaling.ts     # WebSocket 연결 + 채팅 메시지 관리
            └── useWebRTC.ts        # SFU PeerConnection + Renegotiation 처리
```

## 핵심 개념: SFU (Selective Forwarding Unit)

### SFU란?

SFU는 각 클라이언트에서 보낸 미디어 스트림을 **트랜스코딩 없이** 다른 클라이언트에게 그대로 전달(포워딩)하는 서버입니다.

- **P2P 방식**: N명이면 각 클라이언트가 N-1개의 연결 필요 → 대역폭 폭증
- **SFU 방식**: 각 클라이언트는 SFU 서버와 1개의 연결만 유지 → 서버가 분배

### RTP 포워딩

SFU의 핵심은 **RTP 패킷 포워딩**입니다:

1. 클라이언트 A가 영상/음성 트랙을 SFU에 전송
2. SFU가 `TrackLocalStaticRTP`로 로컬 트랙 생성
3. 이 로컬 트랙을 다른 클라이언트(B, C, ...)의 PeerConnection에 추가
4. goroutine에서 `remoteTrack.Read()` → `localTrack.Write()` 루프로 RTP 포워딩

### Renegotiation

다자간 통화에서 참가자가 입장/퇴장하면 트랙이 추가/제거됩니다. 이때 SFU가 기존 클라이언트와 **SDP 재협상(Renegotiation)** 을 수행합니다:

1. SFU가 새 Offer 생성 → 클라이언트에 전송
2. 클라이언트가 자동으로 Answer 생성 → SFU에 응답
3. 새 트랙이 `ontrack` 이벤트로 전달되어 화면에 표시

## 시스템 아키텍처

```mermaid
graph TB
    subgraph Browser A
        A_App["App.tsx"]
        A_Grid["VideoGrid"]
        A_Chat["ChatPanel"]
        A_Sig["useSignaling"]
        A_RTC["useWebRTC"]
        A_App --> A_Grid
        A_App --> A_Chat
        A_App --> A_Sig
        A_App --> A_RTC
    end

    subgraph "Go Backend (SFU)"
        Echo["Echo Server<br/>:8080"]
        Handler["SignalingHandler"]
        Room["RoomManager<br/>(룸당 최대 6명)"]
        SFU["sfu.Peer<br/>(PeerConnection<br/>+ RTP 포워딩)"]
        Echo --> Handler
        Handler --> Room
        Handler --> SFU
    end

    subgraph Browser B
        B_App["App.tsx"]
        B_Grid["VideoGrid"]
        B_Chat["ChatPanel"]
        B_Sig["useSignaling"]
        B_RTC["useWebRTC"]
        B_App --> B_Grid
        B_App --> B_Chat
        B_App --> B_Sig
        B_App --> B_RTC
    end

    subgraph Browser C
        C_App["App.tsx"]
        C_RTC["useWebRTC"]
        C_App --> C_RTC
    end

    A_Sig <-->|"WebSocket<br/>(시그널링/채팅)"| Echo
    B_Sig <-->|"WebSocket<br/>(시그널링/채팅)"| Echo
    C_App <-->|"WebSocket"| Echo

    A_RTC <-->|"WebRTC<br/>(미디어)"| SFU
    B_RTC <-->|"WebRTC<br/>(미디어)"| SFU
    C_RTC <-->|"WebRTC<br/>(미디어)"| SFU
```

## 동작 방식

### 1단계: Room 입장 및 시그널링 연결

```mermaid
sequenceDiagram
    participant A as Browser A
    participant S as SFU Server
    participant B as Browser B

    A->>S: WebSocket 연결 (/ws?roomId=room-123)
    Note over S: Peer A 생성 (PeerConnection + UUID)
    S-->>B: join 알림 (senderId: A)

    A->>S: Offer (SDP)
    Note over S: SetRemoteDescription → CreateAnswer
    S->>A: Answer (SDP)

    A->>S: ICE Candidate
    S->>A: ICE Candidate
    Note over A,S: ICE 연결 완료 → 미디어 전송 시작
```

### 2단계: 새 참가자 입장 및 트랙 포워딩

```mermaid
sequenceDiagram
    participant A as Browser A
    participant S as SFU Server
    participant B as Browser B

    Note over S: A의 OnTrack 발생 → TrackLocalStaticRTP 생성
    Note over S: RTP 포워딩 goroutine 시작

    B->>S: WebSocket 연결 + Offer
    S->>B: Answer
    Note over B,S: ICE 연결 완료

    Note over S: A의 기존 트랙을 B의 PC에 AddTrack
    S->>B: Renegotiation Offer (새 트랙 포함)
    B->>S: Answer
    Note over B: ontrack → Remote Stream 표시

    Note over S: B의 OnTrack 발생 → A의 PC에 AddTrack
    S->>A: Renegotiation Offer
    A->>S: Answer
    Note over A: ontrack → Remote Stream 표시
```

### 3단계: 참가자 퇴장 및 정리

```mermaid
sequenceDiagram
    participant A as Browser A
    participant S as SFU Server
    participant B as Browser B

    Note over A: WebSocket 연결 종료
    Note over S: A의 트랙을 B의 PC에서 RemoveTrack
    S->>B: Renegotiation Offer (트랙 제거)
    B->>S: Answer
    Note over B: 비디오 그리드에서 A 제거

    S->>B: leave 알림 (senderId: A)
    Note over S: Room에서 A 제거, PeerConnection Close
```

## 메시지 프로토콜

### Client → Server

| type | 용도 | payload |
|------|------|---------|
| `offer` | SDP Offer 전송 | `RTCSessionDescription` |
| `answer` | SDP Answer 응답 (renegotiation) | `RTCSessionDescription` |
| `ice` | ICE Candidate 전송 | `RTCIceCandidate` |
| `chat` | 채팅 메시지 | `senderName`, `message` |

### Server → Client

| type | 용도 | payload |
|------|------|---------|
| `offer` | Renegotiation Offer | `RTCSessionDescription` |
| `answer` | SDP Answer 응답 | `RTCSessionDescription` |
| `ice` | ICE Candidate 전송 | `RTCIceCandidate` |
| `join` | 새 참가자 입장 알림 | `senderId` |
| `leave` | 참가자 퇴장 알림 | `senderId` |
| `chat` | 채팅 메시지 브로드캐스트 | `senderId`, `senderName`, `message` |

## 실행 방법

### 사전 요구사항

- Go 1.24+
- Node.js
- WebRTC 지원 브라우저 (Chrome, Firefox, Edge, Safari)

### 실행

```bash
# 터미널 1 - Backend (localhost:8080)
make run-be

# 터미널 2 - Frontend (localhost:5173)
make install-fe  # 최초 1회
make run-fe
```

### 사용법

1. 브라우저에서 `http://localhost:5173` 접속
2. 이름(선택)과 Room ID를 입력하고 **입장** 클릭
3. 다른 브라우저 탭에서 같은 Room ID로 입장 → 자동으로 영상 연결
4. 우측 채팅 패널에서 텍스트 메시지 전송
5. 최대 6명까지 동시 입장 가능
