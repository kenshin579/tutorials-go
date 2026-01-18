# MQTT v5 실시간 디바이스 상태 대시보드 구현 문서

## 1. 프로젝트 구조

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
│   ├── package.json
│   └── tsconfig.json
├── mosquitto/
│   └── config/
│       └── mosquitto.conf
├── docker-compose.yml
├── Makefile
└── README.md
```

## 2. 인프라 구성

### 2.1 Docker Compose

Mosquitto Broker를 Docker로 실행한다.

```yaml
# docker-compose.yml
version: '3.8'
services:
  mosquitto:
    image: eclipse-mosquitto:2
    container_name: mosquitto
    ports:
      - "1883:1883"   # MQTT
      - "9001:9001"   # WebSocket
    volumes:
      - ./mosquitto/config:/mosquitto/config
      - ./mosquitto/data:/mosquitto/data
      - ./mosquitto/log:/mosquitto/log
```

### 2.2 Mosquitto 설정

```conf
# mosquitto/config/mosquitto.conf
listener 1883
listener 9001
protocol websockets

allow_anonymous true
```

### 2.3 연결 방식: TCP vs WebSocket

브라우저는 보안상 직접 TCP 소켓을 열 수 없다. 따라서 Frontend는 **MQTT over WebSocket**을 사용한다.

| 환경 | 연결 방식 | URL | 포트 |
|------|----------|-----|------|
| Backend (Go) | TCP 소켓 | `mqtt://localhost:1883` | 1883 |
| Frontend (브라우저) | WebSocket | `ws://localhost:9001` | 9001 |

```
[React 브라우저]
      │
      ▼ WebSocket (ws://localhost:9001)
      │
[Mosquitto Broker]  ← WebSocket을 MQTT로 변환
      │
      ▼ TCP (mqtt://localhost:1883)
      │
[Go Backend]
```

- `ws://` → WebSocket 연결
- `mqtt.js` → WebSocket 위에 MQTT 프로토콜을 캡슐화
- Mosquitto → WebSocket 메시지를 MQTT로 해석

## 3. Backend 구현

### 3.1 MQTT Client (autopaho)

```go
// internal/mqtt/client.go
package mqtt

import (
    "context"
    "net/url"
    "github.com/eclipse/paho.golang/autopaho"
    "github.com/eclipse/paho.golang/paho"
)

type Client struct {
    conn *autopaho.ConnectionManager
}

func NewClient(ctx context.Context, brokerURL string, onMessage func(topic string, payload []byte)) (*Client, error) {
    u, _ := url.Parse(brokerURL)

    cfg := autopaho.ClientConfig{
        BrokerUrls: []*url.URL{u},
        KeepAlive:  30,
        OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
            // Subscribe on connect
            cm.Subscribe(ctx, &paho.Subscribe{
                Subscriptions: []paho.SubscribeOptions{
                    {Topic: "device/1/command", QoS: 1},
                },
            })
        },
        ClientConfig: paho.ClientConfig{
            ClientID: "go-backend",
            OnPublishReceived: []func(paho.PublishReceived) (bool, error){
                func(pr paho.PublishReceived) (bool, error) {
                    onMessage(pr.Packet.Topic, pr.Packet.Payload)
                    return true, nil
                },
            },
        },
    }

    conn, err := autopaho.NewConnection(ctx, cfg)
    if err != nil {
        return nil, err
    }

    return &Client{conn: conn}, nil
}

func (c *Client) Publish(ctx context.Context, topic string, payload []byte, qos byte, retain bool) error {
    _, err := c.conn.Publish(ctx, &paho.Publish{
        Topic:   topic,
        QoS:     qos,
        Retain:  retain,
        Payload: payload,
    })
    return err
}
```

### 3.2 Device Simulator

```go
// internal/device/simulator.go
package device

import (
    "context"
    "encoding/json"
    "math/rand"
    "sync"
    "time"
)

type State struct {
    DeviceID    string  `json:"deviceId"`
    Status      string  `json:"status"`
    Temperature float64 `json:"temperature"`
    Timestamp   int64   `json:"timestamp"`
}

type Simulator struct {
    mu     sync.RWMutex
    status string
}

func NewSimulator() *Simulator {
    return &Simulator{status: "idle"}
}

func (s *Simulator) GetState() State {
    s.mu.RLock()
    defer s.mu.RUnlock()

    return State{
        DeviceID:    "1",
        Status:      s.status,
        Temperature: 35.0 + rand.Float64()*5.0,
        Timestamp:   time.Now().Unix(),
    }
}

func (s *Simulator) HandleCommand(action string) {
    s.mu.Lock()
    defer s.mu.Unlock()

    switch action {
    case "start":
        s.status = "running"
    case "stop":
        s.status = "idle"
    }
}

func (s *Simulator) ToJSON() ([]byte, error) {
    return json.Marshal(s.GetState())
}
```

### 3.3 Main

```go
// cmd/main.go
package main

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "go-mqtt-dashboard/backend/internal/device"
    "go-mqtt-dashboard/backend/internal/mqtt"
)

func main() {
    ctx := context.Background()
    sim := device.NewSimulator()

    client, err := mqtt.NewClient(ctx, "mqtt://localhost:1883", func(topic string, payload []byte) {
        var cmd struct {
            Action string `json:"action"`
        }
        if err := json.Unmarshal(payload, &cmd); err == nil {
            sim.HandleCommand(cmd.Action)
            log.Printf("Command received: %s", cmd.Action)
        }
    })
    if err != nil {
        log.Fatal(err)
    }

    // Publish state every 2 seconds
    ticker := time.NewTicker(2 * time.Second)
    for range ticker.C {
        state := sim.GetState()
        payload, _ := json.Marshal(state)
        client.Publish(ctx, "device/1/state", payload, 0, true)
        log.Printf("State published: status=%s, temperature=%.1f", state.Status, state.Temperature)
    }
}
```

## 4. Frontend 구현

### 4.1 UI 목업

```
┌─────────────────────────────────────────┐
│                                         │
│        🖥️ Device Dashboard              │
│                                         │
│  ┌───────────────────────────────────┐  │
│  │ Connection: 🟢 Connected          │  │
│  └───────────────────────────────────┘  │
│                                         │
│  ┌───────────────────────────────────┐  │
│  │ Status      │  🔵 running         │  │
│  ├─────────────┼─────────────────────┤  │
│  │ Temperature │  37.2°C             │  │
│  └─────────────┴─────────────────────┘  │
│                                         │
│  ┌─────────────┐    ┌─────────────┐     │
│  │   ▶ Start   │    │   ⏹ Stop    │     │
│  └─────────────┘    └─────────────┘     │
│                                         │
│  ┌───────────────────────────────────┐  │
│  │ 📋 Message Log           [Clear]  │  │
│  ├───────────────────────────────────┤  │
│  │ 10:00:05 ← STATE idle 36.8°C      │  │
│  │ 10:00:03 → CMD   start            │  │
│  │ 10:00:03 ← STATE running 37.2°C   │  │
│  │ 10:00:01 ← STATE running 36.5°C   │  │
│  │ ...                               │  │
│  └───────────────────────────────────┘  │
│                                         │
└─────────────────────────────────────────┘
```

**상태별 화면 변화:**

| 상태 | Connection | Status 표시 | 버튼 |
|------|------------|-------------|------|
| 연결 전 | 🔴 Disconnected | (표시 안 됨) | 비활성화 (회색) |
| idle 상태 | 🟢 Connected | ⚪ idle | 활성화 |
| running 상태 | 🟢 Connected | 🔵 running | 활성화 |

**로그 표시 형식:**

| 방향 | 표시 | 예시 |
|------|------|------|
| 수신 (BE → FE) | `←` | `10:00:01 ← STATE running 37.2°C` |
| 송신 (FE → BE) | `→` | `10:00:03 → CMD start` |

### 4.2 MQTT Hook

```typescript
// src/hooks/useMqtt.ts
import { useState, useEffect, useCallback } from 'react';
import mqtt, { MqttClient } from 'mqtt';

interface DeviceState {
  deviceId: string;
  status: 'idle' | 'running';
  temperature: number;
  timestamp: number;
}

// 로그 엔트리 타입
interface LogEntry {
  id: number;
  time: string;
  direction: 'in' | 'out';  // in: 수신, out: 송신
  type: 'STATE' | 'CMD';
  message: string;
}

const MAX_LOG_SIZE = 50;  // 최대 로그 개수

export function useMqtt(brokerUrl: string) {
  const [client, setClient] = useState<MqttClient | null>(null);
  const [connected, setConnected] = useState(false);
  const [deviceState, setDeviceState] = useState<DeviceState | null>(null);
  const [logs, setLogs] = useState<LogEntry[]>([]);

  // 로그 추가 함수
  const addLog = useCallback((direction: 'in' | 'out', type: 'STATE' | 'CMD', message: string) => {
    const now = new Date();
    const time = now.toLocaleTimeString('ko-KR', { hour12: false });

    setLogs(prev => {
      const newLog: LogEntry = {
        id: Date.now(),
        time,
        direction,
        type,
        message,
      };
      const updated = [newLog, ...prev];
      return updated.slice(0, MAX_LOG_SIZE);  // 최대 개수 제한
    });
  }, []);

  // 로그 초기화 함수
  const clearLogs = useCallback(() => {
    setLogs([]);
  }, []);

  useEffect(() => {
    const mqttClient = mqtt.connect(brokerUrl, {
      protocolVersion: 5,
      reconnectPeriod: 1000,
    });

    mqttClient.on('connect', () => {
      setConnected(true);
      mqttClient.subscribe('device/1/state');
    });

    mqttClient.on('close', () => setConnected(false));

    mqttClient.on('message', (topic, payload) => {
      if (topic === 'device/1/state') {
        const state: DeviceState = JSON.parse(payload.toString());
        setDeviceState(state);
        // 수신 로그 추가
        addLog('in', 'STATE', `${state.status} ${state.temperature.toFixed(1)}°C`);
      }
    });

    setClient(mqttClient);

    return () => {
      mqttClient.end();
    };
  }, [brokerUrl, addLog]);

  const sendCommand = useCallback((action: 'start' | 'stop') => {
    client?.publish('device/1/command', JSON.stringify({ action }), { qos: 1 });
    // 송신 로그 추가
    addLog('out', 'CMD', action);
  }, [client, addLog]);

  return { connected, deviceState, logs, sendCommand, clearLogs };
}
```

### 4.3 DeviceStatus 컴포넌트

```tsx
// src/components/DeviceStatus.tsx
import React from 'react';
import { useMqtt } from '../hooks/useMqtt';
import styles from './DeviceStatus.module.css';

export function DeviceStatus() {
  const { connected, deviceState, logs, sendCommand, clearLogs } = useMqtt('ws://localhost:9001');

  return (
    <div className={styles.container}>
      <h1 className={styles.title}>🖥️ Device Dashboard</h1>

      <div className={styles.connectionStatus}>
        <span>Connection: </span>
        <span className={connected ? styles.connected : styles.disconnected}>
          {connected ? '🟢 Connected' : '🔴 Disconnected'}
        </span>
      </div>

      {deviceState && (
        <table className={styles.stateTable}>
          <tbody>
            <tr>
              <td>Status</td>
              <td className={deviceState.status === 'running' ? styles.running : styles.idle}>
                {deviceState.status === 'running' ? '🔵' : '⚪'} {deviceState.status}
              </td>
            </tr>
            <tr>
              <td>Temperature</td>
              <td>{deviceState.temperature.toFixed(1)}°C</td>
            </tr>
          </tbody>
        </table>
      )}

      <div className={styles.buttonGroup}>
        <button
          className={styles.startButton}
          onClick={() => sendCommand('start')}
          disabled={!connected}
        >
          ▶ Start
        </button>
        <button
          className={styles.stopButton}
          onClick={() => sendCommand('stop')}
          disabled={!connected}
        >
          ⏹ Stop
        </button>
      </div>

      {/* 메시지 로그 영역 */}
      <div className={styles.logSection}>
        <div className={styles.logHeader}>
          <span>📋 Message Log</span>
          <button className={styles.clearButton} onClick={clearLogs}>
            Clear
          </button>
        </div>
        <div className={styles.logList}>
          {logs.length === 0 ? (
            <div className={styles.logEmpty}>No messages yet</div>
          ) : (
            logs.map(log => (
              <div
                key={log.id}
                className={`${styles.logEntry} ${log.direction === 'in' ? styles.logIn : styles.logOut}`}
              >
                <span className={styles.logTime}>{log.time}</span>
                <span className={styles.logDirection}>
                  {log.direction === 'in' ? '←' : '→'}
                </span>
                <span className={styles.logType}>{log.type}</span>
                <span className={styles.logMessage}>{log.message}</span>
              </div>
            ))
          )}
        </div>
      </div>
    </div>
  );
}
```

### 4.4 CSS 스타일링

```css
/* src/components/DeviceStatus.module.css */
.container {
  max-width: 400px;
  margin: 40px auto;
  padding: 24px;
  border: 1px solid #e0e0e0;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.title {
  text-align: center;
  margin-bottom: 24px;
  color: #333;
}

.connectionStatus {
  padding: 12px;
  background: #f5f5f5;
  border-radius: 8px;
  margin-bottom: 16px;
}

.connected {
  color: #2e7d32;
  font-weight: bold;
}

.disconnected {
  color: #c62828;
  font-weight: bold;
}

.stateTable {
  width: 100%;
  border-collapse: collapse;
  margin-bottom: 20px;
}

.stateTable td {
  padding: 12px;
  border: 1px solid #e0e0e0;
}

.stateTable td:first-child {
  background: #fafafa;
  font-weight: bold;
  width: 40%;
}

.running {
  color: #1565c0;
  font-weight: bold;
}

.idle {
  color: #757575;
}

.buttonGroup {
  display: flex;
  gap: 12px;
}

.startButton, .stopButton {
  flex: 1;
  padding: 12px 24px;
  border: none;
  border-radius: 8px;
  font-size: 16px;
  cursor: pointer;
  transition: background 0.2s;
}

.startButton {
  background: #4caf50;
  color: white;
}

.startButton:hover:not(:disabled) {
  background: #43a047;
}

.stopButton {
  background: #f44336;
  color: white;
}

.stopButton:hover:not(:disabled) {
  background: #e53935;
}

.startButton:disabled, .stopButton:disabled {
  background: #bdbdbd;
  cursor: not-allowed;
}

/* 로그 영역 스타일 */
.logSection {
  margin-top: 24px;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  overflow: hidden;
}

.logHeader {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px;
  background: #fafafa;
  border-bottom: 1px solid #e0e0e0;
  font-weight: bold;
}

.clearButton {
  padding: 4px 12px;
  border: 1px solid #bdbdbd;
  border-radius: 4px;
  background: white;
  cursor: pointer;
  font-size: 12px;
}

.clearButton:hover {
  background: #f5f5f5;
}

.logList {
  max-height: 200px;
  overflow-y: auto;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 12px;
}

.logEmpty {
  padding: 24px;
  text-align: center;
  color: #9e9e9e;
}

.logEntry {
  display: flex;
  gap: 8px;
  padding: 6px 12px;
  border-bottom: 1px solid #f0f0f0;
}

.logEntry:last-child {
  border-bottom: none;
}

.logIn {
  background: #e3f2fd;
}

.logOut {
  background: #fff3e0;
}

.logTime {
  color: #757575;
}

.logDirection {
  font-weight: bold;
}

.logIn .logDirection {
  color: #1565c0;
}

.logOut .logDirection {
  color: #e65100;
}

.logType {
  font-weight: bold;
  min-width: 50px;
}

.logMessage {
  color: #424242;
}
```

### 4.5 App

```tsx
// src/App.tsx
import React from 'react';
import { DeviceStatus } from './components/DeviceStatus';

function App() {
  return (
    <div className="App">
      <DeviceStatus />
    </div>
  );
}

export default App;
```

## 5. 의존성

### 5.1 Backend (go.mod)

```go
module go-mqtt-dashboard/backend

go 1.21

require (
    github.com/eclipse/paho.golang v0.21.0
)
```

### 5.2 Frontend (package.json)

```json
{
  "dependencies": {
    "mqtt": "^5.3.0",
    "react": "^18.2.0",
    "react-dom": "^18.2.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.0",
    "typescript": "^5.0.0"
  }
}
```

## 6. Makefile

```makefile
# Makefile
.PHONY: help setup run-broker run-be run-fe stop-broker clean

# 기본 타겟
help:
	@echo "MQTT Dashboard - Available Commands"
	@echo "===================================="
	@echo "make setup        - 프로젝트 초기 설정 (의존성 설치)"
	@echo "make run-broker   - Mosquitto 브로커 실행"
	@echo "make run-be       - Backend 실행"
	@echo "make run-fe       - Frontend 실행"
	@echo "make stop-broker  - Mosquitto 브로커 중지"
	@echo "make clean        - 정리 (컨테이너, node_modules 등)"

# 프로젝트 초기 설정
setup:
	@echo "📦 Installing Backend dependencies..."
	cd backend && go mod download
	@echo "📦 Installing Frontend dependencies..."
	cd frontend && npm install
	@echo "✅ Setup complete!"

# Mosquitto 브로커 실행
run-broker:
	@echo "🚀 Starting Mosquitto broker..."
	docker-compose up -d
	@echo "✅ Mosquitto running on ports 1883 (MQTT), 9001 (WebSocket)"

# Backend 실행
run-be:
	@echo "🚀 Starting Backend..."
	cd backend && go run cmd/main.go

# Frontend 실행
run-fe:
	@echo "🚀 Starting Frontend..."
	cd frontend && npm start

# Mosquitto 브로커 중지
stop-broker:
	@echo "🛑 Stopping Mosquitto broker..."
	docker-compose down

# 정리
clean:
	@echo "🧹 Cleaning up..."
	docker-compose down -v
	rm -rf frontend/node_modules
	rm -rf backend/vendor
	@echo "✅ Cleanup complete!"
```

## 7. 실행 방법

### 7.1 Make 명령어 사용 (권장)

```bash
# 프로젝트 디렉토리로 이동
cd message-queue/go-mqtt-dashboard

# 1. 초기 설정 (최초 1회)
make setup

# 2. Mosquitto 브로커 실행
make run-broker

# 3. Backend 실행 (새 터미널)
make run-be

# 4. Frontend 실행 (새 터미널)
make run-fe
```

### 7.2 개별 명령어 사용

```bash
# 프로젝트 디렉토리로 이동
cd message-queue/go-mqtt-dashboard

# 1. Mosquitto 실행
docker-compose up -d

# 2. Backend 실행
cd backend
go run cmd/main.go

# 3. Frontend 실행
cd frontend
npm install
npm start
```

## 8. 테스트 (MCP Playwright)

개발 완료 후 **MCP Playwright**를 사용하여 E2E 테스트를 수행한다.

### 8.1 테스트 시나리오

```bash
# Claude Code에서 MCP Playwright 도구를 사용하여 테스트

# 1. 브라우저에서 Frontend 접속
mcp__playwright__playwright_navigate url="http://localhost:3000"

# 2. 초기 화면 스크린샷
mcp__playwright__playwright_screenshot name="initial"

# 3. 페이지 텍스트 확인 (Connection 상태)
mcp__playwright__playwright_get_visible_text

# 4. Start 버튼 클릭
mcp__playwright__playwright_click selector="button:has-text('Start')"

# 5. 상태 변경 확인 스크린샷
mcp__playwright__playwright_screenshot name="after-start"

# 6. Stop 버튼 클릭
mcp__playwright__playwright_click selector="button:has-text('Stop')"

# 7. 최종 상태 확인
mcp__playwright__playwright_screenshot name="after-stop"

# 8. 브라우저 종료
mcp__playwright__playwright_close
```

### 8.2 테스트 체크포인트

| 항목 | 확인 내용 |
|------|----------|
| 연결 상태 | "🟢 Connected" 표시 |
| 초기 상태 | status: idle |
| Start 클릭 후 | status: running 으로 변경 |
| Stop 클릭 후 | status: idle 로 변경 |
| 온도 표시 | 2초마다 실시간 업데이트 |
| 메시지 로그 | 수신(←)/송신(→) 메시지 기록 |
| Clear 버튼 | 로그 목록 초기화 |

### 8.3 재연결 테스트

```bash
# 1. Broker 중지
make stop-broker

# 2. Frontend에서 Disconnected 상태 확인
mcp__playwright__playwright_get_visible_text
# Expected: "🔴 Disconnected"

# 3. Broker 재시작
make run-broker

# 4. 자동 재연결 확인 (수 초 대기 후)
mcp__playwright__playwright_get_visible_text
# Expected: "🟢 Connected"
```
