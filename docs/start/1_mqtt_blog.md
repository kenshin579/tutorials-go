# MQTT v5 완벽 가이드: 초보 개발자를 위한 스터디 자료

> 이 문서는 MQTT를 처음 접하는 개발자를 위해 작성되었습니다.
> MQTT v5만 다루며, 실무에서 반드시 필요한 재연결까지 포함합니다.

---

## 목차

- [0장. 스터디 목표와 전제](#0장-스터디-목표와-전제)
- [1장. MQTT v5 개요](#1장-mqtt-v5-개요)
- [2장. MQTT v5 기본 아키텍처](#2장-mqtt-v5-기본-아키텍처)
- [3장. Topic 설계](#3장-topic-설계)
- [4장. MQTT v5 메시지 모델](#4장-mqtt-v5-메시지-모델)
- [5장. QoS 완전 정복](#5장-qos-완전-정복)
- [6장. Session & 연결 관리](#6장-session--연결-관리)
- [7장. 재연결(Reconnect) 전략](#7장-재연결reconnect-전략)
- [8장. MQTT v5 고급 기능](#8장-mqtt-v5-고급-기능)
- [9장. 보안](#9장-보안)
- [10장. Go + Paho (v5) 사용법](#10장-go--paho-v5-사용법)
- [11장. 운영 관점 MQTT v5](#11장-운영-관점-mqtt-v5)
- [12장. MQTT v5 사용 판단 기준](#12장-mqtt-v5-사용-판단-기준)
- [13장. 스터디 마무리](#13장-스터디-마무리)

---

## 0장. 스터디 목표와 전제

### 이 스터디에서 다루는 것

1. **MQTT v5만 사용합니다**
   - v3.1.1은 다루지 않습니다
   - v5에서 추가된 기능들이 실무에서 매우 유용하기 때문입니다

2. **Pub/Sub 시스템의 설계와 운영 관점을 이해합니다**
   - 단순히 "메시지 보내고 받기"가 아닙니다
   - 시스템이 어떻게 동작하는지 전체 그림을 봅니다

3. **"네트워크는 반드시 끊긴다"를 전제로 설계합니다**
   - 이상적인 환경만 가정하지 않습니다
   - 현실 세계에서는 네트워크가 언제든 끊길 수 있습니다

4. **재연결은 옵션이 아니라 기본 기능입니다**
   - 끊어졌을 때 어떻게 복구할지가 핵심입니다

### 왜 MQTT를 배워야 할까요?

여러분이 만드는 시스템에서 이런 요구사항이 있다면 MQTT가 좋은 선택입니다:

- 수천~수만 개의 디바이스가 동시에 연결되어야 함
- 실시간으로 데이터를 주고받아야 함
- 네트워크 환경이 불안정함 (모바일, IoT)
- 서버가 클라이언트에게 먼저 메시지를 보내야 함 (Push)

---

## 1장. MQTT v5 개요

### 1.1 MQTT란 무엇인가

#### 한 줄 정의

> MQTT는 **이벤트 기반 메시징 프로토콜**입니다.

조금 더 풀어서 설명하면:
- **이벤트 기반**: "무언가 일어났을 때" 메시지를 보냅니다
- **메시징**: 데이터를 주고받는 방식입니다
- **프로토콜**: 통신 규칙입니다

#### Broker 중심 구조

MQTT의 가장 큰 특징은 **Broker**가 중간에 있다는 것입니다.

```
[Publisher] ──publish──> [Broker] ──deliver──> [Subscriber]
                            │
                            ├──deliver──> [Subscriber]
                            │
                            └──deliver──> [Subscriber]
```

- **Publisher**: 메시지를 보내는 쪽
- **Subscriber**: 메시지를 받는 쪽
- **Broker**: 메시지를 중계하는 서버

Publisher와 Subscriber는 서로를 직접 알지 못합니다. Broker만 알면 됩니다.

#### HTTP와의 근본적인 차이

| 구분 | HTTP | MQTT |
|------|------|------|
| 통신 방식 | 요청-응답 (Request-Response) | 발행-구독 (Publish-Subscribe) |
| 연결 | 요청할 때마다 연결 | 한 번 연결하면 유지 |
| 방향 | 클라이언트 → 서버 | 양방향 가능 |
| 적합한 경우 | 웹 페이지, REST API | IoT, 실시간 알림, 채팅 |

**비유로 이해하기**

- **HTTP**: 전화 통화와 비슷합니다. 할 말이 있을 때 전화를 걸고, 통화가 끝나면 끊습니다.
- **MQTT**: 라디오 방송과 비슷합니다. 라디오를 한 번 켜두면 방송이 나올 때마다 듣습니다. 여러 채널(Topic)을 선택할 수 있습니다.

### 1.2 MQTT v5가 해결하려는 문제

#### v3의 한계

MQTT v3.1.1은 오랫동안 사용되었지만 몇 가지 한계가 있었습니다:

1. **왜 실패했는지 알 수 없음**
   - 연결이 끊어져도 이유를 모름
   - 메시지 전송 실패 시 원인 파악 불가

2. **메타데이터 전달 불가**
   - 메시지 본문(Payload) 외에 추가 정보를 보낼 방법이 없음

3. **확장성 부족**
   - 대규모 시스템에서 필요한 기능들이 부족

#### v5에서 추가된 핵심 기능

**1. Reason Code (관측 가능성)**

모든 응답에 "왜?"를 알려줍니다.

```
v3: 연결 실패 (이유 모름)
v5: 연결 실패 - Reason Code 134 (Bad User Name or Password)
```

**2. User Properties (확장성)**

메시지에 원하는 정보를 추가할 수 있습니다.

```
메시지 본문: {"temperature": 25}
User Properties:
  - device_id: "sensor-001"
  - trace_id: "abc-123"
  - version: "1.0"
```

**3. 기타 유용한 기능들**
- Message Expiry Interval: 메시지 유효 시간
- Request/Response 패턴 지원
- Shared Subscription: 로드 밸런싱

---

## 2장. MQTT v5 기본 아키텍처

### 2.1 구성 요소

#### Client (클라이언트)

MQTT 시스템에서 메시지를 보내거나 받는 모든 것이 Client입니다.

- 센서, 모바일 앱, 서버 애플리케이션 모두 Client
- Publisher이면서 동시에 Subscriber일 수 있음
- Broker에 연결해서 동작

```
[모바일 앱] ──── Client
[온도 센서] ──── Client
[백엔드 서버] ── Client
```

#### Broker (브로커)

메시지를 중계하는 서버입니다.

**주요 역할:**
1. Client 연결 관리
2. Topic별 메시지 라우팅
3. 메시지 저장 (필요시)
4. 인증/인가

**대표적인 Broker:**
- **Mosquitto**: 가볍고 무료, 학습/소규모에 적합
- **EMQX**: 대규모 상용, 클러스터링 지원
- **HiveMQ**: 엔터프라이즈급

#### Topic (토픽)

메시지를 분류하는 "주소"입니다.

```
home/livingroom/temperature
home/livingroom/humidity
home/bedroom/temperature
```

- 슬래시(/)로 계층 구분
- 대소문자 구분함
- UTF-8 문자열

### 2.2 메시지 흐름

#### Publish 흐름 (메시지 보내기)

```
1. Publisher가 Broker에 연결
2. Publisher가 Topic과 함께 메시지 발행
3. Broker가 메시지를 수신하고 저장
4. Broker가 해당 Topic 구독자들에게 전달
```

```
[센서] ──PUBLISH──> [Broker]
         topic: home/temp
         payload: 25
```

#### Subscribe 흐름 (메시지 받기)

```
1. Subscriber가 Broker에 연결
2. Subscriber가 관심 있는 Topic 구독 요청
3. Broker가 구독 정보 저장
4. 해당 Topic에 메시지가 오면 전달
```

```
[앱] ──SUBSCRIBE──> [Broker]
      topic: home/temp

[Broker] ──MESSAGE──> [앱]
          topic: home/temp
          payload: 25
```

#### Broker 내부 역할

```
                    ┌─────────────────────┐
                    │       Broker        │
                    │                     │
[Publisher] ───────>│  1. 메시지 수신     │
                    │  2. Topic 매칭      │
                    │  3. 구독자 탐색     │───────> [Subscriber A]
                    │  4. 메시지 전달     │───────> [Subscriber B]
                    │  5. ACK 처리        │───────> [Subscriber C]
                    │                     │
                    └─────────────────────┘
```

Broker는 "똑똑한 우체국"과 같습니다:
- 보낸 사람이 받는 사람을 몰라도 됨
- 주소(Topic)만 보고 배달
- 같은 주소로 여러 명에게 동시 배달 가능

---

## 3장. Topic 설계

Topic 설계는 MQTT 시스템의 **가장 중요한 설계 결정**입니다.

### 3.1 Topic 구조와 규칙

#### 계층적 네이밍

Topic은 슬래시(/)로 계층을 나눕니다.

```
# 좋은 예
company/building/floor/room/sensor_type

# 실제 예시
acme/hq/3f/meeting-room-a/temperature
acme/hq/3f/meeting-room-a/humidity
acme/hq/3f/meeting-room-b/temperature
```

#### 네이밍 규칙

1. **소문자 사용 권장**
   ```
   home/temperature (O)
   Home/Temperature (비추천 - 대소문자 구분되어 혼란)
   ```

2. **공백 대신 하이픈 또는 언더스코어**
   ```
   meeting-room-a (O)
   meeting_room_a (O)
   meeting room a (X)
   ```

3. **의미 있는 순서**
   ```
   # 일반적인 순서: 큰 범위 → 작은 범위 → 타입
   지역/건물/층/방/센서종류
   ```

### 3.2 Wildcard

Wildcard는 **Subscribe할 때만** 사용할 수 있습니다. Publish할 때는 사용할 수 없습니다.

#### + (Single-Level Wildcard)

한 단계만 대체합니다.

```
구독: home/+/temperature

매칭됨:
  home/livingroom/temperature ✓
  home/bedroom/temperature ✓

매칭 안됨:
  home/floor1/room1/temperature ✗ (2단계)
  home/temperature ✗ (0단계)
```

#### # (Multi-Level Wildcard)

해당 위치부터 모든 하위 레벨을 대체합니다.

```
구독: home/#

매칭됨:
  home/livingroom ✓
  home/livingroom/temperature ✓
  home/floor1/room1/sensor1 ✓

주의:
  # 는 반드시 마지막에만 올 수 있음
  home/#/temperature (X) - 잘못된 사용
```

#### 왜 Subscribe 전용인가?

Publish할 때 Wildcard를 쓸 수 있다면?

```
# 만약 이게 가능하다면...
PUBLISH topic: home/+/temperature, payload: 25

# Broker 입장에서 어디로 보내야 할지 모름!
# - home/livingroom/temperature?
# - home/bedroom/temperature?
# - 둘 다?
```

Wildcard는 "여러 곳에서 받겠다"는 의미이지, "여러 곳에 보내겠다"는 의미가 아닙니다.

### 3.3 Topic 설계 Best Practice

#### Command / Event / State 분리

메시지의 성격에 따라 Topic을 분리하세요.

```
# Command: 명령 (누군가가 해야 할 일)
device/light-001/cmd/turn_on
device/light-001/cmd/set_brightness

# Event: 이벤트 (발생한 일)
device/light-001/event/button_pressed
device/light-001/event/error_occurred

# State: 상태 (현재 상태)
device/light-001/state/power
device/light-001/state/brightness
```

**왜 분리해야 할까요?**

- Command는 반드시 처리되어야 함 → QoS 1 이상
- Event는 놓쳐도 될 수 있음 → QoS 0 가능
- State는 최신 값만 중요 → Retained Message 사용

#### 버전 관리 전략

API처럼 Topic에도 버전을 넣을 수 있습니다.

```
# 버전 포함
v1/device/sensor-001/temperature
v2/device/sensor-001/temperature

# 또는 시스템 수준에서
mycompany/v1/device/sensor-001/temperature
```

**언제 버전을 올릴까요?**
- Payload 구조가 바뀔 때
- 의미가 바뀔 때
- Breaking change가 있을 때

#### 과도한 Wildcard의 문제

```
# 위험한 구독
구독: #

# 이러면 Broker의 모든 메시지를 받음
# - 부하 증가
# - 보안 문제
# - 처리할 수 없는 양의 데이터
```

**원칙**: 필요한 만큼만 구독하세요.

```
# 나쁜 예: 너무 넓은 구독
home/#

# 좋은 예: 필요한 것만 구독
home/livingroom/temperature
home/livingroom/humidity
```

---

## 4장. MQTT v5 메시지 모델

### 4.1 Payload

Payload는 메시지의 **본문**입니다. 실제 데이터가 들어갑니다.

#### JSON 형식

가장 많이 사용되는 형식입니다.

```json
{
  "temperature": 25.5,
  "humidity": 60,
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**장점:**
- 사람이 읽기 쉬움
- 디버깅 용이
- 유연한 구조

**단점:**
- 크기가 큼
- 파싱 비용

#### Binary 형식

배터리로 동작하는 소형 IoT 기기에서 사용합니다.

```
// 예: 온도 25.5를 2바이트로 표현
[0x00, 0xFF] // 255 = 25.5 * 10
```

**장점:**
- 크기가 작음
- 파싱 빠름

**단점:**
- 사람이 읽기 어려움
- 스키마 관리 필요

#### Schema 없는 통신의 책임

MQTT는 Payload의 형식을 강제하지 않습니다.

```
// Broker 입장에서 이 둘은 동일하게 처리됨
{"temp": 25}
hello world
0x00 0x01 0x02
```

**따라서 개발자가 책임져야 할 것들:**
1. Publisher와 Subscriber가 같은 형식 사용
2. 버전 관리
3. 유효성 검증
4. 문서화

### 4.2 User Properties

v5에서 추가된 기능으로, 메시지에 **메타데이터**를 추가할 수 있습니다.

#### 메타데이터 전달

```
Payload: {"temperature": 25}

User Properties:
  content-type: application/json
  device-id: sensor-001
  firmware-version: 1.2.3
```

Payload를 건드리지 않고 추가 정보를 전달할 수 있습니다.

#### Correlation 정보

Request/Response 패턴에서 요청과 응답을 매칭할 때 사용합니다.

```
# Request
User Properties:
  correlation-id: req-12345

# Response
User Properties:
  correlation-id: req-12345  # 같은 ID로 매칭
```

#### Trace ID 전달 패턴

분산 시스템에서 로그 추적에 유용합니다.

```
# 모든 메시지에 Trace ID 포함
User Properties:
  trace-id: abc-xyz-123
  span-id: span-456

# 로그에서 이 ID로 전체 흐름 추적 가능
[sensor] trace-id=abc-xyz-123 -> Published temperature
[broker] trace-id=abc-xyz-123 -> Routing to 3 subscribers
[server] trace-id=abc-xyz-123 -> Received and processed
```

### 4.3 Message Expiry Interval

메시지의 **유효 시간(TTL)**을 설정합니다.

#### TTL 개념

```
PUBLISH
  topic: alert/fire
  payload: "Fire detected!"
  message_expiry_interval: 60  # 60초 후 만료
```

**동작 방식:**
1. Publisher가 메시지 발행 시 TTL 설정
2. Broker가 메시지를 저장할 때 타이머 시작
3. TTL 내에 전달되지 않으면 메시지 삭제
4. Subscriber가 받을 때 남은 시간 확인 가능

#### 늦게 도착한 메시지 처리 전략

```
# 상황: Subscriber가 오프라인이었다가 연결됨
# 5분 전 메시지가 저장되어 있음

# 옵션 1: TTL로 자동 만료
message_expiry_interval: 300  # 5분

# 옵션 2: Subscriber에서 판단
if (now - message.timestamp) > threshold:
    discard(message)

# 옵션 3: 항상 최신 값만 사용 (Retained)
# 이전 메시지 무시하고 마지막 값만 사용
```

**Best Practice:**
- 실시간 알림: 짧은 TTL (10~60초)
- 상태 업데이트: 중간 TTL (5~30분)
- 중요 명령: TTL 없음 또는 매우 긴 TTL

---

## 5장. QoS 완전 정복

QoS(Quality of Service)는 메시지 **전달 보장 수준**입니다.

### 5.1 QoS 0 / 1 / 2 동작 원리

#### QoS 0: At Most Once (최대 한 번)

"보내고 잊어버린다" 방식입니다.

```
[Publisher] ──PUBLISH──> [Broker] ──PUBLISH──> [Subscriber]
             (끝)                    (끝)
```

**특징:**
- 가장 빠름
- 메시지 유실 가능
- ACK 없음

**비유**: 엽서 보내기 - 보냈는지 확인 안 함

#### QoS 1: At Least Once (최소 한 번)

"받았다고 확인할 때까지 재전송" 방식입니다.

```
[Publisher] ──PUBLISH──> [Broker]
[Publisher] <──PUBACK─── [Broker]  # ACK 받으면 끝

[Broker] ──PUBLISH──> [Subscriber]
[Broker] <──PUBACK─── [Subscriber]  # ACK 받으면 끝
```

**특징:**
- 메시지 전달 보장
- 중복 가능 (ACK 유실 시 재전송)
- 가장 많이 사용됨

**비유**: 등기 우편 - 받았다는 확인 필요

#### QoS 2: Exactly Once (정확히 한 번)

"중복 없이 정확히 한 번 전달" 방식입니다.

```
[Publisher] ──PUBLISH──> [Broker]
[Publisher] <──PUBREC─── [Broker]  # 받았음
[Publisher] ──PUBREL──> [Broker]   # 삭제해도 됨
[Publisher] <──PUBCOMP── [Broker]  # 완료

# Broker → Subscriber도 동일한 4단계
```

**특징:**
- 중복 없음 보장
- 가장 느림 (4번의 핸드셰이크)
- 거의 사용되지 않음

**비유**: 은행 송금 - 정확히 한 번만 실행되어야 함

#### 한눈에 비교

| QoS | 이름 | 전달 보장 | 중복 가능 | 속도 |
|-----|------|-----------|-----------|------|
| 0 | At Most Once | X | X | 빠름 |
| 1 | At Least Once | O | O | 보통 |
| 2 | Exactly Once | O | X | 느림 |

### 5.2 QoS 선택 기준

#### 상태 보고: QoS 0 또는 1

```
# 예: 온도 센서가 1초마다 값 전송
topic: sensor/temp
payload: 25.5
qos: 0  # 하나쯤 놓쳐도 다음 값이 옴
```

**판단 기준:**
- 주기적으로 전송됨 → QoS 0
- 가끔 전송되고 중요함 → QoS 1

#### 이벤트: QoS 1

```
# 예: 문 열림 이벤트
topic: door/event/opened
payload: {"time": "10:30:00"}
qos: 1  # 이벤트는 놓치면 안 됨
```

이벤트는 보통 한 번 발생하면 끝이므로 놓치면 복구가 어렵습니다.

#### 명령: QoS 1 또는 2

```
# 예: 조명 끄기 명령
topic: light/cmd/off
payload: {}
qos: 1  # 반드시 전달되어야 함
```

**중복 실행이 문제가 되는 경우:**
```
# 예: 결제 요청
topic: payment/process
payload: {"amount": 10000}
qos: 2  # 정확히 한 번만 실행
# 또는 QoS 1 + Idempotent 처리
```

### 5.3 QoS와 중복 처리

#### At-Least-Once의 현실

QoS 1을 사용하면 중복이 발생할 수 있습니다.

```
# 시나리오
1. Publisher가 메시지 전송
2. Broker가 받고 저장
3. Broker가 PUBACK 전송
4. 네트워크 문제로 PUBACK 유실
5. Publisher가 메시지 재전송 (타임아웃)
6. Broker가 같은 메시지를 또 받음 → 중복!
```

#### Idempotent Consumer 설계

중복 메시지를 받아도 문제없게 설계하는 것이 **멱등성(Idempotency)**입니다.

**방법 1: 메시지 ID로 중복 체크**
```go
func handleMessage(msg Message) {
    // 이미 처리한 메시지인지 확인
    if processed[msg.ID] {
        return  // 무시
    }

    processMessage(msg)
    processed[msg.ID] = true
}
```

**방법 2: 상태 기반 처리**
```go
// 나쁜 예: 잔액 증가 (중복되면 문제)
balance += amount

// 좋은 예: 상태 설정 (중복되어도 같은 결과)
balance = newBalance
status = "completed"
```

**방법 3: 타임스탬프 활용**
```go
func handleState(msg StateMessage) {
    // 오래된 메시지는 무시
    if msg.Timestamp < lastTimestamp {
        return
    }

    updateState(msg)
    lastTimestamp = msg.Timestamp
}
```

---

## 6장. Session & 연결 관리

### 6.1 Session Expiry Interval

세션은 Client와 Broker 간의 **연결 상태 정보**입니다.

#### Clean Start vs Session 유지

**Clean Start = true (새 세션)**
```
연결 시:
  - 이전 세션 정보 삭제
  - 구독 정보 초기화
  - 저장된 메시지 삭제

사용 케이스:
  - 임시 연결
  - 상태가 필요 없는 Publisher
```

**Clean Start = false (세션 유지)**
```
연결 시:
  - 이전 세션 정보 복원
  - 구독 정보 유지
  - 오프라인 동안의 메시지 전달

사용 케이스:
  - 지속적인 구독자
  - 메시지를 놓치면 안 되는 경우
```

#### Session Expiry Interval

세션을 얼마나 유지할지 설정합니다.

```go
// 세션 설정 예시
SessionExpiryInterval: 3600  // 1시간

// 동작
1. Client 연결 끊김
2. Broker가 1시간 동안 세션 유지
3. 1시간 내 재연결 → 세션 복원, 밀린 메시지 전달
4. 1시간 후 재연결 → 새 세션 시작
```

**권장 값:**
- 모바일 앱: 1-24시간
- IoT 기기: 필요에 따라 (분~일)
- 임시 연결: 0 (세션 유지 안 함)

#### 오프라인 메시지

Session이 유지되는 동안 Broker가 메시지를 저장합니다.

```
1. Subscriber가 오프라인
2. Publisher가 메시지 발행 (QoS 1)
3. Broker가 메시지 저장 (Subscriber 세션이 살아있으므로)
4. Subscriber가 재연결
5. Broker가 저장된 메시지 전달
```

**주의사항:**
- QoS 0 메시지는 저장되지 않음
- 저장 용량에 제한이 있을 수 있음
- Session Expiry 전에 재연결해야 함

### 6.2 Keep Alive

연결이 살아있는지 확인하는 메커니즘입니다.

#### Ping 메커니즘

```
Keep Alive = 60초로 설정

[Client] ──PINGREQ──> [Broker]  # 60초 동안 통신 없으면
[Client] <──PINGRESP── [Broker]

# 응답 없으면 연결 끊김으로 판단
```

**동작 방식:**
1. Client가 Keep Alive 간격 설정 (예: 60초)
2. 해당 시간 동안 메시지가 없으면 PINGREQ 전송
3. Broker가 PINGRESP로 응답
4. Keep Alive * 1.5 시간 내 응답 없으면 연결 종료

#### 네트워크 품질과의 관계

```
# 안정적인 네트워크
keep_alive: 60~120초

# 불안정한 네트워크 (모바일, IoT)
keep_alive: 15~30초
# 더 자주 체크하지만 오버헤드 증가

# 매우 안정적인 환경 (데이터센터 내)
keep_alive: 300초 이상
```

**Trade-off:**
- 짧은 Keep Alive: 빠른 끊김 감지, 높은 오버헤드
- 긴 Keep Alive: 낮은 오버헤드, 느린 끊김 감지

### 6.3 Retained Message

Topic에 **마지막 메시지를 저장**하는 기능입니다.

#### Last Known State 패턴

```
# 온도 센서가 Retained 메시지 발행
PUBLISH
  topic: sensor/temperature
  payload: 25
  retain: true

# Broker가 이 메시지를 저장

# 나중에 새 Subscriber가 구독하면
SUBSCRIBE topic: sensor/temperature
# → 즉시 마지막 값(25)을 받음
```

**왜 유용한가:**
- 새로 연결한 Client도 현재 상태를 즉시 알 수 있음
- 센서가 자주 전송하지 않아도 됨
- "현재 상태가 뭐야?" 질문에 답할 수 있음

#### 오용 사례

```
# 나쁜 사용: 이벤트에 Retain
PUBLISH
  topic: door/event/opened
  payload: {"time": "10:30:00"}
  retain: true  # 잘못됨!

# 문제: 새 구독자가 "문이 열렸다"는 과거 이벤트를 받음
# 현재 문 상태인지, 과거 이벤트인지 구분 불가
```

**Retain을 써야 하는 경우:**
- 상태 (온도, 습도, 전원 상태)
- 설정 값
- 현재 위치

**Retain을 쓰면 안 되는 경우:**
- 이벤트 (버튼 클릭, 문 열림)
- 명령
- 로그

---

## 7장. 재연결(Reconnect) 전략

> 이 장은 실무에서 **가장 중요한** 부분입니다.

### 7.1 재연결이 반드시 필요한 이유

#### 현실 세계의 네트워크

이상적인 세계에서는 한 번 연결하면 영원히 유지됩니다.
하지만 현실은 다릅니다:

```
# 네트워크 끊김 원인들
- Wi-Fi → LTE 전환 (모바일)
- 터널, 엘리베이터 (모바일)
- 라우터 재시작
- ISP 장애
- Broker 재시작
- 로드밸런서 타임아웃
- 메모리 부족으로 인한 강제 종료
```

#### 환경별 특성

**모바일**
```
- 수시로 네트워크 전환
- 백그라운드 진입 시 OS가 연결 끊음
- 배터리 절약으로 인한 제한
```

**로봇/차량**
```
- 이동 중 기지국 전환
- 음영 지역 통과
- 하드웨어 재부팅
```

**IoT 센서**
```
- 전원 불안정
- 무선 간섭
- 펌웨어 업데이트로 재시작
```

#### Broker 장애

Broker도 죽을 수 있습니다:
```
- 메모리 부족
- 디스크 가득 참
- 업그레이드/패치
- 하드웨어 장애
```

**결론**: 재연결은 "만약"이 아니라 "언제" 발생하느냐의 문제입니다.

### 7.2 재연결 시 발생하는 문제들

#### 구독 유실

Clean Start 설정에 따라 구독이 사라질 수 있습니다.

```
# 시나리오
1. Client가 topic/a, topic/b 구독 중
2. 연결 끊김
3. Session Expiry 지남 또는 Clean Start=true로 재연결
4. 구독 정보 사라짐
5. 메시지를 못 받음!
```

#### 중복 메시지

재연결 시점에 따라 같은 메시지를 여러 번 받을 수 있습니다.

```
# 시나리오
1. Broker가 메시지 전송
2. Client가 받았지만 ACK 전송 전 연결 끊김
3. 재연결
4. Broker가 ACK 못 받았으므로 재전송
5. 같은 메시지 2번 받음
```

#### 메시지 순서 깨짐

```
# 시나리오
1. 메시지 A 전송됨
2. 연결 끊김
3. 메시지 B, C가 Broker에 저장됨
4. 재연결
5. 저장된 B, C가 먼저 옴
6. 순서: B → C → D (A는 이미 처리됨)

# 문제: A 처리 후 연결 끊기 전에 온 메시지는?
```

### 7.3 재연결 설계 전략

#### Auto Reconnect

대부분의 MQTT 클라이언트 라이브러리는 자동 재연결을 지원합니다.

```go
// Paho v5 예시
config := autopaho.ClientConfig{
    ConnectRetryDelay: 10 * time.Second,  // 재시도 간격
    // ...
}
```

**자동 재연결이 하는 일:**
1. 연결 끊김 감지
2. 일정 시간 대기
3. 재연결 시도
4. 실패하면 다시 대기 후 재시도

#### Backoff 전략

재연결 실패 시 대기 시간을 점점 늘리는 전략입니다.

```
# Fixed Backoff (고정)
시도 1: 1초 대기
시도 2: 1초 대기
시도 3: 1초 대기
...

# Exponential Backoff (지수)
시도 1: 1초 대기
시도 2: 2초 대기
시도 3: 4초 대기
시도 4: 8초 대기
...

# Exponential Backoff with Jitter (+ 랜덤)
시도 1: 1초 + random(0~500ms)
시도 2: 2초 + random(0~500ms)
...
```

**왜 Jitter가 필요한가:**
```
# 시나리오: Broker 재시작
1. 1000개 Client가 동시에 끊김
2. 모두 1초 후 재연결 시도
3. Broker에 1000개 연결 요청 폭주
4. Broker 과부하

# Jitter 적용 시
1. 1000개 Client가 동시에 끊김
2. 각자 1초 + 랜덤 시간 후 재연결
3. 연결 요청이 분산됨
4. Broker 안정적 처리
```

#### Session 유지 vs 초기화

```go
// Session 유지 (권장)
CleanStart: false
SessionExpiryInterval: 3600  // 1시간

// 장점:
// - 구독 정보 유지
// - 오프라인 메시지 받음

// Session 초기화
CleanStart: true

// 필요한 경우:
// - 완전히 새로 시작해야 할 때
// - 문제가 발생해서 리셋할 때
```

### 7.4 재연결 후 처리 로직

#### 재구독 전략

Session이 만료되었거나 Clean Start를 사용한 경우, 재구독이 필요합니다.

```go
// 재연결 성공 시 콜백
func onConnect(client *paho.Client) {
    // 필요한 Topic들 재구독
    topics := []string{
        "device/+/state",
        "command/mydevice/#",
    }

    for _, topic := range topics {
        client.Subscribe(topic, qos)
    }
}
```

**Best Practice: 구독 목록 관리**
```go
type SubscriptionManager struct {
    subscriptions map[string]byte  // topic -> qos
}

func (sm *SubscriptionManager) Resubscribe(client *paho.Client) {
    for topic, qos := range sm.subscriptions {
        client.Subscribe(topic, qos)
    }
}
```

#### 미처리 메시지 처리

재연결 후 밀린 메시지를 받을 때 고려사항:

```go
func onMessage(msg Message) {
    // 1. 메시지 나이 확인
    age := time.Since(msg.Timestamp)
    if age > maxMessageAge {
        log.Warn("Discarding old message", age)
        return
    }

    // 2. 중복 확인
    if isProcessed(msg.ID) {
        return
    }

    // 3. 처리
    processMessage(msg)
    markAsProcessed(msg.ID)
}
```

#### 상태 동기화 패턴

재연결 후 현재 상태를 동기화하는 패턴입니다.

**방법 1: Retained Message 활용**
```
# 구독하면 마지막 상태 즉시 수신
SUBSCRIBE topic: device/+/state
→ 각 디바이스의 마지막 상태 수신
```

**방법 2: 명시적 상태 요청**
```
# 재연결 후 상태 요청
PUBLISH topic: device/mydevice/cmd/get_state
→ 디바이스가 현재 상태 응답
```

**방법 3: 시퀀스 번호 기반**
```go
// 마지막 처리한 시퀀스 저장
lastSequence := loadLastSequence()

// 재연결 후
for _, msg := range messages {
    if msg.Sequence <= lastSequence {
        continue  // 이미 처리함
    }
    processMessage(msg)
    saveLastSequence(msg.Sequence)
}
```

---

## 8장. MQTT v5 고급 기능

### 8.1 Shared Subscription

여러 Subscriber가 **메시지를 나눠서** 처리하는 기능입니다.

#### 개념

```
# 일반 구독
[Subscriber A] ←─ message ─┐
[Subscriber B] ←─ message ─┤ Broker (같은 메시지를 모두에게)
[Subscriber C] ←─ message ─┘

# Shared Subscription
[Subscriber A] ←─ message 1 ─┐
[Subscriber B] ←─ message 2 ─┤ Broker (메시지를 분배)
[Subscriber C] ←─ message 3 ─┘
```

**사용 방법:**
```
# Topic 앞에 $share/그룹명/ 추가
$share/mygroup/sensor/temperature

# 같은 그룹 내에서 메시지 분배
# 다른 그룹은 독립적으로 모든 메시지 수신
```

#### 로드 분산 vs 순서 보장

**로드 분산 관점:**
```
처리량 = 단일 Subscriber 처리량 × Subscriber 수
```

**순서 보장 관점:**
```
# 문제
Message 1 → Subscriber A
Message 2 → Subscriber B
Message 3 → Subscriber A

# Subscriber A 입장에서 순서: 1, 3 (2가 없음)
# 전체 순서 보장 안 됨!
```

#### 언제 쓰면 안 되는가

```
# 쓰면 안 되는 경우
1. 메시지 순서가 중요할 때
   - 거래 처리
   - 상태 변경 추적

2. 상태를 유지해야 할 때
   - 특정 디바이스의 모든 메시지를 한 곳에서 처리

# 적합한 경우
1. 독립적인 메시지 처리
   - 로그 수집
   - 이미지 처리
   - 알림 발송
```

### 8.2 Request / Response 패턴

MQTT로 HTTP처럼 요청-응답을 구현하는 패턴입니다.

#### Response Topic

응답을 받을 Topic을 요청에 포함시킵니다.

```
# 요청
PUBLISH
  topic: device/cmd/get_status
  response_topic: reply/client-123/status
  correlation_data: req-001

# 응답
PUBLISH
  topic: reply/client-123/status  # 요청의 response_topic
  correlation_data: req-001        # 요청의 correlation_data
  payload: {"status": "ok"}
```

#### Correlation Data

요청과 응답을 매칭하는 데이터입니다.

```go
// 요청 보내기
request := &paho.Publish{
    Topic:   "device/cmd",
    Payload: []byte("get_status"),
    Properties: &paho.PublishProperties{
        ResponseTopic:   "reply/my-client",
        CorrelationData: []byte("req-12345"),
    },
}

// 응답 받기 (reply/my-client 구독 중)
func onMessage(msg Message) {
    correlationID := string(msg.Properties.CorrelationData)
    // correlationID == "req-12345" 로 매칭
}
```

#### Timeout 처리

응답이 안 오면 어떻게 할까요?

```go
func requestWithTimeout(request Message, timeout time.Duration) (Response, error) {
    // 1. 응답 채널 생성
    responseChan := make(chan Response)
    pending[request.CorrelationID] = responseChan

    // 2. 요청 전송
    client.Publish(request)

    // 3. 응답 또는 타임아웃 대기
    select {
    case resp := <-responseChan:
        return resp, nil
    case <-time.After(timeout):
        delete(pending, request.CorrelationID)
        return Response{}, ErrTimeout
    }
}
```

**Best Practice:**
- 항상 타임아웃 설정
- 타임아웃 시 재시도 또는 에러 처리
- 오래된 응답은 무시

### 8.3 Reason Code

모든 응답에 **성공/실패 이유**를 포함합니다.

#### 성공/실패 세분화

```
# 연결 응답 Reason Code 예시
0   = Success
128 = Unspecified error
129 = Malformed Packet
130 = Protocol Error
131 = Implementation specific error
132 = Unsupported Protocol Version
133 = Client Identifier not valid
134 = Bad User Name or Password
135 = Not authorized
```

#### 장애 원인 파악

```go
// 연결 실패 시
func onConnectError(err error, reasonCode byte) {
    switch reasonCode {
    case 134:
        log.Error("인증 실패: 사용자명/비밀번호 확인 필요")
    case 135:
        log.Error("권한 없음: ACL 설정 확인 필요")
    case 137:
        log.Error("서버 사용 불가: 잠시 후 재시도")
    default:
        log.Error("연결 실패", reasonCode)
    }
}
```

**v3과의 차이:**
```
v3: "연결 실패" (왜?)
v5: "연결 실패 - Reason Code 134: Bad User Name or Password"
```

---

## 9장. 보안

### 9.1 인증

Client가 **누구인지** 확인합니다.

#### Username / Password

가장 기본적인 인증 방식입니다.

```go
// 연결 시 인증 정보 제공
config := paho.Connect{
    ClientID: "my-device",
    Username: "device-001",
    Password: []byte("secret-password"),
}
```

**주의사항:**
- 평문으로 전송됨 (TLS 필수)
- 비밀번호 관리 필요
- 각 디바이스별 고유 인증 정보 권장

#### Token 기반 인증

JWT 등의 토큰을 사용하는 방식입니다.

```go
// JWT 토큰을 Password로 사용
token := generateJWT(deviceID, expiry)
config := paho.Connect{
    ClientID: "my-device",
    Username: "jwt",
    Password: []byte(token),
}
```

**장점:**
- 만료 시간 설정 가능
- 추가 정보 포함 가능 (권한 등)
- 비밀번호 저장 불필요

### 9.2 인가

인증된 Client가 **무엇을 할 수 있는지** 결정합니다.

#### Topic 기반 ACL

```
# ACL 설정 예시 (Mosquitto)
user sensor-001
topic read sensor/+/state     # 읽기만 가능
topic write sensor/001/#      # 자기 토픽만 쓰기

user admin
topic readwrite #             # 모든 토픽 읽기/쓰기
```

#### Publish / Subscribe 권한 분리

```
# 센서는 자기 데이터만 발행
sensor-001:
  publish: sensor/001/data
  subscribe: command/001/#

# 대시보드는 모든 센서 데이터 구독, 명령 발행
dashboard:
  publish: command/+/#
  subscribe: sensor/+/data
```

### 9.3 TLS

통신 내용을 **암호화**합니다.

#### 언제 필요한가

```
# TLS 필수
- 인터넷을 통한 통신
- 민감한 데이터 전송
- 인증 정보 보호

# TLS 선택적
- 폐쇄망 내부 통신
- 성능이 극히 중요한 경우
```

#### 성능과 보안의 균형

```
# TLS 오버헤드
- 핸드셰이크 지연 (초기 연결)
- CPU 사용량 증가 (암호화/복호화)
- 메모리 사용량 증가

# 경량 디바이스의 경우
- TLS 1.3 사용 (더 가벼움)
- 또는 VPN으로 네트워크 레벨 보안
```

---

## 10장. Go + Paho (v5) 사용법

### 10.1 Paho v5 구조 이해

Go에서 MQTT v5를 사용하려면 `eclipse/paho.golang` 패키지를 사용합니다.

#### 주요 패키지

```go
import (
    "github.com/eclipse/paho.golang/paho"           // 기본 클라이언트
    "github.com/eclipse/paho.golang/autopaho"       // 자동 재연결
)
```

#### ClientConfig (autopaho)

```go
config := autopaho.ClientConfig{
    // Broker 주소
    BrokerUrls: []*url.URL{brokerURL},

    // Keep Alive 간격
    KeepAlive: 30,

    // 재연결 간격
    ConnectRetryDelay: 10 * time.Second,

    // 연결 성공 시 콜백
    OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
        // 구독 설정
    },

    // 연결 끊김 시 콜백
    OnConnectError: func(err error) {
        log.Error("Connection error", err)
    },
}
```

#### ConnectionManager

```go
// 연결 시작
cm, err := autopaho.NewConnection(ctx, config)

// 연결 대기
err = cm.AwaitConnection(ctx)

// 연결 종료
err = cm.Disconnect(ctx)
```

#### Handler 구조

```go
// 메시지 수신 핸들러
func messageHandler(msg *paho.Publish) {
    fmt.Printf("Topic: %s, Payload: %s\n",
        msg.Topic, string(msg.Payload))
}

// Router 설정
router := paho.NewStandardRouter()
router.RegisterHandler("sensor/#", messageHandler)
```

### 10.2 기본 사용 흐름

#### Connect (연결)

```go
package main

import (
    "context"
    "log"
    "net/url"

    "github.com/eclipse/paho.golang/autopaho"
    "github.com/eclipse/paho.golang/paho"
)

func main() {
    ctx := context.Background()

    brokerURL, _ := url.Parse("mqtt://localhost:1883")

    config := autopaho.ClientConfig{
        BrokerUrls: []*url.URL{brokerURL},
        KeepAlive:  30,

        ConnectUsername: "user",
        ConnectPassword: []byte("password"),

        ClientConfig: paho.ClientConfig{
            ClientID: "my-client",
            Router:   paho.NewStandardRouter(),
        },
    }

    cm, err := autopaho.NewConnection(ctx, config)
    if err != nil {
        log.Fatal(err)
    }

    err = cm.AwaitConnection(ctx)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Connected!")
}
```

#### Subscribe (구독)

```go
func setupSubscription(cm *autopaho.ConnectionManager, router *paho.StandardRouter) {
    // 핸들러 등록
    router.RegisterHandler("sensor/+/temperature", func(msg *paho.Publish) {
        log.Printf("Temperature: %s", msg.Payload)
    })

    // 구독 요청
    cm.Subscribe(context.Background(), &paho.Subscribe{
        Subscriptions: []paho.SubscribeOptions{
            {Topic: "sensor/+/temperature", QoS: 1},
        },
    })
}
```

#### Publish (발행)

```go
func publishMessage(cm *autopaho.ConnectionManager) {
    msg := &paho.Publish{
        Topic:   "sensor/001/temperature",
        QoS:     1,
        Payload: []byte(`{"value": 25.5}`),
        Properties: &paho.PublishProperties{
            UserProperties: []paho.UserProperty{
                {Key: "device-id", Value: "sensor-001"},
            },
        },
    }

    _, err := cm.Publish(context.Background(), msg)
    if err != nil {
        log.Error("Publish failed", err)
    }
}
```

### 10.3 재연결 구현 방식

#### 자동 재연결 설정

```go
config := autopaho.ClientConfig{
    // 재연결 간격
    ConnectRetryDelay: 10 * time.Second,

    // 최대 재연결 간격 (Backoff)
    // 기본적으로 Exponential Backoff 적용됨
}
```

#### OnConnectionUp

연결 성공 시 호출됩니다. 재구독에 사용합니다.

```go
config.OnConnectionUp = func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
    log.Println("Connected!")

    // Session이 새로 시작되었는지 확인
    if !connAck.SessionPresent {
        log.Println("Session not present, resubscribing...")
        resubscribe(cm)
    }
}

func resubscribe(cm *autopaho.ConnectionManager) {
    topics := []paho.SubscribeOptions{
        {Topic: "sensor/+/temperature", QoS: 1},
        {Topic: "command/my-device/#", QoS: 1},
    }

    cm.Subscribe(context.Background(), &paho.Subscribe{
        Subscriptions: topics,
    })
}
```

#### OnServerDisconnect

Broker가 연결을 끊었을 때 호출됩니다.

```go
config.ClientConfig.OnServerDisconnect = func(d *paho.Disconnect) {
    if d.ReasonCode != 0 {
        log.Printf("Server disconnected: reason=%d", d.ReasonCode)
    }
}
```

#### OnClientError

클라이언트 에러 발생 시 호출됩니다.

```go
config.ClientConfig.OnClientError = func(err error) {
    log.Printf("Client error: %v", err)
}
```

### 10.4 안전한 Handler 설계

#### Blocking 금지

메시지 핸들러에서 오래 걸리는 작업을 하면 안 됩니다.

```go
// 나쁜 예: 핸들러에서 직접 처리
func badHandler(msg *paho.Publish) {
    result := heavyProcessing(msg.Payload)  // 10초 걸림
    saveToDatabase(result)                   // 1초 걸림
    // 이 동안 다른 메시지 처리 못함!
}

// 좋은 예: 채널로 전달하고 바로 리턴
func goodHandler(msg *paho.Publish) {
    messageQueue <- msg  // 즉시 리턴
}

// 별도 고루틴에서 처리
go func() {
    for msg := range messageQueue {
        result := heavyProcessing(msg.Payload)
        saveToDatabase(result)
    }
}()
```

#### Worker Pool 패턴

```go
type MessageProcessor struct {
    queue   chan *paho.Publish
    workers int
}

func NewMessageProcessor(workers, queueSize int) *MessageProcessor {
    mp := &MessageProcessor{
        queue:   make(chan *paho.Publish, queueSize),
        workers: workers,
    }

    // Worker 시작
    for i := 0; i < workers; i++ {
        go mp.worker(i)
    }

    return mp
}

func (mp *MessageProcessor) worker(id int) {
    for msg := range mp.queue {
        log.Printf("Worker %d processing: %s", id, msg.Topic)
        processMessage(msg)
    }
}

func (mp *MessageProcessor) Enqueue(msg *paho.Publish) {
    select {
    case mp.queue <- msg:
        // 큐에 추가됨
    default:
        log.Warn("Queue full, dropping message")
    }
}

// 핸들러에서 사용
func handler(msg *paho.Publish) {
    processor.Enqueue(msg)
}
```

---

## 11장. 운영 관점 MQTT v5

### 11.1 모니터링 포인트

#### 연결 수

```
# 모니터링 항목
- 현재 활성 연결 수
- 연결/해제 비율 (churn rate)
- 연결 실패 수

# 경고 기준 예시
- 연결 수 급증: 1분 내 50% 이상 증가
- 연결 실패율: 1% 이상
```

#### 메시지 처리율

```
# 모니터링 항목
- 초당 수신 메시지 수 (messages/sec)
- 초당 발송 메시지 수
- 평균 메시지 크기
- 대기 중인 메시지 수

# 경고 기준 예시
- 처리율 저하: 평소 대비 30% 이상 감소
- 대기열 증가: 1000개 이상
```

#### 재연결 빈도

```
# 모니터링 항목
- 재연결 횟수 / 시간
- Client별 재연결 패턴
- 재연결 실패율

# 경고 기준 예시
- 특정 Client가 1분에 10회 이상 재연결
- 전체 재연결률 급증
```

### 11.2 장애 시나리오별 대응

#### Broker 재시작

```
# 현상
- 모든 Client 연결 끊김
- 동시 재연결 시도

# 대응
1. Client에 Exponential Backoff + Jitter 적용
2. Session Expiry 충분히 설정
3. Broker 클러스터링 고려
```

#### 네트워크 Flap

```
# 현상
- 연결/끊김 반복
- 메시지 중복 발생

# 대응
1. 재연결 간격 조정
2. Idempotent 처리
3. 회로 차단기 패턴 적용
```

#### Client 폭증

```
# 현상
- 연결 수 급증
- Broker 응답 지연
- 메모리/CPU 급증

# 대응
1. 연결 속도 제한 (rate limiting)
2. Broker 스케일 아웃
3. 불필요한 연결 정리
```

---

## 12장. MQTT v5 사용 판단 기준

### MQTT를 써야 하는 경우

1. **실시간 양방향 통신이 필요할 때**
   ```
   - 채팅
   - 실시간 알림
   - 원격 제어
   ```

2. **많은 디바이스가 연결될 때**
   ```
   - IoT 센서 네트워크
   - 스마트 홈
   - 차량 관제
   ```

3. **네트워크가 불안정할 때**
   ```
   - 모바일 환경
   - 저전력 무선
   - 원격지
   ```

4. **서버 → 클라이언트 Push가 필요할 때**
   ```
   - 상태 변경 알림
   - 명령 전달
   - 이벤트 브로드캐스트
   ```

### MQTT를 쓰면 안 되는 경우

1. **단순 요청-응답만 필요할 때**
   ```
   → HTTP/REST 사용
   ```

2. **파일 전송이 필요할 때**
   ```
   → HTTP, FTP, S3 등 사용
   MQTT는 작은 메시지에 최적화됨
   ```

3. **강력한 트랜잭션이 필요할 때**
   ```
   → 메시지 큐 (RabbitMQ, Kafka) 사용
   MQTT는 메시지 순서 보장이 약함
   ```

4. **브라우저 직접 연결이 필요할 때**
   ```
   → WebSocket 직접 사용 또는 MQTT over WebSocket
   ```

### HTTP / gRPC와의 경계

| 기준 | HTTP | gRPC | MQTT |
|------|------|------|------|
| 통신 패턴 | 요청-응답 | 요청-응답, 스트리밍 | Pub/Sub |
| 연결 | 단발성 | 지속 가능 | 지속 |
| 다수 수신자 | 어려움 | 어려움 | 쉬움 |
| 서버 Push | 폴링 필요 | 스트리밍 가능 | 기본 지원 |
| 적합한 곳 | 웹 API | 마이크로서비스 | IoT, 실시간 |

---

## 13장. 스터디 마무리

### 13.1 핵심 요약

#### MQTT v5의 본질

1. **Pub/Sub 패턴**
   - Publisher와 Subscriber가 서로 몰라도 됨
   - Broker가 메시지를 중계
   - Topic으로 메시지 분류

2. **경량 프로토콜**
   - 작은 오버헤드
   - 불안정한 네트워크에 적합
   - 저사양 디바이스 지원

3. **v5의 개선점**
   - Reason Code로 디버깅 용이
   - User Properties로 확장성
   - Shared Subscription으로 로드 분산

#### 신뢰성은 애플리케이션 책임

MQTT가 보장하는 것:
- QoS에 따른 전달 보장
- 세션 유지

MQTT가 보장하지 않는 것:
- 메시지 순서 (여러 Topic 간)
- 중복 방지 (QoS 1 사용 시)
- 비즈니스 로직 정합성

**따라서 애플리케이션에서:**
- Idempotent 처리 구현
- 타임스탬프/시퀀스 기반 정렬
- 재연결 후 상태 동기화

### 13.2 체크리스트

스터디를 마치고 이 질문들에 답할 수 있어야 합니다.

#### Topic 설계

- [ ] 내 시스템의 Topic 구조를 설계했는가?
- [ ] Command / Event / State를 구분했는가?
- [ ] Wildcard 사용이 적절한가?

#### QoS 선택

- [ ] 각 Topic의 QoS를 결정했는가?
- [ ] 그 이유를 설명할 수 있는가?
- [ ] 중복 처리 방안은 있는가?

#### 재연결 전략

- [ ] 네트워크 끊김을 고려했는가?
- [ ] 자동 재연결을 구현했는가?
- [ ] 재연결 후 재구독을 구현했는가?
- [ ] Backoff 전략을 적용했는가?

#### 세션 관리

- [ ] Session Expiry를 설정했는가?
- [ ] Clean Start 정책을 결정했는가?
- [ ] 오프라인 메시지 처리 방안은 있는가?

#### 보안

- [ ] 인증 방식을 결정했는가?
- [ ] Topic별 권한을 설정했는가?
- [ ] TLS 사용 여부를 결정했는가?

---

## 부록: 실습 환경 설정

### Mosquitto Broker 설치 (Docker)

```bash
# Mosquitto 실행
docker run -d --name mosquitto \
  -p 1883:1883 \
  -p 9001:9001 \
  eclipse-mosquitto

# 설정 파일 사용
docker run -d --name mosquitto \
  -p 1883:1883 \
  -v $(pwd)/mosquitto.conf:/mosquitto/config/mosquitto.conf \
  eclipse-mosquitto
```

### 기본 설정 파일 (mosquitto.conf)

```
listener 1883
allow_anonymous true
```

### 테스트 명령어

```bash
# 구독 (터미널 1)
mosquitto_sub -h localhost -t "test/#" -v

# 발행 (터미널 2)
mosquitto_pub -h localhost -t "test/hello" -m "Hello MQTT!"
```

### Go 의존성

```bash
go get github.com/eclipse/paho.golang@latest
```

---

## 참고 자료

- [MQTT v5 스펙](https://docs.oasis-open.org/mqtt/mqtt/v5.0/mqtt-v5.0.html)
- [Eclipse Paho Go Client](https://github.com/eclipse/paho.golang)
- [EMQX 문서](https://www.emqx.io/docs)
- [Mosquitto 문서](https://mosquitto.org/documentation/)

---

> 이 문서는 계속 업데이트됩니다.
> 질문이나 피드백은 이슈로 남겨주세요.
