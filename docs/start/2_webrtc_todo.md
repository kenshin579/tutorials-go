# TODO: Pion SFU 기반 다자간 화상 통화

## Phase 1: 프로젝트 초기 설정

- [x] 기존 `webrtc/` 하위 파일(Makefile, README.md, backend/, frontend/)을 `webrtc/simple-p2p/`로 이동
- [x] 이동 후 `webrtc/simple-p2p/` 내 import 경로 및 go.mod 수정
- [x] `webrtc/multi-users-sfu/backend/` 디렉터리 생성 및 `go mod init`
- [x] `go.mod`에 의존성 추가 (echo, gorilla/websocket, pion/webrtc/v4)
- [x] `webrtc/multi-users-sfu/frontend/` Vite + React 18 + TypeScript 프로젝트 생성
- [x] Frontend 불필요한 boilerplate 정리

## Phase 2: Backend - Room 관리 및 Signaling

- [ ] `room/manager.go` 구현 (Room 구조체, Join/Leave/Broadcast)
- [ ] `handler/signaling.go` 구현 (WebSocket 핸들러, 메시지 파싱 및 라우팅)
- [ ] `main.go` 구현 (Echo 서버 부트스트랩, 라우트 등록)
- [ ] WebSocket 연결 및 메시지 송수신 동작 확인

## Phase 3: Backend - SFU 핵심 로직 (Pion WebRTC)

- [ ] `sfu/peer.go` 구현 (Peer 구조체, PeerConnection 생성, ICE 콜백)
- [ ] `sfu/router.go` 구현 (OnTrack → RTP 포워딩 루프)
- [ ] Signaling 핸들러에서 offer/answer/ice 메시지 → SFU PeerConnection 연동
- [ ] 새 Peer 입장 시 기존 트랙 구독 (AddTrack + Renegotiation)
- [ ] Peer 퇴장 시 트랙 제거 및 리소스 정리 (RemoveTrack + PeerConnection Close)

## Phase 4: Frontend - WebSocket Signaling Hook

- [ ] `hooks/useSignaling.ts` 구현 (WebSocket 연결, 메시지 타입 확장)
- [ ] 채팅 메시지 WebSocket 송수신 로직 구현

## Phase 5: Frontend - WebRTC Hook (SFU 연동)

- [ ] `hooks/useWebRTC.ts` 구현 (PeerConnection 생성, 로컬 미디어 획득)
- [ ] Offer 생성 → SFU 전송 → Answer 수신 흐름 구현
- [ ] ICE Candidate 교환 구현
- [ ] `ontrack` 이벤트에서 Remote Stream Map 관리
- [ ] SFU Renegotiation Offer 수신 → 자동 Answer 응답 구현

## Phase 6: Frontend - UI 컴포넌트

- [ ] `components/VideoGrid.tsx` 구현 (다자간 비디오 그리드, 동적 레이아웃)
- [ ] `components/ChatPanel.tsx` 구현 (WebSocket 채팅, 발신자 이름 표시)
- [ ] `App.tsx` 구현 (Room 입장 화면 + 통화 화면 조합)

## Phase 7: 통합 테스트 (MCP Playwright)

- [ ] Backend 서버 시작 및 Frontend 개발 서버 시작
- [ ] MCP Playwright로 Room 입장 화면 렌더링 확인
- [ ] MCP Playwright로 Room 입장 후 로컬 비디오 표시 확인
- [ ] 브라우저 탭 3개에서 동일 Room 입장 → 다자간 영상 송출 확인
- [ ] 채팅 메시지 전송 → Room 내 전체 수신 확인
- [ ] 참가자 퇴장 → 비디오 제거 및 UI 업데이트 확인
- [ ] 콘솔 로그에서 ICE/RTP 포워딩 흐름 확인
