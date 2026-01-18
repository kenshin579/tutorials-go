# MQTT v5 실시간 디바이스 상태 대시보드 PRD (MVP)

## 1. 개요

### 1.1 프로젝트 목적

본 프로젝트는 **MQTT v5 학습을 위한 최소 기능(MVP) 샘플 프로젝트**이다.
프론트엔드(FE)와 백엔드(BE)가 **동등한 MQTT Client**로 동작하며, Pub/Sub 기반 메시징, 재연결, 상태 동기화 개념을 실습하는 것을 목표로 한다.

본 프로젝트는 실제 서비스 목적이 아닌 **학습·데모·레퍼런스 코드** 제공을 위한 프로젝트이다.

### 1.2 핵심 학습 목표

| 목표 | 설명 |
|------|------|
| Pub/Sub 양방향 흐름 | MQTT v5 Publish / Subscribe 패턴 이해 |
| WebSocket 연동 | MQTT over WebSocket (FE) 사용 경험 |
| 메시지 분리 설계 | 상태(State) / 명령(Command) 분리 패턴 |
| Retained Message | Retained Message를 통한 상태 복구 |
| 재연결 처리 | Reconnect 시나리오 이해 및 구현 |

### 1.3 기술 스택

| 구분 | 기술 |
|------|------|
| Frontend | React, mqtt.js (MQTT v5 지원), MQTT over WebSocket |
| Backend | Go 1.25, Eclipse Paho MQTT v5 (autopaho) |
| Broker | Eclipse Mosquitto v5 (WebSocket Listener 활성화) |

## 2. 시스템 아키텍처

```
[ React FE ]  <--MQTT(WebSocket)-->  [ Mosquitto Broker ]  <--MQTT-->  [ Go BE ]
```

### 2.1 구성 요소

| 구성 요소 | 역할 |
|----------|------|
| React FE | MQTT Client (WebSocket), 상태 표시 및 명령 발행 |
| Mosquitto Broker | MQTT v5 메시지 브로커 |
| Go BE | MQTT Client, 가상 디바이스 시뮬레이션 |

## 3. MQTT Topic 설계

### 3.1 Topic 목록

| Topic | 방향 | 설명 |
|-------|------|------|
| `device/1/state` | BE → FE | 디바이스 현재 상태 |
| `device/1/command` | FE → BE | 디바이스 제어 명령 |

### 3.2 Topic 설계 원칙

- 단일 디바이스(`device/1`)만 사용 (MVP 단순화)
- Command / State 명확히 분리
- Wildcard 사용 금지

## 4. 메시지 스펙

### 4.1 Device State 메시지

| 항목 | 값 |
|------|-----|
| Topic | `device/1/state` |
| Publish 주기 | 2초 |
| QoS | 0 |
| Retained | true |

```json
{
  "deviceId": "1",
  "status": "idle | running",
  "temperature": 36.5,
  "timestamp": 1700000000
}
```

### 4.2 Device Command 메시지

| 항목 | 값 |
|------|-----|
| Topic | `device/1/command` |
| Publish 시점 | FE 버튼 클릭 |
| QoS | 1 |
| Retained | false |

```json
{
  "action": "start | stop"
}
```

## 5. Backend (Go) 요구사항

### 5.1 역할

- 가상 디바이스 시뮬레이션
- 상태(State) 주기적 Publish
- 명령(Command) Subscribe 및 처리

### 5.2 기능 요구사항

| 요구사항 ID | 설명 |
|------------|------|
| BE-001 | MQTT Broker에 연결한다 |
| BE-002 | 연결 성공 시 `device/1/command` Topic을 Subscribe 한다 |
| BE-003 | 내부 상태를 유지한다 (status: idle/running, temperature: 랜덤 값) |
| BE-004 | 2초마다 상태를 `device/1/state`로 Publish 한다 |
| BE-005 | Command 수신 시: `start` → status = running, `stop` → status = idle |

### 5.3 재연결 요구사항

| 요구사항 ID | 설명 |
|------------|------|
| BE-RC-001 | 네트워크 단절 시 자동 재연결 |
| BE-RC-002 | 재연결 후 Subscribe 자동 복구 |
| BE-RC-003 | 재연결 후 상태 Publish 재개 |

## 6. Frontend (React) 요구사항

### 6.1 역할

- 디바이스 상태 실시간 표시
- 디바이스 제어 명령 발행

### 6.2 기능 요구사항

| 요구사항 ID | 설명 |
|------------|------|
| FE-001 | MQTT Broker(WebSocket) 연결 |
| FE-002 | `device/1/state` Topic Subscribe |
| FE-003 | 상태 메시지 수신 시 UI 업데이트 |
| FE-004 | Start / Stop 버튼 제공 |
| FE-005 | 버튼 클릭 시 Command Publish |

### 6.3 UI 요구사항 (MVP)

| 요구사항 ID | 설명 |
|------------|------|
| FE-UI-001 | 연결 상태 표시 (Connected / Disconnected) |
| FE-UI-002 | 현재 Status 표시 (idle / running) |
| FE-UI-003 | 현재 Temperature 표시 |
| FE-UI-004 | Start 버튼 |
| FE-UI-005 | Stop 버튼 |

## 7. 프로젝트 구조

```
message-queue/go-mqtt-dashboard/
├── backend/
│   ├── cmd/
│   │   └── main.go
│   ├── internal/
│   │   ├── device/
│   │   │   └── simulator.go
│   │   └── mqtt/
│   │       └── client.go
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── App.tsx
│   │   ├── components/
│   │   │   ├── DeviceStatus.tsx
│   │   │   └── DeviceStatus.module.css
│   │   └── hooks/
│   │       └── useMqtt.ts
│   └── package.json
├── mosquitto/
│   └── config/
│       └── mosquitto.conf
├── docker-compose.yml
├── Makefile
└── README.md
```

## 8. 성공 기준 (Definition of Done)

| 기준 ID | 설명 |
|--------|------|
| DOD-001 | FE / BE 모두 MQTT 연결 성공 |
| DOD-002 | 상태 메시지가 실시간으로 화면에 표시됨 |
| DOD-003 | Command 버튼으로 상태 변경 가능 |
| DOD-004 | Broker 재시작 후 자동 복구 확인 |

## 9. 향후 확장 아이디어 (Optional)

| 아이디어 | 설명 |
|---------|------|
| 멀티 디바이스 | 여러 디바이스 시뮬레이션 |
| Shared Subscription | MQTT v5 Shared Subscription 활용 |
| Message Expiry | Message Expiry Interval 실험 |
| 인증 / ACL | Broker 인증 및 접근 제어 추가 |

## 10. 비고

본 프로젝트는 **학습용 레퍼런스 구현**이며, 모든 설계는 단순성과 가독성을 우선한다.
