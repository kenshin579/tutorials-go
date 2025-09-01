# 목적

- keycloak을 이용해 로그인/로그아웃 샘플 페이지를 만들어 테스트해보고 싶다

## 요구사항

### 기능적 요구사항
- 인증 서버로 keycloak 를 사용한다
- 로그인 시 사용자의 정보(이름, 이메일)를 보여주는 페이지가 있어야 한다
- 로그아웃도 가능해야 한다
- 로그인 실패 시 에러 메시지 표시
- 세션 만료 시 자동 로그아웃

### 비기능적 요구사항
- 로컬 환경에서만 실행 (HTTP 통신)
- 최소한의 구현

## 기술 스택

### Backend
- Go 1.25
- Echo framework
- Keycloak Go 어댑터 또는 JWT 라이브러리

### Frontend
- React
- Keycloak JavaScript 어댑터
- React Router (페이지 보호용)

## 개발 환경 구성
- Keycloak: Docker로 localhost:8080에서 실행
- Backend (Echo): localhost:8081
- Frontend (React): localhost:3000
- 모든 통신은 HTTP로 진행 (로컬 환경)

## API 설계
- GET /api/user - 인증된 사용자 정보 반환 (이름, 이메일)

## 구현 플로우
1. Frontend에서 Keycloak JavaScript 어댑터로 로그인 처리
2. 로그인 성공 시 JWT 토큰을 받아서 Backend API 호출 시 헤더에 포함
3. Backend(Echo)에서 JWT 토큰을 검증하고 사용자 정보 반환
4. Frontend에서 사용자 정보 표시 및 로그아웃 기능 제공

## 테스트 시나리오
- Keycloak 관리자 콘솔에서 테스트 사용자 생성
- 로그인/로그아웃 플로우 테스트
- 토큰 만료 시나리오 테스트

## 구현 지침
- 구현은 최소한으로 한다
- FE 개발은 react 이용해서 개발하고 frontend 폴더에 작성한다
- BE 개발은 golang으로 개발하고 backend 폴더에 작성한다
- library 사용시 최신 코드는 mcp context7 사용해서 최신 코드로 작성한다