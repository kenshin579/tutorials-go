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

이 스터디는 MQTT를 처음 접하는 개발자가 실무에서 바로 활용할 수 있도록 구성되었습니다. 단순히 "메시지를 보내고 받는 방법"을 넘어서, 시스템 설계와 운영 관점에서 MQTT를 이해하는 것을 목표로 합니다. 특히 네트워크 불안정성을 전제로 한 재연결 전략까지 다루므로, 학습 후에는 프로덕션 환경에서도 안정적으로 동작하는 MQTT 기반 시스템을 구축할 수 있습니다.

### 이 스터디에서 다루는 것

1. **MQTT v5만 사용합니다**
   - v3.1.1은 다루지 않습니다. v5에서 추가된 기능들이 실무에서 매우 유용하기 때문입니다.
   - v5는 2019년에 발표된 최신 버전으로, Reason Code, User Properties, Shared Subscription 등 운영에 필수적인 기능들이 대폭 추가되었습니다. 새로운 프로젝트를 시작한다면 v5를 선택하는 것이 합리적입니다.

2. **Pub/Sub 시스템의 설계와 운영 관점을 이해합니다**
   - 단순히 "메시지 보내고 받기"가 아닙니다. 시스템이 어떻게 동작하는지 전체 그림을 봅니다.
   - Topic 네이밍 규칙, QoS 선택 기준, 세션 관리 등 설계 단계에서 결정해야 할 사항들을 다룹니다. 이러한 결정들은 나중에 변경하기 어려우므로 처음부터 올바르게 설계하는 것이 중요합니다.

3. **"네트워크는 반드시 끊긴다"를 전제로 설계합니다**
   - 이상적인 환경만 가정하지 않습니다. 현실 세계에서는 네트워크가 언제든 끊길 수 있습니다.
   - 모바일 환경에서는 Wi-Fi와 LTE 사이를 오가고, IoT 디바이스는 전원이 불안정할 수 있습니다. 이러한 상황을 처음부터 고려하지 않으면 프로덕션에서 예상치 못한 문제에 직면하게 됩니다.

4. **재연결은 옵션이 아니라 기본 기능입니다**
   - 끊어졌을 때 어떻게 복구할지가 핵심입니다.
   - 많은 MQTT 튜토리얼이 연결과 메시지 전송만 다루지만, 실제 서비스에서는 재연결 로직이 코드의 상당 부분을 차지합니다. 이 스터디에서는 재연결 전략을 별도의 장으로 다룹니다.

### 왜 MQTT를 배워야 할까요?

MQTT는 1999년 IBM에서 석유 파이프라인 모니터링을 위해 개발된 프로토콜입니다. 위성 네트워크처럼 대역폭이 제한되고 불안정한 환경에서도 효율적으로 데이터를 전송하기 위해 설계되었습니다. 이러한 특성 덕분에 오늘날 IoT, 모바일 앱, 실시간 시스템에서 널리 사용되고 있습니다.

여러분이 만드는 시스템에서 이런 요구사항이 있다면 MQTT가 좋은 선택입니다:

- **수천~수만 개의 디바이스가 동시에 연결되어야 함**: MQTT Broker는 단일 인스턴스에서도 수십만 개의 동시 연결을 처리할 수 있습니다. HTTP 기반 폴링 방식으로는 이 정도 규모를 효율적으로 처리하기 어렵습니다.
- **실시간으로 데이터를 주고받아야 함**: MQTT는 연결을 유지하므로 메시지 전달 지연이 밀리초 단위입니다. 새로운 연결을 맺는 오버헤드가 없어 실시간성이 중요한 시스템에 적합합니다.
- **네트워크 환경이 불안정함 (모바일, IoT)**: 작은 패킷 크기와 내장된 재연결 메커니즘으로 불안정한 네트워크에서도 안정적으로 동작합니다.
- **서버가 클라이언트에게 먼저 메시지를 보내야 함 (Push)**: HTTP에서는 클라이언트가 먼저 요청해야 하지만, MQTT에서는 서버(Publisher)가 언제든 클라이언트(Subscriber)에게 메시지를 보낼 수 있습니다.

---

## 1장. MQTT v5 개요

이 장에서는 MQTT가 무엇인지, 왜 필요한지, 그리고 v5에서 어떤 점이 개선되었는지 알아봅니다. MQTT의 기본 개념을 이해하면 이후 장에서 다루는 Topic 설계, QoS 선택, 재연결 전략 등을 더 쉽게 이해할 수 있습니다.

### 1.1 MQTT란 무엇인가

#### 한 줄 정의

> MQTT는 **이벤트 기반 메시징 프로토콜**입니다.

MQTT는 "Message Queuing Telemetry Transport"의 약자입니다. 이름에 "Message Queuing"이 들어있지만, 실제로는 메시지 큐(RabbitMQ, Kafka 등)와는 다른 Pub/Sub 패턴을 사용합니다. 경량 프로토콜로 설계되어 헤더 크기가 최소 2바이트에 불과하며, 이는 HTTP 헤더(수백 바이트)와 비교하면 매우 작은 크기입니다.

조금 더 풀어서 설명하면:
- **이벤트 기반**: "무언가 일어났을 때" 메시지를 보냅니다. 주기적으로 상태를 확인하는 폴링 방식이 아니라, 변화가 발생했을 때만 데이터를 전송합니다.
- **메시징**: 데이터를 주고받는 방식입니다. 메시지는 Topic이라는 주소를 통해 분류되고 전달됩니다.
- **프로토콜**: 통신 규칙입니다. TCP/IP 위에서 동작하며, WebSocket을 통해 브라우저에서도 사용할 수 있습니다.

#### Broker 중심 구조

MQTT의 가장 큰 특징은 **Broker**가 중간에 있다는 것입니다. 이러한 구조를 "허브 앤 스포크(Hub and Spoke)" 패턴이라고도 합니다. 클라이언트들은 서로 직접 통신하지 않고, 모든 메시지가 중앙의 Broker를 통해 전달됩니다. 이 구조 덕분에 클라이언트는 다른 클라이언트의 존재나 위치를 알 필요가 없으며, Broker만 알면 됩니다.

```
[Publisher] ──publish──> [Broker] ──deliver──> [Subscriber]
                            │
                            ├──deliver──> [Subscriber]
                            │
                            └──deliver──> [Subscriber]
```

- **Publisher**: 메시지를 보내는 쪽. 데이터를 생성하고 특정 Topic으로 발행합니다.
- **Subscriber**: 메시지를 받는 쪽. 관심 있는 Topic을 구독하고 해당 메시지를 수신합니다.
- **Broker**: 메시지를 중계하는 서버. Publisher로부터 메시지를 받아 적절한 Subscriber에게 전달합니다.

Publisher와 Subscriber는 서로를 직접 알지 못합니다. Broker만 알면 됩니다. 이러한 느슨한 결합(Loose Coupling) 덕분에 시스템 확장이 용이합니다. 새로운 Subscriber를 추가해도 Publisher를 수정할 필요가 없고, 반대로 Publisher를 추가해도 기존 Subscriber에 영향을 주지 않습니다.

#### HTTP와의 근본적인 차이

MQTT를 이해하는 가장 좋은 방법은 익숙한 HTTP와 비교하는 것입니다. 두 프로토콜은 근본적으로 다른 통신 패턴을 사용하며, 각각 적합한 사용 사례가 있습니다.

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

MQTT v5는 2019년에 OASIS 표준으로 발표되었습니다. v3.1.1이 2014년에 발표된 이후 5년간의 실무 경험을 바탕으로 많은 개선이 이루어졌습니다. 특히 대규모 IoT 시스템을 운영하면서 발견된 문제점들을 해결하는 데 초점을 맞췄습니다.

#### v3의 한계

MQTT v3.1.1은 오랫동안 사용되었지만 실무에서 몇 가지 한계가 드러났습니다. 특히 장애 상황에서 원인을 파악하거나 시스템을 확장하는 데 어려움이 있었습니다.

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

이 장에서는 MQTT 시스템을 구성하는 핵심 요소들과 메시지가 어떻게 흘러가는지 살펴봅니다. Client, Broker, Topic의 관계를 이해하면 MQTT 기반 시스템을 설계할 때 각 구성 요소의 역할과 책임을 명확히 정의할 수 있습니다.

### 2.1 구성 요소

#### Client (클라이언트)

MQTT 시스템에서 메시지를 보내거나 받는 모든 것이 Client입니다. 여기서 중요한 점은 "클라이언트"라는 용어가 HTTP에서의 의미와 다르다는 것입니다. HTTP에서 클라이언트는 서버에 요청을 보내는 쪽이지만, MQTT에서 클라이언트는 Broker에 연결하는 모든 참여자를 의미합니다. 백엔드 서버도 Broker에 연결하면 MQTT 클라이언트가 됩니다.

- **센서, 모바일 앱, 서버 애플리케이션 모두 Client**: 온도 센서가 데이터를 보내든, 모바일 앱이 알림을 받든, 백엔드 서버가 명령을 내리든 모두 동등한 Client입니다.
- **Publisher이면서 동시에 Subscriber일 수 있음**: 하나의 Client가 특정 Topic에 메시지를 발행하면서 동시에 다른 Topic의 메시지를 구독할 수 있습니다. 예를 들어 스마트 조명은 명령을 구독하면서 자신의 상태를 발행합니다.
- **Broker에 연결해서 동작**: 모든 Client는 Broker와의 연결을 통해서만 메시지를 주고받습니다. Client들 사이에 직접 연결은 없습니다.

```
[모바일 앱] ──── Client
[온도 센서] ──── Client
[백엔드 서버] ── Client
```

#### Broker (브로커)

메시지를 중계하는 서버입니다. MQTT 시스템의 심장부로, 모든 메시지가 Broker를 통해 흐릅니다. Broker의 안정성과 성능이 전체 시스템의 품질을 결정하므로, 프로덕션 환경에서는 Broker 선택과 운영에 많은 주의를 기울여야 합니다.

**주요 역할:**
1. **Client 연결 관리**: 수천에서 수백만 개의 동시 연결을 관리합니다. 각 연결의 인증, 세션 상태, Keep Alive를 추적합니다.
2. **Topic별 메시지 라우팅**: Publisher가 보낸 메시지를 해당 Topic을 구독하는 모든 Subscriber에게 전달합니다. 이 과정이 효율적으로 이루어져야 지연 시간을 최소화할 수 있습니다.
3. **메시지 저장 (필요시)**: QoS 1, 2 메시지는 전달이 확인될 때까지 저장합니다. Retained Message는 새 구독자가 즉시 마지막 상태를 받을 수 있도록 저장합니다.
4. **인증/인가**: 클라이언트가 누구인지 확인하고(인증), 어떤 Topic에 접근할 수 있는지 결정합니다(인가). 보안이 중요한 시스템에서는 필수입니다.

**대표적인 Broker:**

| Broker | 언어 | 라이선스 | 클러스터링 | 적합한 환경 |
|--------|------|----------|------------|-------------|
| Mosquitto | C | EPL/EDL (오픈소스) | X (기본) | 학습, 소규모, 에지 |
| EMQX | Erlang/OTP | Apache 2.0 / 상용 | O | 대규모 IoT, 엔터프라이즈 |
| HiveMQ | Java | 상용 / CE 무료 | O | 엔터프라이즈, 미션크리티컬 |
| VerneMQ | Erlang | Apache 2.0 | O | 중대규모, 클라우드 |
| NanoMQ | C | MIT | X | 에지, 임베디드 |

**1. Mosquitto (Eclipse)**

- 특징: 가장 널리 사용되는 오픈소스 브로커, 매우 가벼움 (메모리 수 MB), MQTT v5 완벽 지원
- 장점: 빠른 시작 가능, 리소스 사용량 최소, 안정적이고 검증된 구현
- 단점: 단일 노드만 지원, 클러스터링 미지원 (외부 솔루션 필요)
- 추천: 개발/테스트 환경, 소규모 IoT (수백~수천 연결), 에지 게이트웨이

**2. EMQX**

- 특징: Erlang/OTP 기반 (높은 동시성), 단일 클러스터에서 수백만 연결 지원, 규칙 엔진/데이터 브리지 내장
- 장점: 수평 확장 용이, 고가용성 기본 지원, REST API/대시보드 제공, Kafka/MySQL 등 네이티브 연동
- 단점: Mosquitto 대비 리소스 사용량 높음, 고급 기능은 엔터프라이즈 버전 필요
- 추천: 대규모 IoT 플랫폼, 상용 서비스, 클라우드 네이티브 환경

**3. HiveMQ**

- 특징: Java 기반 엔터프라이즈 브로커, Extension 시스템으로 기능 확장, Community Edition 무료 제공
- 장점: 안정성과 성능, 상용 지원 및 SLA, AWS/Azure 마켓플레이스 제공
- 단점: 상용 라이선스 비용, Community Edition 기능 제한
- 추천: 미션 크리티컬 시스템, 기업 IoT 인프라, 금융/의료 등 규제 산업

**4. VerneMQ**

- 특징: Erlang 기반 오픈소스, 분산 아키텍처 설계, 플러그인 시스템
- 장점: 완전 오픈소스, 클러스터링 무료 지원, 유연한 인증 플러그인
- 단점: EMQX 대비 커뮤니티 작음, 문서화 상대적 부족
- 추천: 오픈소스 기반 중대규모 시스템, 비용 절감이 중요한 환경

**5. NanoMQ**

- 특징: EMQ에서 개발한 초경량 브로커, 에지 컴퓨팅 특화, 브리지 기능 내장
- 장점: 극히 낮은 리소스 사용, 빠른 시작 시간, 에지-클라우드 브리지 용이
- 단점: 기능 제한적, 대규모 연결 미지원
- 추천: 에지 디바이스, 라즈베리파이 등 SBC, 로컬 게이트웨이

**선택 가이드:**

| 상황 | 추천 브로커 |
|------|-------------|
| 학습/개발 | Mosquitto |
| 소규모 (< 1만 연결) | Mosquitto, NanoMQ |
| 중규모 (1만~10만 연결) | EMQX, VerneMQ |
| 대규모 (> 10만 연결) | EMQX Enterprise, HiveMQ |
| 에지/게이트웨이 | NanoMQ, Mosquitto |
| 엔터프라이즈/SLA 필요 | HiveMQ, EMQX Enterprise |

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

Topic 설계는 MQTT 시스템의 **가장 중요한 설계 결정**입니다. 한 번 정해진 Topic 구조는 나중에 변경하기 매우 어렵습니다. 이미 운영 중인 시스템에서 Topic을 변경하려면 모든 Publisher와 Subscriber를 동시에 수정해야 하기 때문입니다. 따라서 처음부터 확장성과 유지보수성을 고려한 설계가 필수입니다.

이 장에서는 Topic 네이밍 규칙, Wildcard 사용법, 그리고 실무에서 검증된 Best Practice를 다룹니다. 이 내용을 숙지하면 수천 개의 디바이스가 연결된 시스템에서도 효율적으로 메시지를 관리할 수 있습니다.

### 3.1 Topic 구조와 규칙

#### 계층적 네이밍

Topic은 슬래시(/)로 계층을 나눕니다. 이 구조는 파일 시스템의 디렉터리 구조와 유사합니다. 계층적 구조를 사용하면 Wildcard를 통해 특정 범위의 메시지만 구독할 수 있어 매우 유연한 메시지 필터링이 가능합니다. 예를 들어, 3층의 모든 센서 데이터만 구독하거나, 특정 건물의 모든 온도 데이터만 구독하는 것이 가능해집니다.

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

Wildcard는 **Subscribe할 때만** 사용할 수 있습니다. Publish할 때는 사용할 수 없습니다. Wildcard를 사용하면 여러 Topic을 한 번에 구독할 수 있어 코드가 간결해지고 관리가 용이해집니다. 하지만 과도한 Wildcard 사용은 불필요한 메시지를 수신하게 되어 성능 저하를 일으킬 수 있으므로 주의가 필요합니다.

#### + (Single-Level Wildcard)

한 단계만 대체합니다. 정확히 하나의 Topic 레벨을 대신하며, 빈 레벨은 매칭되지 않습니다. 특정 위치의 값만 다른 여러 Topic을 구독할 때 유용합니다.

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

MQTT 메시지는 단순히 데이터만 담는 것이 아닙니다. v5에서는 Payload 외에도 User Properties, Message Expiry Interval 등 다양한 메타데이터를 함께 전송할 수 있습니다. 이 장에서는 메시지를 구성하는 요소들과 각각의 활용 방법을 알아봅니다. 올바른 메시지 모델링은 시스템의 확장성과 유지보수성에 큰 영향을 미칩니다.

### 4.1 Payload

Payload는 메시지의 **본문**입니다. 실제 데이터가 들어가며, MQTT 프로토콜은 Payload의 형식을 강제하지 않습니다. JSON, XML, Binary, 심지어 단순 문자열도 가능합니다. 이러한 유연성은 장점이자 단점입니다. 형식의 자유도가 높은 만큼 Publisher와 Subscriber 간의 명확한 약속(계약)이 필요합니다.

#### JSON 형식

가장 많이 사용되는 형식입니다. 대부분의 프로그래밍 언어에서 JSON 파싱 라이브러리를 제공하므로 구현이 쉽고, 사람이 읽을 수 있어 디버깅에 유리합니다.

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

QoS(Quality of Service)는 메시지 **전달 보장 수준**입니다. MQTT에서 가장 중요한 개념 중 하나로, 네트워크 상황과 메시지의 중요도에 따라 적절한 QoS를 선택해야 합니다. QoS 선택은 시스템의 신뢰성과 성능 사이의 트레이드오프입니다. 높은 QoS는 더 많은 네트워크 오버헤드와 지연을 발생시키지만, 메시지 전달을 더 강력하게 보장합니다.

이 장에서는 각 QoS 레벨의 동작 원리를 상세히 설명하고, 실무에서 어떤 상황에 어떤 QoS를 선택해야 하는지 알아봅니다. 특히 QoS 1에서 발생할 수 있는 중복 메시지 처리 방법도 다룹니다.

### 5.1 QoS 0 / 1 / 2 동작 원리

#### QoS 0: At Most Once (최대 한 번)

"보내고 잊어버린다" 방식입니다. 메시지를 한 번 전송하고 응답을 기다리지 않습니다. 네트워크 문제로 메시지가 유실되어도 재전송하지 않습니다. 가장 빠르고 가벼운 방식이지만, 메시지 전달을 보장하지 않습니다.

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

MQTT에서 세션(Session)은 단순히 TCP 연결을 넘어서는 개념입니다. 세션에는 구독 정보, 전달되지 않은 메시지, QoS 흐름 상태 등이 포함됩니다. 올바른 세션 관리는 네트워크가 불안정한 환경에서 메시지 손실을 방지하는 핵심입니다. 이 장에서는 세션의 생명주기와 Keep Alive 메커니즘, 그리고 Retained Message 활용법을 다룹니다.

### 6.1 Session Expiry Interval

세션은 Client와 Broker 간의 **연결 상태 정보**입니다. v5에서는 Session Expiry Interval을 통해 연결이 끊어진 후에도 세션을 얼마나 유지할지 세밀하게 제어할 수 있습니다. 이 기능은 모바일 앱처럼 연결이 자주 끊어지는 환경에서 특히 유용합니다.

#### Clean Start vs Session 유지

Clean Start 플래그는 연결 시 이전 세션을 어떻게 처리할지 결정합니다. 이 설정은 시스템의 동작 방식에 큰 영향을 미치므로 신중하게 선택해야 합니다.

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

많은 MQTT 튜토리얼이 연결과 메시지 전송만 다루지만, 실제 프로덕션 코드에서는 재연결 로직이 전체 코드의 상당 부분을 차지합니다. 네트워크는 반드시 끊기며, 이에 대한 준비 없이는 안정적인 서비스를 운영할 수 없습니다. 이 장에서는 재연결이 필요한 이유, 재연결 시 발생하는 문제들, 그리고 검증된 재연결 전략을 상세히 다룹니다.

### 7.1 재연결이 반드시 필요한 이유

#### 현실 세계의 네트워크

이상적인 세계에서는 한 번 연결하면 영원히 유지됩니다. 하지만 현실은 다릅니다. 네트워크 연결은 다양한 이유로 끊어질 수 있으며, 이는 버그가 아닌 정상적인 운영 환경의 일부입니다. 따라서 재연결은 예외 처리가 아니라 핵심 기능으로 설계해야 합니다.

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

MQTT v5에서는 실무에서 자주 필요한 고급 기능들이 추가되었습니다. 이 장에서는 로드 밸런싱을 위한 Shared Subscription, HTTP 스타일의 Request/Response 패턴, 그리고 디버깅에 필수적인 Reason Code를 다룹니다. 이 기능들을 활용하면 더 확장성 있고 운영하기 쉬운 시스템을 구축할 수 있습니다.

### 8.1 Shared Subscription

여러 Subscriber가 **메시지를 나눠서** 처리하는 기능입니다. 일반적인 MQTT 구독에서는 같은 Topic을 구독하는 모든 Subscriber가 동일한 메시지를 받습니다. 하지만 Shared Subscription을 사용하면 메시지가 구독자들 사이에 분배되어 로드 밸런싱 효과를 얻을 수 있습니다. 이는 대량의 메시지를 처리해야 하는 시스템에서 수평 확장을 가능하게 합니다.

#### 개념

Shared Subscription의 핵심은 같은 그룹 내의 Subscriber들이 메시지를 분배받는다는 것입니다. 이를 통해 단일 Subscriber의 처리 한계를 넘어서는 대량의 메시지를 처리할 수 있습니다.

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

응답을 받을 Topic을 요청에 포함시킵니다. 전체 흐름을 이해하는 것이 중요합니다.

![image-20260118172008215](/Users/user/GolandProjects/tutorials-go/docs/start/image-20260118172008215.png)



```
[요청자 Client A]                              [응답자 Client B]
       │                                              │
       ├─ 1. SUBSCRIBE: reply/client-123/status       │
       │     (응답 받을 Topic 미리 구독)              │
       │                                              │
       ├─ 2. PUBLISH ─────────────────────────────────┼──► 요청 수신
       │      topic: device/cmd/get_status            │
       │      response_topic: reply/client-123/status │
       │      correlation_data: req-001               │
       │                                              │
       │                                         처리 후
       │                                              │
       ◄──────────────────────────────────────────────┼─ 3. PUBLISH (응답)
         응답 수신                                    │      topic: reply/client-123/status
                                                      │      correlation_data: req-001
                                                      │      payload: {"status": "ok"}
```

**핵심 포인트:**

- **요청자**는 요청 전에 `response_topic`을 미리 **SUBSCRIBE** 해야 함
- **응답자**는 요청의 `response_topic`으로 **PUBLISH**하여 응답
- 요청과 응답 모두 PUBLISH이며, 구독은 응답을 받기 위한 사전 준비

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

MQTT 시스템의 보안은 세 가지 축으로 구성됩니다: 인증(Authentication), 인가(Authorization), 그리고 암호화(Encryption). 인증은 "당신이 누구인가"를 확인하고, 인가는 "무엇을 할 수 있는가"를 결정하며, 암호화는 "통신 내용이 노출되지 않는가"를 보장합니다. 특히 IoT 환경에서는 수많은 디바이스가 연결되므로 보안 설계가 더욱 중요합니다.

### 9.1 인증

Client가 **누구인지** 확인합니다. MQTT에서는 연결 시점에 인증이 이루어지며, 한 번 인증된 연결은 세션이 유지되는 동안 유효합니다. 인증에 실패하면 Broker는 연결을 거부하고, v5에서는 Reason Code를 통해 실패 원인을 알려줍니다.

#### Username / Password

가장 기본적인 인증 방식입니다. 설정이 간단하여 개발 및 테스트 환경에서 많이 사용됩니다. 하지만 보안 수준이 높지 않으므로 프로덕션에서는 TLS와 함께 사용하거나 다른 인증 방식을 고려해야 합니다.

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

인증된 Client가 **무엇을 할 수 있는지** 결정합니다. ACL(Access Control List)은 **Broker 측에서 설정**하며, 클라이언트 코드가 아닌 Broker의 설정 파일이나 관리 시스템에서 구성합니다.

#### Broker별 ACL 설정 방식

| Broker | 설정 방식 |
|--------|----------|
| **Mosquitto** | `acl_file` 설정 파일 (텍스트) |
| **EMQX** | Dashboard UI, REST API, 또는 외부 DB 연동 |
| **HiveMQ** | XML 설정 또는 Extension |
| **VerneMQ** | `vmq.acl` 파일 또는 플러그인 |

#### Topic 기반 ACL

**Mosquitto 설정 예시:**

```bash
# mosquitto.conf에서 ACL 파일 지정
password_file /mosquitto/config/passwd
acl_file /mosquitto/config/acl
```

```
# /mosquitto/config/acl
user sensor-001
topic read sensor/+/state     # 읽기만 가능
topic write sensor/001/#      # 자기 토픽만 쓰기

user admin
topic readwrite #             # 모든 토픽 읽기/쓰기
```

#### ACL 변경 적용 방법

ACL 파일을 수정한 후에는 Broker에 변경 사항을 적용해야 합니다.

**Mosquitto 적용 방법:**

| 방법 | 명령어 | 설명 |
|------|--------|------|
| 재시작 | `docker restart mosquitto` | 모든 연결 끊김 |
| 설정 리로드 | `kill -SIGHUP $(pidof mosquitto)` | 연결 유지하며 리로드 (Linux) |
| Dynamic Security | REST API 호출 | 런타임 변경 가능 (v2.0+) |

**Broker별 동적 변경 지원:**

| Broker | 동적 변경 | 방법 |
|--------|----------|------|
| **Mosquitto** | △ (플러그인 필요) | Dynamic Security 플러그인 |
| **EMQX** | ✅ | Dashboard/REST API로 즉시 반영 |
| **HiveMQ** | ✅ | Control Center에서 즉시 반영 |
| **VerneMQ** | △ | `vmq-admin` CLI로 리로드 |

#### Mosquitto Dynamic Security 플러그인

Mosquitto 2.0부터 제공되는 **Dynamic Security 플러그인**을 사용하면 Broker 재시작 없이 런타임에 사용자, 그룹, ACL을 관리할 수 있습니다.

**활성화 방법:**

```bash
# mosquitto.conf
listener 1883
allow_anonymous false
plugin /usr/lib/mosquitto_dynamic_security.so
plugin_opt_config_file /mosquitto/config/dynamic-security.json
```

**초기 설정:**

```bash
# 관리자 계정 생성
mosquitto_ctrl dynsec init /mosquitto/config/dynamic-security.json admin-user

# 클라이언트 추가
mosquitto_ctrl dynsec createClient sensor-001 -p password123

# ACL 설정
mosquitto_ctrl dynsec addClientRole sensor-001 sensor-role
```

**장점:**
- Broker 재시작 없이 사용자/권한 관리
- JSON 기반 설정으로 백업/복원 용이
- `mosquitto_ctrl` CLI 또는 MQTT 메시지로 관리 가능

#### mosquitto-go-auth 플러그인

외부 시스템과 연동하여 동적으로 ACL을 체크하려면 **mosquitto-go-auth** 플러그인을 사용할 수 있습니다. 이 오픈소스 플러그인은 매 요청마다 외부 Backend에서 권한을 조회하므로, Broker 재시작 없이 실시간으로 권한 변경이 반영됩니다.

- GitHub: https://github.com/iegomez/mosquitto-go-auth

**지원하는 Backend:**

| Backend | 설명 |
|---------|------|
| **HTTP** | 외부 API 호출로 인증/인가 |
| **PostgreSQL** | DB에서 사용자/ACL 조회 |
| **MySQL** | DB에서 사용자/ACL 조회 |
| **Redis** | 캐시 기반 빠른 조회 |
| **MongoDB** | Document DB 연동 |
| **JWT** | 토큰 기반 인증 |
| **SQLite** | 경량 DB |

**HTTP Backend 설정 예시:**

```bash
# mosquitto.conf
auth_plugin /mosquitto/go-auth.so
auth_opt_backends http
auth_opt_http_host your-auth-server.com
auth_opt_http_port 8080
auth_opt_http_aclcheck_uri /mqtt/acl
auth_opt_http_usercheck_uri /mqtt/user
```

**인증 서버 구현 예시 (Go):**

```go
// POST /mqtt/acl
// Body: {"username": "sensor-001", "topic": "sensor/001/data", "acc": 2}
// acc: 1=subscribe, 2=publish

func checkACL(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Username string `json:"username"`
        Topic    string `json:"topic"`
        Acc      int    `json:"acc"`  // 1: subscribe, 2: publish
    }
    json.NewDecoder(r.Body).Decode(&req)

    // DB에서 권한 조회 후 동적으로 판단
    allowed := checkPermissionFromDB(req.Username, req.Topic, req.Acc)

    if allowed {
        w.WriteHeader(http.StatusOK)       // 200: 허용
    } else {
        w.WriteHeader(http.StatusForbidden) // 403: 거부
    }
}
```

**장점:**
- 매 요청마다 실시간 ACL 체크
- Broker 재시작 없이 권한 변경 즉시 반영
- 비즈니스 로직에 맞는 복잡한 권한 체크 가능
- 기존 인증 시스템(LDAP, OAuth 등)과 통합 용이

프로덕션 환경에서는 동적 변경이 가능한 방식을 선택하는 것이 운영에 유리합니다.

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

통신 내용을 **암호화**합니다. MQTT는 기본적으로 평문 통신을 하므로, 인터넷을 통한 통신이나 민감한 데이터 전송 시 TLS를 반드시 적용해야 합니다.

#### MQTT 포트 규약

| 포트 | 프로토콜 | 설명 |
|------|----------|------|
| **1883** | MQTT (평문) | 개발/테스트 또는 폐쇄망 |
| **8883** | MQTTS (TLS) | 프로덕션 표준 |
| **8084** | WSS (WebSocket + TLS) | 브라우저 연결 시 |

#### TLS 인증 방식

| 방식 | 설명 | 사용 사례 |
|------|------|----------|
| **Server Auth Only** | 클라이언트가 서버 인증서 검증 | 일반적인 웹 서비스와 동일 |
| **Mutual TLS (mTLS)** | 서버와 클라이언트 상호 인증 | 높은 보안이 필요한 IoT |

#### Mosquitto TLS 설정

**1. 인증서 생성 (테스트용 Self-signed)**

```bash
# CA 인증서 생성
openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -sha256 -days 365 \
    -out ca.crt -subj "/CN=MQTT CA"

# 서버 인증서 생성
openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr \
    -subj "/CN=mqtt.example.com"
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key \
    -CAcreateserial -out server.crt -days 365 -sha256
```

**2. Mosquitto 설정 (mosquitto.conf)**

```bash
# 평문 포트 (개발용, 프로덕션에서는 비활성화 권장)
listener 1883 localhost

# TLS 포트
listener 8883
cafile /mosquitto/certs/ca.crt
certfile /mosquitto/certs/server.crt
keyfile /mosquitto/certs/server.key

# TLS 버전 설정 (1.2 이상 권장)
tls_version tlsv1.2

# Mutual TLS (클라이언트 인증서 요구 시)
# require_certificate true
# use_identity_as_username true
```

**3. Docker Compose 예시**

```yaml
version: '3'
services:
  mosquitto:
    image: eclipse-mosquitto:2
    ports:
      - "1883:1883"
      - "8883:8883"
    volumes:
      - ./mosquitto.conf:/mosquitto/config/mosquitto.conf
      - ./certs:/mosquitto/certs
```

#### Go 클라이언트 TLS 설정

```go
import (
    "crypto/tls"
    "crypto/x509"
    "os"
)

func createTLSConfig(caFile string) (*tls.Config, error) {
    // CA 인증서 로드
    caCert, err := os.ReadFile(caFile)
    if err != nil {
        return nil, err
    }

    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    return &tls.Config{
        RootCAs:            caCertPool,
        MinVersion:         tls.VersionTLS12,
        InsecureSkipVerify: false,  // 프로덕션에서는 반드시 false
    }, nil
}

// autopaho에서 사용
func main() {
    tlsConfig, _ := createTLSConfig("/path/to/ca.crt")

    brokerURL, _ := url.Parse("tls://mqtt.example.com:8883")

    config := autopaho.ClientConfig{
        BrokerUrls: []*url.URL{brokerURL},
        TlsCfg:     tlsConfig,  // TLS 설정 적용
        // ...
    }
}
```

#### Mutual TLS (mTLS) 설정

클라이언트도 인증서를 제출해야 하는 양방향 인증입니다.

```go
func createMutualTLSConfig(caFile, certFile, keyFile string) (*tls.Config, error) {
    // CA 인증서 로드
    caCert, err := os.ReadFile(caFile)
    if err != nil {
        return nil, err
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    // 클라이언트 인증서 로드
    clientCert, err := tls.LoadX509KeyPair(certFile, keyFile)
    if err != nil {
        return nil, err
    }

    return &tls.Config{
        RootCAs:      caCertPool,
        Certificates: []tls.Certificate{clientCert},
        MinVersion:   tls.VersionTLS12,
    }, nil
}
```

#### 성능과 보안의 균형

| 항목 | TLS 1.2 | TLS 1.3 |
|------|---------|---------|
| 핸드셰이크 | 2-RTT | 1-RTT (더 빠름) |
| 암호화 스위트 | 다양함 | 간소화됨 |
| 0-RTT 재연결 | X | O |
| 권장 환경 | 레거시 호환 필요 시 | 신규 시스템 |

**경량 디바이스 고려사항:**
- TLS 1.3 사용 권장 (핸드셰이크 오버헤드 감소)
- 하드웨어 암호화 가속 지원 여부 확인
- 리소스 제약 시 VPN으로 네트워크 레벨 보안 대체 고려

### 9.4 MQTT over WebSocket

브라우저는 TCP 소켓을 직접 사용할 수 없기 때문에, 웹 애플리케이션에서 MQTT를 사용하려면 WebSocket으로 감싸야 합니다. MQTT over WebSocket을 사용하면 프론트엔드에서도 실시간으로 MQTT Topic을 구독하고 메시지를 발행할 수 있습니다.

#### 동작 원리

```
[IoT 디바이스]                   [MQTT Broker]                  [웹 브라우저]
      │                              │                              │
      ├── TCP:1883 ─────────────────►│◄───────────── WSS:8084 ──────┤
      │   (MQTT 원본)                │   (MQTT over WebSocket)      │
      │                              │                              │
      ├── PUBLISH sensor/temp ──────►│──────────────────────────────►│
      │                              │         (실시간 전달)         │
      │                              │                              │
      │◄─────────────────────────────┤◄── PUBLISH command/device ───┤
      │   (명령 수신)                │                              │
```

**핵심 포인트:**
- WebSocket 클라이언트도 일반 MQTT 클라이언트와 **동일하게 동작**
- 같은 Topic을 공유하며 서로 메시지를 주고받을 수 있음
- Broker 입장에서는 연결 방식만 다를 뿐 동일한 클라이언트

#### Mosquitto WebSocket 설정

```bash
# mosquitto.conf

# 일반 MQTT (TCP)
listener 1883

# WebSocket (평문) - 개발용
listener 8083
protocol websockets

# WebSocket (TLS) - 프로덕션
listener 8084
protocol websockets
cafile /mosquitto/certs/ca.crt
certfile /mosquitto/certs/server.crt
keyfile /mosquitto/certs/server.key
```

#### 프론트엔드 연동 (JavaScript)

**MQTT.js 설치:**

```bash
npm install mqtt
```

**React/Vue 등에서 사용:**

```javascript
import mqtt from 'mqtt';

// WebSocket으로 MQTT Broker 연결
const client = mqtt.connect('wss://mqtt.example.com:8084/mqtt', {
    clientId: 'web-dashboard-' + Math.random().toString(16).substr(2, 8),
    username: 'dashboard-user',
    password: 'secret',
    clean: true,
});

// 연결 성공
client.on('connect', () => {
    console.log('MQTT Connected!');

    // Topic 구독
    client.subscribe('sensor/+/temperature', { qos: 1 });
    client.subscribe('sensor/+/humidity', { qos: 1 });
});

// 실시간 메시지 수신
client.on('message', (topic, message) => {
    const data = JSON.parse(message.toString());
    console.log(`${topic}:`, data);

    // 예: sensor/livingroom/temperature: { value: 25.5, unit: "°C" }
    // UI 업데이트 로직
    updateDashboard(topic, data);
});

// 메시지 발행 (디바이스 제어)
function sendCommand(deviceId, command) {
    client.publish(
        `command/${deviceId}`,
        JSON.stringify(command),
        { qos: 1 }
    );
}

// 사용 예: sendCommand('light-001', { action: 'turn_on', brightness: 80 });
```

**연결 해제:**

```javascript
// 컴포넌트 언마운트 시
client.end();
```

#### React Hook 예시

```javascript
import { useEffect, useState } from 'react';
import mqtt from 'mqtt';

function useMQTT(brokerUrl, topics) {
    const [messages, setMessages] = useState({});
    const [client, setClient] = useState(null);
    const [isConnected, setIsConnected] = useState(false);

    useEffect(() => {
        const mqttClient = mqtt.connect(brokerUrl);

        mqttClient.on('connect', () => {
            setIsConnected(true);
            topics.forEach(topic => mqttClient.subscribe(topic));
        });

        mqttClient.on('message', (topic, message) => {
            setMessages(prev => ({
                ...prev,
                [topic]: JSON.parse(message.toString())
            }));
        });

        mqttClient.on('error', (err) => console.error('MQTT Error:', err));
        mqttClient.on('close', () => setIsConnected(false));

        setClient(mqttClient);

        return () => mqttClient.end();
    }, [brokerUrl]);

    const publish = (topic, message) => {
        if (client && isConnected) {
            client.publish(topic, JSON.stringify(message));
        }
    };

    return { messages, isConnected, publish };
}

// 사용 예시
function Dashboard() {
    const { messages, isConnected, publish } = useMQTT(
        'wss://mqtt.example.com:8084/mqtt',
        ['sensor/+/temperature', 'sensor/+/humidity']
    );

    return (
        <div>
            <p>연결 상태: {isConnected ? '✅ 연결됨' : '❌ 끊김'}</p>
            <p>거실 온도: {messages['sensor/livingroom/temperature']?.value}°C</p>
            <button onClick={() => publish('command/ac', { action: 'turn_on' })}>
                에어컨 켜기
            </button>
        </div>
    );
}
```

#### 실제 사용 사례

| 사용 사례 | 설명 |
|----------|------|
| **IoT 대시보드** | 센서 데이터 실시간 시각화 |
| **스마트홈 앱** | 조명, 에어컨 등 원격 제어 |
| **실시간 알림** | 이벤트 발생 시 즉시 알림 |
| **채팅 애플리케이션** | 메시지 실시간 전송/수신 |
| **라이브 모니터링** | 서버/장비 상태 실시간 확인 |
| **협업 도구** | 문서 동시 편집 상태 공유 |

#### WebSocket vs HTTP Polling

| 항목 | WebSocket (MQTT) | HTTP Polling |
|------|------------------|--------------|
| 지연 시간 | 수 ms (실시간) | 폴링 간격에 의존 |
| 서버 부하 | 낮음 (연결 유지) | 높음 (반복 요청) |
| 양방향 통신 | ✅ 기본 지원 | ❌ 별도 구현 필요 |
| 배터리 소모 | 낮음 | 높음 |

**결론:** 실시간성이 중요한 웹 애플리케이션에서는 MQTT over WebSocket이 HTTP Polling보다 효율적입니다.

---

## 10장. Go + Paho (v5) 사용법

이 장에서는 Go 언어로 MQTT v5 클라이언트를 구현하는 방법을 다룹니다. Eclipse Paho 프로젝트에서 제공하는 `paho.golang` 패키지를 사용하며, 특히 자동 재연결을 지원하는 `autopaho` 패키지의 사용법을 중심으로 설명합니다. 앞서 배운 개념들을 실제 코드로 구현하는 방법을 익히면 바로 프로덕션에 적용할 수 있습니다.

### 10.1 Paho v5 구조 이해

Go에서 MQTT v5를 사용하려면 `eclipse/paho.golang` 패키지를 사용합니다. 이 패키지는 두 가지 레벨의 API를 제공합니다. `paho` 패키지는 저수준 API로 세밀한 제어가 가능하고, `autopaho` 패키지는 자동 재연결 등 편의 기능이 포함된 고수준 API입니다. 실무에서는 대부분 `autopaho`를 사용하는 것이 좋습니다.

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

MQTT 시스템을 프로덕션에서 안정적으로 운영하려면 적절한 모니터링과 장애 대응 전략이 필요합니다. 이 장에서는 반드시 모니터링해야 할 핵심 지표와 흔히 발생하는 장애 시나리오별 대응 방법을 다룹니다. 사전에 이러한 상황들을 준비해두면 장애 발생 시 빠르게 대응할 수 있습니다.

### 11.1 모니터링 포인트

MQTT 시스템의 건강 상태를 파악하기 위해 다음 지표들을 모니터링해야 합니다. 대부분의 Broker가 이러한 메트릭을 제공하며, EMQX나 HiveMQ 같은 엔터프라이즈 Broker는 대시보드를 통해 시각화할 수 있습니다.

#### 연결 수

연결 수는 시스템 부하를 가장 직접적으로 나타내는 지표입니다. 갑작스러운 연결 수 변화는 네트워크 장애나 클라이언트 문제를 의미할 수 있습니다.

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

### 11.3 Mosquitto 모니터링 도구

Mosquitto를 사용하는 경우 다양한 방법으로 Broker 상태를 모니터링할 수 있습니다. 환경과 규모에 따라 적합한 도구를 선택하세요.

#### $SYS Topic (내장 기능)

Mosquitto는 자체 상태 정보를 `$SYS/#` Topic으로 발행합니다. 별도 설치 없이 바로 사용할 수 있어 빠른 상태 확인에 유용합니다.

```bash
# 모든 시스템 메트릭 구독
mosquitto_sub -h localhost -t '$SYS/#' -v
```

**주요 메트릭:**

| Topic | 설명 |
|-------|------|
| `$SYS/broker/clients/connected` | 현재 연결된 클라이언트 수 |
| `$SYS/broker/clients/total` | 총 등록된 클라이언트 수 |
| `$SYS/broker/messages/received` | 수신한 총 메시지 수 |
| `$SYS/broker/messages/sent` | 발송한 총 메시지 수 |
| `$SYS/broker/load/messages/received/1min` | 1분간 수신 메시지 비율 |
| `$SYS/broker/load/publish/sent/1min` | 1분간 발송 메시지 비율 |
| `$SYS/broker/uptime` | Broker 가동 시간 (초) |
| `$SYS/broker/bytes/received` | 수신한 총 바이트 |
| `$SYS/broker/bytes/sent` | 발송한 총 바이트 |

**활성화 설정 (mosquitto.conf):**

```bash
# $SYS 메트릭 발행 간격 (초, 기본값 10)
sys_interval 10
```

#### MQTT Explorer (GUI 도구)

개발 및 테스트 환경에서 가장 쉽게 사용할 수 있는 데스크톱 앱입니다.

- **다운로드**: https://mqtt-explorer.com
- **주요 기능**:
  - Topic 트리 시각화
  - 실시간 메시지 모니터링
  - 메시지 발행/구독 테스트
  - Payload 히스토리 및 차트
  - Retained Message 관리

```
# 연결 설정 예시
Host: localhost
Port: 1883
Username: (선택)
Password: (선택)
```

#### Prometheus + Grafana

프로덕션 환경에서 권장하는 방식입니다. 메트릭 수집, 저장, 시각화, 알림까지 통합 관리할 수 있습니다.

**mosquitto-exporter 사용:**

```yaml
# docker-compose.yml
version: '3'
services:
  mosquitto:
    image: eclipse-mosquitto:2
    ports:
      - "1883:1883"
    volumes:
      - ./mosquitto.conf:/mosquitto/config/mosquitto.conf

  mosquitto-exporter:
    image: sapcc/mosquitto-exporter
    ports:
      - "9234:9234"
    environment:
      - BROKER_ENDPOINT=tcp://mosquitto:1883
    depends_on:
      - mosquitto

  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml

  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
```

**prometheus.yml:**

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'mosquitto'
    static_configs:
      - targets: ['mosquitto-exporter:9234']
```

**Grafana 대시보드 설정:**
1. Grafana 접속 (http://localhost:3000)
2. Data Source에 Prometheus 추가
3. Dashboard Import에서 Mosquitto 템플릿 검색 또는 직접 생성

#### Cedalo Management Center

Mosquitto를 만든 Cedalo에서 제공하는 공식 상용 관리 도구입니다.

- **사이트**: https://cedalo.com/mqtt-management-center
- **주요 기능**:
  - 웹 기반 대시보드
  - 실시간 클라이언트 관리
  - ACL 동적 관리 (GUI)
  - 클러스터 모니터링
  - 감사 로그

#### 환경별 추천 도구

| 환경 | 추천 도구 | 이유 |
|------|----------|------|
| **개발/테스트** | MQTT Explorer | 설치 쉽고 직관적인 GUI |
| **소규모 프로덕션** | $SYS Topic + 스크립트 | 추가 인프라 불필요 |
| **중규모 프로덕션** | Prometheus + Grafana | 알림, 히스토리, 대시보드 |
| **대규모/엔터프라이즈** | Cedalo 또는 EMQX 전환 | 전문 지원, 클러스터링 |

**$SYS Topic 모니터링 스크립트 예시 (Go):**

```go
func monitorBroker(cm *autopaho.ConnectionManager) {
    topics := []string{
        "$SYS/broker/clients/connected",
        "$SYS/broker/messages/received",
        "$SYS/broker/load/messages/received/1min",
    }

    for _, topic := range topics {
        cm.Subscribe(context.Background(), &paho.Subscribe{
            Subscriptions: []paho.SubscribeOptions{
                {Topic: topic, QoS: 0},
            },
        })
    }
}

// 메시지 핸들러에서 메트릭 수집
func handleSysMessage(msg *paho.Publish) {
    switch msg.Topic {
    case "$SYS/broker/clients/connected":
        clientCount, _ := strconv.Atoi(string(msg.Payload))
        if clientCount > threshold {
            alertSlack("클라이언트 수 임계치 초과: " + string(msg.Payload))
        }
    }
}
```

---

## 12장. MQTT v5 사용 판단 기준

모든 기술에는 적합한 사용처가 있습니다. MQTT는 강력한 프로토콜이지만, 모든 상황에 적합한 것은 아닙니다. 이 장에서는 MQTT를 선택해야 하는 상황과 다른 기술을 선택해야 하는 상황을 명확히 구분합니다. 잘못된 기술 선택은 프로젝트 전체에 영향을 미치므로, 프로젝트 초기에 올바른 판단을 내리는 것이 중요합니다.

### MQTT를 써야 하는 경우

다음과 같은 요구사항이 있다면 MQTT가 좋은 선택입니다. 하나 이상 해당된다면 MQTT를 검토해볼 가치가 있습니다.

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

이 스터디를 통해 MQTT v5의 핵심 개념부터 실무 적용까지 전체적인 그림을 그릴 수 있게 되었기를 바랍니다. 마지막으로 배운 내용을 정리하고, 실무에 적용하기 전에 확인해야 할 체크리스트를 제공합니다.

### 13.1 핵심 요약

지금까지 배운 내용 중 가장 중요한 포인트들을 정리합니다. 이 내용들은 MQTT 기반 시스템을 설계하고 구현할 때 항상 염두에 두어야 합니다.

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

---

## FAQ

### Q: Topic에 Wildcard를 지원하나요?

**A: 네, 지원합니다. 단, Subscribe할 때만 사용 가능합니다.**

MQTT는 두 가지 Wildcard를 제공합니다:

| Wildcard | 이름 | 설명 | 예시 |
|----------|------|------|------|
| `+` | Single-Level | 한 단계만 대체 | `home/+/temperature` → `home/livingroom/temperature` |
| `#` | Multi-Level | 해당 위치부터 모든 하위 레벨 대체 | `home/#` → `home/livingroom/temperature` |

**주의사항:**
- Publish할 때는 Wildcard 사용 불가 (정확한 Topic 명시 필요)
- `#`는 반드시 Topic의 마지막에만 위치해야 함
- `home/#/temperature`와 같은 형태는 잘못된 사용

자세한 내용은 [3.2 Wildcard](#32-wildcard) 섹션을 참고하세요.

### Q: Wildcard로 구독했을 때 실제 매칭된 Topic을 알 수 있나요?

**A: 네, 알 수 있습니다. 메시지에 항상 실제 Topic이 포함되어 전달됩니다.**

Wildcard로 구독하더라도 메시지를 받을 때는 정확한 Topic 정보가 함께 옵니다.

```
# 구독
SUBSCRIBE topic: home/+/temperature

# 메시지 수신 시
Message 1:
  topic: home/livingroom/temperature  ← 실제 Topic
  payload: 25

Message 2:
  topic: home/bedroom/temperature     ← 실제 Topic
  payload: 22
```

**Go Paho 예시:**

```go
router.RegisterHandler("home/+/temperature", func(msg *paho.Publish) {
    // msg.Topic에 실제 매칭된 Topic이 들어있음
    fmt.Printf("Topic: %s, Payload: %s\n", msg.Topic, msg.Payload)
    // 출력: Topic: home/livingroom/temperature, Payload: 25
    // 출력: Topic: home/bedroom/temperature, Payload: 22
})
```

이를 활용하면 Topic에서 방 이름 등을 파싱하여 처리할 수 있습니다.

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
