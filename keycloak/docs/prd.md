# 목적
- keycloak 를 인증 서버로 사용해서 FE, BE 개발을 하고 싶다

## 요구사항
- keycloak 기반으로 로그인, 로그아웃 기능을 개발한다
- 로그인 시 페이지에 사용자의 이름을 보여준다
- 로컬 개발에 맞게 keycloak 설정 필요

> 범위 안내: 본 저장소는 샘플 구현 목적입니다. 배포 및 운영(프로덕션 구성, 모니터링, 장애 대응, CI/CD, 비용/보안 최적화 등) 관련 작업은 현재 범위에 포함하지 않습니다. 문서 및 예시는 로컬 개발 환경 기준이며, Keycloak은 `infra/docker_run.sh`로 로컬에서 실행하는 것을 가정합니다.

### 현재 keycloak 설정
도커로 설치 이후 아래처럼 설정해둔 상태이다

- realm: myrealm 
- client: myclient
  - authentication flow: standard flow
  - valid redirect URIs: 미설정
  - web origins: 미설정
- 사용자: frank

## 구현 지침
- FE react를 사용해줘
- BE는 golang, echo 로 개발을 해줘
  - 구조는 clean architecture 사용해줘

## 요구사항 분석

### 1. Keycloak 설정 개선 필요사항
- `valid redirect URIs` 설정 필요
  - React 앱 URL: `http://localhost:3000/*`
  - 로그인 후 리다이렉트 URL 설정
- `web origins` 설정 필요
  - React 앱: `http://localhost:3000`
  - CORS 설정을 위한 origin 허용
- 로컬 개발 환경에 맞는 CORS 설정

### 2. 프론트엔드 (React) 요구사항
- Keycloak JavaScript Adapter 사용
- 로그인/로그아웃 기능 구현
- 사용자 정보 표시 (이름)
- 토큰 관리 및 갱신
- 보호된 라우트 구현

### 3. 백엔드 (Golang + Echo) 요구사항
- Clean Architecture 구조 적용
  - Domain Layer: 핵심 비즈니스 로직
  - UseCase Layer: 애플리케이션 서비스
  - Repository Layer: 데이터 접근
  - Handler Layer: HTTP 요청 처리
- Keycloak 토큰 검증 미들웨어
- 사용자 정보 조회 API
- CORS 설정
- 환경 설정 관리

### 4. 개발 환경 참고 (선택, 샘플 범위 밖)
- Docker Compose로 전체 환경 구성
- 로컬 개발용 설정
- 환경별 설정 분리 (dev, prod)

## 프로젝트 구조

```
keycloak/
├── docs/
│   └── prd.md
├── infra/
│   ├── docker_run.sh
│   ├── docker-compose.yml (선택, 샘플 범위 밖)
│   └── keycloak-config/ (선택, 샘플 범위 밖)
│       ├── realm-export.json
│       └── client-config.json
├── backend/ (새로 생성)
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── domain/
│   │   │   ├── user.go
│   │   │   └── auth.go
│   │   ├── usecase/
│   │   │   ├── user_usecase.go
│   │   │   └── auth_usecase.go
│   │   ├── repository/
│   │   │   └── keycloak_repository.go
│   │   └── handler/
│   │       ├── user_handler.go
│   │       └── auth_handler.go
│   ├── pkg/
│   │   ├── middleware/
│   │   │   └── auth.go
│   │   └── config/
│   │       └── config.go
│   ├── go.mod
│   └── go.sum
└── frontend/ (새로 생성)
    ├── src/
    │   ├── components/
    │   │   ├── Login.js
    │   │   ├── Logout.js
    │   │   └── UserInfo.js
    │   ├── services/
    │   │   └── keycloak.js
    │   ├── App.js
    │   └── index.js
    ├── public/
    │   └── index.html
    └── package.json
```

## 구현 우선순위

1. **Keycloak 설정 개선** (redirect URIs, web origins 설정)
2. **백엔드 Clean Architecture 구조** 생성
3. **프론트엔드 React 앱** 생성
4. (선택) **Docker Compose 환경** 구성
5. **통합 테스트** 및 검증

## 기술 스택

### 프론트엔드
- React 18+
- Keycloak JavaScript Adapter
- React Router (라우팅)
- Axios (HTTP 클라이언트)

### 백엔드
- Go 1.21+
- Echo Framework
- Keycloak Go Client
- JWT 토큰 검증
- Viper (설정 관리)

### 인프라
- Docker & Docker Compose
- Keycloak 26.3.2
- PostgreSQL (Keycloak DB)
