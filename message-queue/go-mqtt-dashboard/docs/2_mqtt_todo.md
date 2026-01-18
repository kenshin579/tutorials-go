# MQTT v5 실시간 디바이스 상태 대시보드 TODO

## Phase 1: 인프라 설정

- [x] 프로젝트 디렉토리 생성 (`message-queue/go-mqtt-dashboard/`)
- [x] Docker Compose 파일 작성
- [x] Mosquitto 설정 파일 작성 (MQTT + WebSocket 리스너)
- [x] Makefile 작성 (setup, run-broker, run-be, run-fe, stop-broker, clean)
- [ ] Mosquitto 컨테이너 실행 및 연결 테스트 (`make run-broker`)

## Phase 2: Backend 구현

### 2.1 프로젝트 초기화
- [x] Go 모듈 초기화 (`go mod init`)
- [x] autopaho 의존성 추가

### 2.2 MQTT Client 구현
- [x] `backend/internal/mqtt/client.go` 생성
- [x] MQTT 연결 로직 구현
- [x] Subscribe 로직 구현 (device/1/command)
- [x] Publish 로직 구현 (device/1/state)
- [x] 자동 재연결 설정

### 2.3 Device Simulator 구현
- [x] `backend/internal/device/simulator.go` 생성
- [x] State 구조체 정의
- [x] GetState 메서드 구현 (랜덤 온도 생성)
- [x] HandleCommand 메서드 구현 (start/stop)

### 2.4 Main 구현
- [x] `backend/cmd/main.go` 생성
- [x] MQTT Client 초기화
- [x] 2초 주기 상태 Publish 로직
- [x] Command 수신 및 처리 연동

### 2.5 Backend 테스트
- [ ] MQTT 연결 테스트
- [ ] State Publish 확인 (mosquitto_sub로 검증)
- [ ] Command 처리 확인 (mosquitto_pub로 검증)

## Phase 3: Frontend 구현

### 3.1 프로젝트 초기화
- [x] React 프로젝트 생성 (Create React App 또는 Vite)
- [x] TypeScript 설정
- [x] mqtt.js 의존성 추가

### 3.2 MQTT Hook 구현
- [x] `frontend/src/hooks/useMqtt.ts` 생성
- [x] MQTT 연결 로직 (WebSocket)
- [x] Subscribe 로직 (device/1/state)
- [x] Publish 로직 (device/1/command)
- [x] 연결 상태 관리
- [x] 자동 재연결 설정
- [x] 메시지 로그 히스토리 관리 (최대 50개)
- [x] 로그 초기화 함수 구현

### 3.3 UI 컴포넌트 구현
- [x] `frontend/src/components/DeviceStatus.tsx` 생성
- [x] `frontend/src/components/DeviceStatus.module.css` 생성
- [x] 연결 상태 표시 (Connected/Disconnected)
- [x] 디바이스 Status 표시 (아이콘 포함)
- [x] Temperature 표시
- [x] Start 버튼 구현 (녹색 스타일링)
- [x] Stop 버튼 구현 (빨간색 스타일링)
- [x] 버튼 비활성화 스타일 적용
- [x] 메시지 로그 영역 구현 (수신/송신 구분)
- [x] 로그 Clear 버튼 구현

### 3.4 App 통합
- [x] `frontend/src/App.tsx` 수정
- [x] DeviceStatus 컴포넌트 연동

## Phase 4: 통합 테스트

> **MCP Playwright**를 사용하여 E2E 테스트 수행

### 4.1 테스트 환경 준비
- [x] 전체 시스템 실행 (`make run-broker`, `make run-be`, `make run-fe`)
- [x] Frontend 접속 확인 (http://localhost:5173)

### 4.2 MCP Playwright 테스트 시나리오

```
# 테스트 시나리오 예시 (Claude Code에서 실행)

1. playwright_navigate: http://localhost:3000 접속
2. playwright_screenshot: 초기 화면 캡처
3. playwright_get_visible_text: Connection 상태 확인 (Connected)
4. playwright_click: Start 버튼 클릭
5. playwright_screenshot: status 변경 확인 (running)
6. playwright_click: Stop 버튼 클릭
7. playwright_screenshot: status 변경 확인 (idle)
8. playwright_get_visible_text: 로그 영역에 메시지 표시 확인
```

### 4.3 테스트 체크리스트
- [x] 연결 상태 표시 확인 (🟢 Connected)
- [x] 디바이스 상태 실시간 업데이트 확인
- [x] Start 버튼 클릭 → status: running 변경
- [x] Stop 버튼 클릭 → status: idle 변경
- [x] 메시지 로그에 수신/송신 기록 표시
- [x] Clear 버튼으로 로그 초기화
- [x] Broker 재시작 후 자동 재연결 확인

## Phase 5: 문서화

- [x] README.md 작성 (실행 방법, 아키텍처 설명)
- [x] 코드 주석 추가 (학습 목적)
