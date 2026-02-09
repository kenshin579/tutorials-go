**좋아, 이건 ****Claude Code에게 그대로 던져도 잘 나올 수 있게**

**👉 ****명확하고 과하지 않은 PRD**로 정리해줄게요.

(“왜 만드는지 / 무엇을 만드는지 / 어디까지 만드는지”가 분명한 버전)

---

# **PRD: Echo + WebRTC + React 기반 1:1 화상 통화 실습 프로젝트**

## **1. 개요 (Overview)**

**본 프로젝트는 ****WebRTC 스터디의 최종 실습 과제**로,

Go(Echo) 기반 Signaling 서버와 React 18 프론트엔드를 사용하여

**맥북 내장 카메라를 활용한 1:1 화상 통화 애플리케이션**을 구현한다.

본 프로젝트의 목적은:

* WebRTC 핵심 개념을 실제 코드로 연결하고
* 실무에서 사용 가능한 구조(Echo + React)를 경험하며
* Claude Code를 활용해 **설계 → 구현 흐름을 자동화**하는 것이다.

---

## **2. 목표 (Goals)**

### **기능적 목표**

* **브라우저 간 ****1:1 WebRTC 화상 통화** 구현
* **WebRTC ****DataChannel 기반 텍스트 채팅** 구현
* **Go Echo 서버를 통한 ****WebRTC Signaling 처리**

### **학습 목표**

* WebRTC Offer / Answer / ICE 흐름 이해
* Signaling Server의 역할과 책임 분리
* React 최신 Hook 패턴에서 WebRTC 객체 관리 방법 이해
* Claude Code를 활용한 코드 생성 및 구조 설계

---

## **3. 범위 (Scope)**

### **포함 (In Scope)**

* 1:1 통신만 지원
* STUN 서버만 사용
* 로컬 개발 환경(macOS) 기준
* 인증 없이 roomId 기반 연결
* 단순 UI (기능 위주)

### **제외 (Out of Scope)**

* TURN 서버 구성
* SFU / MCU 구조
* 다자간 통화 (1:N, N:N)
* 사용자 인증 / 권한 관리
* 모바일 대응
* 배포 환경 구성

---

## **4. 전체 아키텍처**

### **구성 요소**

* **Frontend**
  * React 18 + Vite + TypeScript
  * WebRTC Peer 역할
* **Backend**
  * Go + Echo
  * WebSocket 기반 Signaling Server
* **Media**
  * 맥북 내장 카메라 및 마이크
  * WebRTC P2P 전송

### **아키텍처 다이어그램 (논리적)**

```
React Client A ──┐
                  │ WebSocket (Signaling)
      Echo Server ─┤
                  │
React Client B ──┘

Media / Data → WebRTC P2P
```

---

## **5. 기능 요구사항 (Functional Requirements)**

### **5.1 Room 연결**

* 사용자는 **roomId**를 기준으로 방에 입장한다.
* 동일한 **roomId**에 최대 2명만 입장 가능하다.
* 2명이 입장하면 WebRTC 연결을 시작한다.

---

### **5.2 WebRTC Signaling**

* Signaling은 WebSocket을 통해 수행한다.
* 지원하는 메시지 타입:
  * offer
  * answer
  * ice

#### **메시지 형식**

```
{
  "type": "offer | answer | ice",
  "payload": {}
}
```

---

### **5.3 Media Stream**

* **브라우저에서 **getUserMedia()**를 사용해**
  * video
  * audio
    스트림을 획득한다.
* Local video와 Remote video를 각각 화면에 표시한다.

---

### **5.4 DataChannel 채팅**

* WebRTC DataChannel을 통해 텍스트 메시지를 송수신한다.
* 채팅 메시지는 실시간으로 UI에 반영된다.

---

## **6. Frontend 요구사항 (React)**

### **기술 스택**

* React 18
* TypeScript
* Vite
* Web API (WebRTC, MediaDevices)

### **구조 요구사항**

* WebRTC 관련 로직은 **Custom Hook**으로 분리한다.
  * useWebRTC
  * useSignaling
* RTCPeerConnection** 객체는 **useRef**로 관리한다.**
* React re-render 시 WebRTC 객체가 재생성되지 않아야 한다.

### **주요 UI**

* Local Video 영역
* Remote Video 영역
* 채팅 입력창 및 메시지 목록

---

## **7. Backend 요구사항 (Go + Echo)**

### **기술 스택**

* Go
* Echo framework
* Gorilla WebSocket
* pion/webrtc (서버에서는 사용하지 않음, signaling만 담당)

### **역할**

* WebSocket 연결 관리
* roomId 기준 클라이언트 그룹 관리
* Signaling 메시지 중계 (relay)
* Media 처리 ❌

### **API**

* GET /ws?roomId={roomId}
  * WebSocket 업그레이드
  * Signaling 통신 채널

---

## **8. 비기능 요구사항 (Non-Functional)**

* 코드 가독성과 학습 목적을 우선한다.
* 최소한의 외부 의존성만 사용한다.
* 로컬 환경에서 바로 실행 가능해야 한다.
* 로그를 통해 Signaling 및 ICE 흐름을 확인할 수 있어야 한다.

---

## **9. 개발 환경**

* OS: macOS
* Browser: Chrome (기준)
* Node.js: LTS
* Go: 1.21+

---

## **10. 성공 기준 (Success Metrics)**

* 두 브라우저 탭에서 서로 영상/음성이 정상 송출된다.
* ICE Candidate 교환이 로그로 확인된다.
* DataChannel 채팅이 지연 없이 동작한다.
* React Hook 구조가 명확하게 분리되어 있다.
* Claude Code로 생성된 코드가 큰 수정 없이 동작한다.

---

## **11. 확장 아이디어 (Post-Scope)**

* TURN 서버 추가
* 1:N 구조로 확장 (SFU 개념 도입)
* JWT 기반 room 접근 제어
* 로봇 카메라 스트리밍으로 확장

---

### **✅ Claude Code 사용 팁 (한 줄 가이드)**

> “이 PRD를 기반으로,

> \*\*Go(Echo) 백엔드와 React(TypeScript) 프론트엔드 코드를 각각 생성하라.

> 학습용 예제이므로 가독성과 구조를 우선하라.”

---

원하면 다음으로:

* **Claude Code에 넣을 프롬프트 버전**
* **PRD → Task breakdown**
* **Backend / Frontend PRD 분리본**
* **블로그용 PRD 요약본**

중에서 뭐가 필요해?
