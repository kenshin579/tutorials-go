# TODO: Echo + WebRTC + React 1:1 화상 통화

## Phase 1: 프로젝트 초기화

- [x] `webrtc/backend/` 디렉토리 생성 및 `go mod init` 실행
- [x] Echo v4, Gorilla WebSocket 의존성 추가
- [x] `webrtc/frontend/` Vite + React 18 + TypeScript 프로젝트 생성
- [x] Backend CORS 설정 (localhost:5173 허용)

## Phase 2: Backend - Signaling 서버

- [x] Room Manager 구현 (roomId → connections 매핑, 최대 2명 제한)
- [x] WebSocket 핸들러 구현 (`GET /ws?roomId={roomId}`)
- [x] Signaling 메시지 relay 구현 (offer/answer/ice)
- [x] 연결/해제 시 Room 정리 로직 구현
- [x] Signaling 이벤트 로그 출력 확인

## Phase 3: Frontend - WebSocket Signaling

- [x] `useSignaling` Hook 구현 (WebSocket 연결, 메시지 송수신)
- [x] roomId 입력 UI 구현 (App.tsx)
- [x] WebSocket 연결 상태 표시
- [x] MCP Playwright로 roomId 입력 → WebSocket 연결 테스트

## Phase 4: Frontend - WebRTC 화상 통화

- [x] `useWebRTC` Hook 구현 (RTCPeerConnection을 useRef로 관리)
- [x] `getUserMedia()` → Local Video 표시
- [x] Offer/Answer 생성 및 교환 로직 구현
- [x] ICE Candidate 교환 로직 구현
- [x] Remote Video 수신 및 표시
- [x] VideoPanel 컴포넌트 구현 (local + remote video)

## Phase 5: Frontend - DataChannel 채팅

- [x] DataChannel 생성 (Caller) / 수신 (Callee) 구현
- [x] ChatPanel 컴포넌트 구현 (메시지 목록 + 입력창)
- [x] 실시간 메시지 송수신 연동

## Phase 6: 통합 테스트

- [x] Backend 서버 실행 확인 (`go run main.go`)
- [x] Frontend 개발 서버 실행 확인 (`npm run dev`)
- [x] MCP Playwright로 roomId 입력 → 통화 화면 전환 테스트
- [x] WebSocket Signaling 연결 확인 (콘솔 로그: `[Signaling] connected`)
- [ ] MCP Playwright로 두 브라우저 탭에서 영상/음성 P2P 연결 확인
- [ ] MCP Playwright로 DataChannel 채팅 송수신 테스트
