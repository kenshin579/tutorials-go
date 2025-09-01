# TODO List - Keycloak Tutorial Project

## Phase 1: 환경 설정
- [x] **Keycloak 설정**
  - [x] Docker로 Keycloak 실행
  - [x] 관리자 콘솔 접속 확인 (admin/admin)
  - [x] Realm 생성 (`myrealm`)
  - [x] Client 생성 (`myclient`)
  - [x] Test User 생성 (`myuser`)

## Phase 2: Backend 개발 (Go + Echo)
- [x] **프로젝트 초기 설정**
  - [x] `backend/` 디렉토리 생성
  - [x] `go.mod` 파일 생성 및 의존성 추가
  - [x] 프로젝트 구조 생성 (handlers, middleware, models)

- [x] **핵심 기능 구현**
  - [x] JWT 토큰 검증 미들웨어 구현 (`middleware/auth.go`)
  - [x] 사용자 정보 모델 정의 (`models/user.go`)
  - [x] 사용자 정보 핸들러 구현 (`handlers/user.go`)
  - [x] 메인 서버 구현 (`main.go`)
  - [x] CORS 설정 추가

- [x] **테스트 및 검증**
  - [x] 서버 실행 테스트 (localhost:8081)
  - [x] JWT 토큰 검증 로직 테스트
  - [x] API 엔드포인트 테스트 (`GET /api/user`)

## Phase 3: Frontend 개발 (React)
- [x] **프로젝트 초기 설정**
  - [x] `frontend/` 디렉토리에 React 앱 생성
  - [x] 필요한 의존성 설치 (keycloak-js, axios, react-router-dom)
  - [x] 프로젝트 구조 생성 (components, services)

- [x] **Keycloak 통합**
  - [x] Keycloak 설정 파일 생성 (`services/keycloak.ts`)
  - [x] API 서비스 구현 (`services/api.ts`)
  - [x] 토큰 자동 추가 인터셉터 구현

- [x] **컴포넌트 구현**
  - [x] 로그인 컴포넌트 (`components/Login.tsx`)
  - [x] 사용자 프로필 컴포넌트 (`components/UserProfile.tsx`)
  - [x] 보호된 라우트 컴포넌트 (`components/ProtectedRoute.tsx`)
  - [x] 메인 App 컴포넌트 업데이트

- [x] **라우팅 설정**
  - [x] React Router 설정
  - [x] 보호된 라우트 적용
  - [x] 로그인/로그아웃 플로우 구현

## Phase 4: 통합 테스트
- [ ] **기본 기능 테스트**
  - [ ] 전체 시스템 실행 (Keycloak + Backend + Frontend)
  - [ ] 로그인 플로우 테스트
  - [ ] 사용자 정보 표시 테스트
  - [ ] 로그아웃 플로우 테스트

- [ ] **에러 처리 테스트**
  - [ ] 로그인 실패 시나리오 테스트
  - [ ] 토큰 만료 시나리오 테스트
  - [ ] 네트워크 오류 처리 테스트
  - [ ] Backend 서버 중단 시 처리 테스트

- [ ] **사용자 경험 개선**
  - [ ] 로딩 상태 표시
  - [ ] 에러 메시지 표시
  - [ ] 기본적인 스타일링 적용

## Phase 5: 문서화 및 정리
- [ ] **실행 가이드 작성**
  - [ ] 환경 설정 가이드
  - [ ] 실행 순서 문서화
  - [ ] 트러블슈팅 가이드

- [ ] **코드 정리**
  - [ ] 코드 리뷰 및 리팩토링
  - [ ] 주석 추가
  - [ ] 불필요한 파일 정리

- [ ] **최종 테스트**
  - [ ] 전체 시나리오 재테스트
  - [ ] 문서와 실제 구현 일치성 확인

## 우선순위
1. **High**: Phase 1, Phase 2 핵심 기능, Phase 3 핵심 기능
2. **Medium**: Phase 4 기본 기능 테스트, Phase 3 라우팅
3. **Low**: Phase 4 에러 처리, Phase 5 문서화

## 예상 소요 시간
- Phase 1: 30분
- Phase 2: 2-3시간
- Phase 3: 3-4시간
- Phase 4: 1-2시간
- Phase 5: 1시간

**총 예상 시간: 7-10시간**

## 참고사항
- 각 Phase 완료 후 동작 확인 필수
- 문제 발생 시 implementation.md 참조
- 최신 라이브러리 버전 사용 (mcp context7 활용)
- 로컬 환경에서만 동작하므로 HTTPS, 복잡한 보안 설정 불필요
