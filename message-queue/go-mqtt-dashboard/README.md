# MQTT v5 실시간 디바이스 상태 대시보드

MQTT v5 프로토콜을 활용한 실시간 디바이스 모니터링 시스템입니다.

## 아키텍처

```
┌─────────────────┐     WebSocket(9001)     ┌─────────────────┐
│    Frontend     │◄──────────────────────►│                 │
│  (React + TS)   │                         │    Mosquitto    │
└─────────────────┘                         │   MQTT Broker   │
                                            │                 │
┌─────────────────┐     TCP(1883)           │                 │
│    Backend      │◄──────────────────────►│                 │
│  (Go + autopaho)│                         └─────────────────┘
└─────────────────┘
```

### 토픽 설계

| 토픽 | Publisher | Subscriber | QoS | 설명 |
|------|-----------|------------|-----|------|
| `device/1/state` | Backend | Frontend | 0 | 디바이스 상태 (2초 주기) |
| `device/1/command` | Frontend | Backend | 1 | 명령 (start/stop) |

### 메시지 형식

**State (Backend → Frontend)**
```json
{
  "deviceId": "device-1",
  "status": "running",
  "temperature": 23.5,
  "timestamp": 1705580400000
}
```

**Command (Frontend → Backend)**
```json
{
  "action": "start"
}
```

## 기술 스택

- **Frontend**: React + TypeScript + mqtt.js + Vite
- **Backend**: Go + eclipse/paho.golang (autopaho)
- **Broker**: Eclipse Mosquitto v2

## 사전 요구사항

- Docker & Docker Compose
- Go 1.21+
- Node.js 18+

## 실행 방법

### 1. Mosquitto 브로커 실행

```bash
make run-broker
```

브로커가 실행되면:
- TCP: localhost:1883
- WebSocket: localhost:9001

### 2. Backend 실행

```bash
make run-be
```

2초마다 디바이스 상태를 publish합니다.

### 3. Frontend 실행

```bash
make run-fe
```

http://localhost:3000 에서 대시보드에 접근합니다.

## 주요 기능

- 실시간 디바이스 상태 모니터링
- Start/Stop 명령 전송
- 연결 상태 표시 (Connected/Disconnected)
- 메시지 로그 히스토리 (수신/송신 구분)
- 자동 재연결

## 프로젝트 구조

```
go-mqtt-dashboard/
├── backend/
│   ├── cmd/
│   │   └── main.go          # 진입점
│   └── internal/
│       ├── device/
│       │   └── simulator.go # 디바이스 시뮬레이터
│       └── mqtt/
│           └── client.go    # MQTT 클라이언트 래퍼
├── frontend/
│   └── src/
│       ├── hooks/
│       │   └── useMqtt.ts   # MQTT 훅
│       └── components/
│           ├── DeviceStatus.tsx
│           └── DeviceStatus.module.css
├── mosquitto/
│   └── config/
│       └── mosquitto.conf   # 브로커 설정
├── docker-compose.yml
├── Makefile
└── README.md
```

## Makefile 명령어

```bash
make help        # 도움말
make setup       # 의존성 설치
make run-broker  # Mosquitto 시작
make run-be      # Backend 실행
make run-fe      # Frontend 실행
make stop-broker # Mosquitto 중지
make clean       # 정리
```

## 학습 포인트

### MQTT v5 주요 개념

1. **Pub/Sub 패턴**: Publisher는 토픽에 메시지를 발행하고, Subscriber는 토픽을 구독
2. **QoS (Quality of Service)**:
   - QoS 0: At most once (빠르지만 손실 가능)
   - QoS 1: At least once (재전송으로 확실한 전달)
3. **WebSocket 지원**: 브라우저에서 MQTT 사용 가능
4. **자동 재연결**: autopaho, mqtt.js 모두 내장 지원

### Frontend (mqtt.js)

```typescript
const client = mqtt.connect('ws://localhost:9001', {
  protocolVersion: 5,
  reconnectPeriod: 1000,
});

client.subscribe('device/1/state');
client.publish('device/1/command', JSON.stringify({ action: 'start' }));
```

### Backend (autopaho)

```go
cfg := autopaho.ClientConfig{
    ServerUrls: []*url.URL{brokerURL},
    KeepAlive:  30,
    OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
        cm.Subscribe(ctx, &paho.Subscribe{
            Subscriptions: []paho.SubscribeOptions{
                {Topic: "device/1/command", QoS: 1},
            },
        })
    },
}
```

## 테스트

### mosquitto 클라이언트로 수동 테스트

```bash
# 상태 구독
mosquitto_sub -h localhost -p 1883 -t "device/1/state" -v

# 명령 발행
mosquitto_pub -h localhost -p 1883 -t "device/1/command" -m '{"action":"start"}'
```

## 라이선스

MIT
