# Keycloak 기반 인증 시스템 구현 TODO

> 참고: 본 저장소는 샘플 구현 목적입니다. 배포 및 운영(프로덕션 환경 구성, 모니터링, CI/CD 등) 관련 작업은 현재 범위에 포함되지 않습니다. 모든 체크리스트는 로컬 개발/학습 환경 기준이며, Keycloak은 `infra/docker_run.sh`로 로컬 실행하는 것을 가정합니다.

## ▶ Next Actions (다음 작업)
- Phase 1.1: Keycloak 서버 실행 및 상태 확인 (infra/docker_run.sh 실행, http://localhost:8080 접속)
- Phase 1.2: Realm/클라이언트 설정
  - Realm: myrealm 생성/확인
  - myclient (React) redirect URIs/Web Origins 설정
  - mybackend (백엔드) confidential + Service Accounts ON, client secret 확인
- backend/pkg/config/config.go의 ClientSecret 값을 Keycloak의 실제 시크릿으로 교체
- 백엔드 서버 실행: cd keycloak/backend && go run cmd/server/main.go
- 프론트엔드 Phase 3 시작: React 앱 생성 및 keycloak-js 연동

## 📋 Phase별 구현 계획

### 🚀 Phase 1: Keycloak 인프라 설정 및 기본 설정
**목표**: Keycloak 서버 실행 및 기본 클라이언트 설정 완료

#### 1.1 Keycloak 서버 실행
- [ ] Docker 환경에서 Keycloak 실행
  ```bash
  cd infra
  ./docker_run.sh
  ```
- [ ] Keycloak 서버 상태 확인 (http://localhost:8080)
- [ ] Admin Console 접속 테스트 (http://localhost:8080/admin)

**테스트 항목**:
- [ ] Keycloak 서버가 정상적으로 8080 포트에서 실행되는지 확인
- [ ] Admin Console에 admin/admin 계정으로 로그인 가능한지 확인
- [ ] 서버 로그에 에러가 없는지 확인

#### 1.2 Realm 및 클라이언트 설정
- [ ] Realm `myrealm` 생성/확인
- [ ] React 클라이언트 (`myclient`) 설정
  - [ ] Client Protocol: `openid-connect`
  - [ ] Access Type: `public`
  - [ ] Standard Flow Enabled: `ON`
  - [ ] Valid Redirect URIs: `http://localhost:3000/*`, `http://localhost:3000`
  - [ ] Web Origins: `http://localhost:3000`
- [ ] 백엔드 클라이언트 (`mybackend`) 설정
  - [ ] Client Protocol: `openid-connect`
  - [ ] Access Type: `confidential`
  - [ ] Service Accounts Enabled: `ON`
  - [ ] Valid Redirect URIs: `http://localhost:8081/api/*`, `http://localhost:8081/api`
  - [ ] Web Origins: `http://localhost:3000`
- [ ] 사용자 `frank` 계정 생성 및 비밀번호 설정

**테스트 항목**:
- [ ] `myrealm` Realm이 정상적으로 생성되었는지 확인
- [ ] `myclient` 클라이언트가 public 타입으로 설정되었는지 확인
- [ ] `mybackend` 클라이언트가 confidential 타입으로 설정되었는지 확인
- [ ] `frank` 사용자로 로그인 테스트
- [ ] 클라이언트 설정에서 CORS 에러가 발생하지 않는지 확인

---

### 🏗️ Phase 2: 백엔드 Clean Architecture 구조 구현
**목표**: Go 백엔드 서버의 기본 구조 및 인증 로직 구현

#### 2.1 프로젝트 구조 생성
- [x] 디렉토리 구조 생성
  ```bash
  mkdir -p keycloak/backend/{cmd/server,internal/{domain,usecase,repository,handler},pkg/{middleware,config}}
  ```
- [x] Go 모듈 초기화
  ```bash
  cd keycloak/backend
  go mod init github.com/kenshin579/tutorials-go/keycloak/backend
  ```
- [x] 의존성 설치
  ```bash
  go get github.com/labstack/echo/v4@v4.13.4
  go get github.com/golang-jwt/jwt/v5@v5.3.0
  ```

**테스트 항목**:
- [x] 모든 디렉토리가 정상적으로 생성되었는지 확인
- [x] `go mod init` 명령이 성공적으로 실행되는지 확인
- [x] 의존성 설치 시 에러가 없는지 확인
- [x] `go mod tidy` 실행 시 문제가 없는지 확인

#### 2.2 Domain Layer 구현
- [x] `internal/domain/user.go` 작성
  - [x] User 구조체 정의
  - [x] UserRepository 인터페이스 정의
  - [x] UserUseCase 인터페이스 정의
- [x] `internal/domain/auth.go` 작성
  - [x] AuthRequest 구조체 정의
  - [x] AuthResponse 구조체 정의
  - [x] AuthRepository 인터페이스 정의

**테스트 항목**:
- [x] Go 컴파일 에러가 없는지 확인
- [x] 구조체 필드 타입이 올바른지 확인
- [x] 인터페이스 메서드 시그니처가 일치하는지 확인

#### 2.3 Repository Layer 구현
- [x] `internal/repository/keycloak_repository.go` 작성
  - [x] KeycloakRepository 구조체 정의
  - [x] NewKeycloakRepository 함수 구현
  - [x] GetUserInfo 메서드 구현 (UserInfo 엔드포인트 호출)
  - [x] ValidateToken 메서드 구현 (Token Introspection 엔드포인트 호출)
  - [x] GetUserByID 메서드 구현 (Admin API 호출)
  - [x] getAdminToken 메서드 구현

**테스트 항목**:
- [ ] Keycloak API 연결 테스트
- [ ] GetUserInfo 메서드가 올바른 사용자 정보를 반환하는지 확인
- [ ] ValidateToken 메서드가 유효한 토큰을 올바르게 검증하는지 확인
- [ ] GetUserByID 메서드가 올바른 사용자 정보를 반환하는지 확인
- [ ] getAdminToken 메서드가 유효한 관리자 토큰을 반환하는지 확인

#### 2.4 UseCase Layer 구현
- [x] `internal/usecase/user_usecase.go` 작성
  - [x] UserUseCaseImpl 구조체 정의
  - [x] NewUserUseCase 함수 구현
  - [x] GetUserInfo 메서드 구현
  - [x] ValidateToken 메서드 구현

**테스트 항목**:
- [ ] UseCase 메서드들이 Repository를 올바르게 호출하는지 확인
- [ ] 에러 처리가 적절히 구현되었는지 확인
- [ ] 비즈니스 로직이 올바르게 동작하는지 확인

#### 2.5 Handler Layer 구현
- [x] `internal/handler/user_handler.go` 작성
  - [x] UserHandler 구조체 정의
  - [x] NewUserHandler 함수 구현
  - [x] GetUserInfo 핸들러 구현
  - [x] ValidateToken 핸들러 구현
  - [x] extractToken 헬퍼 함수 구현

**테스트 항목**:
- [ ] HTTP 요청이 올바른 핸들러로 라우팅되는지 확인
- [ ] 응답 상태 코드가 올바른지 확인
- [ ] JSON 응답 형식이 올바른지 확인
- [ ] 에러 응답이 적절한 형식으로 반환되는지 확인

#### 2.6 Middleware 구현
- [x] `pkg/middleware/auth.go` 작성
  - [x] AuthMiddleware 함수 구현
  - [x] extractToken 헬퍼 함수 구현

**테스트 항목**:
- [ ] 인증되지 않은 요청이 적절히 차단되는지 확인
- [ ] 유효한 토큰이 포함된 요청이 통과되는지 확인
- [ ] 잘못된 토큰이 포함된 요청이 적절히 처리되는지 확인

#### 2.7 Configuration 구현
- [x] `pkg/config/config.go` 작성
  - [x] Config 구조체 정의
  - [x] ServerConfig 구조체 정의
  - [x] KeycloakConfig 구조체 정의
  - [x] NewConfig 함수 구현

**테스트 항목**:
- [ ] 설정 파일이 올바르게 로드되는지 확인
- [ ] 환경 변수가 올바르게 파싱되는지 확인
- [ ] 기본값이 올바르게 설정되는지 확인

#### 2.8 Main Application 구현
- [x] `cmd/server/main.go` 작성
  - [x] 설정 로드
  - [x] Echo 인스턴스 생성
  - [x] 미들웨어 설정 (Logger, Recover, CORS)
  - [x] Repository, UseCase, Handler 생성
  - [x] 라우트 설정
  - [x] 서버 시작

**테스트 항목**:
- [ ] 서버가 8081 포트에서 정상적으로 시작되는지 확인
- [ ] 모든 미들웨어가 올바르게 적용되는지 확인
- [ ] 라우트가 올바르게 등록되는지 확인
- [ ] 서버 종료 시 정상적으로 종료되는지 확인

---

### 🎨 Phase 3: 프론트엔드 React 앱 구현
**목표**: React 기반 사용자 인터페이스 및 Keycloak 연동 구현

#### 3.1 프로젝트 생성 및 기본 설정
- [ ] React 앱 생성
  ```bash
  npx create-react-app keycloak/frontend
  cd keycloak/frontend
  ```
- [ ] 의존성 설치
  ```bash
  npm install keycloak-js axios react-router-dom
  ```

**테스트 항목**:
- [ ] React 앱이 정상적으로 생성되는지 확인
- [ ] 모든 의존성이 올바르게 설치되는지 확인
- [ ] 기본 React 앱이 정상적으로 실행되는지 확인

#### 3.2 Keycloak 설정 및 연동
- [ ] `src/services/keycloak.js` 작성
  - [ ] Keycloak 설정 객체 정의
  - [ ] initKeycloak 함수 구현
  - [ ] PKCE 방식 설정

**테스트 항목**:
- [ ] Keycloak 초기화가 성공적으로 완료되는지 확인
- [ ] 로그인 페이지로 올바르게 리다이렉트되는지 확인
- [ ] 인증 토큰이 올바르게 저장되는지 확인

#### 3.3 API 서비스 구현
- [ ] `src/services/api.js` 작성
  - [ ] axios 인스턴스 생성
  - [ ] 요청 인터셉터 (토큰 추가)
  - [ ] 응답 인터셉터 (토큰 갱신)
  - [ ] getUserInfo 함수 구현
  - [ ] validateToken 함수 구현

**테스트 항목**:
- [ ] API 요청에 인증 토큰이 자동으로 포함되는지 확인
- [ ] 토큰 만료 시 자동으로 갱신되는지 확인
- [ ] API 응답이 올바르게 처리되는지 확인

#### 3.4 컴포넌트 구현
- [ ] `src/components/Login.js` 작성
  - [ ] 로그인 버튼 컴포넌트
  - [ ] Keycloak 로그인 호출
- [ ] `src/components/Logout.js` 작성
  - [ ] 로그아웃 버튼 컴포넌트
  - [ ] Keycloak 로그아웃 호출
- [ ] `src/components/UserInfo.js` 작성
  - [ ] 사용자 정보 표시 컴포넌트
  - [ ] API 호출하여 사용자 정보 가져오기
  - [ ] 로딩 상태 처리
  - [ ] 에러 상태 처리

**테스트 항목**:
- [ ] 로그인 버튼 클릭 시 Keycloak 로그인 페이지로 이동하는지 확인
- [ ] 로그아웃 버튼 클릭 시 세션이 정상적으로 종료되는지 확인
- [ ] 사용자 정보가 올바르게 표시되는지 확인
- [ ] 로딩 상태와 에러 상태가 적절히 표시되는지 확인

#### 3.5 메인 앱 컴포넌트 및 라우팅
- [ ] `src/App.js` 수정
  - [ ] Keycloak 초기화
  - [ ] 인증 상태 관리
  - [ ] 라우팅 설정
  - [ ] 조건부 렌더링 (로그인/사용자 정보)

**테스트 항목**:
- [ ] 인증 상태에 따라 올바른 컴포넌트가 렌더링되는지 확인
- [ ] 라우팅이 올바르게 동작하는지 확인
- [ ] 인증 상태 변경 시 UI가 적절히 업데이트되는지 확인

#### 3.6 스타일링 및 UI/UX
- [ ] `src/App.css` 작성
  - [ ] 기본 레이아웃 스타일
  - [ ] 버튼 스타일
  - [ ] 사용자 정보 카드 스타일

**테스트 항목**:
- [ ] 모든 컴포넌트가 올바르게 스타일링되었는지 확인
- [ ] 반응형 디자인이 올바르게 동작하는지 확인
- [ ] 사용자 경험이 직관적인지 확인

---

### 🔗 Phase 4: 통합 테스트 및 검증
**목표**: 전체 시스템의 통합 테스트

#### 4.1 환경 시작 및 기본 연결 테스트
- [ ] Keycloak 실행 확인 (http://localhost:8080)
- [ ] 백엔드 서버 시작
  ```bash
  cd keycloak/backend
  go run cmd/server/main.go
  ```
- [ ] 프론트엔드 서버 시작
  ```bash
  cd keycloak/frontend
  npm start
  ```

**테스트 항목**:
- [ ] 모든 서비스가 정상적으로 시작되는지 확인
- [ ] 서비스 간 네트워크 연결이 정상적인지 확인
- [ ] 포트 충돌이 없는지 확인

#### 4.2 인증 플로우 통합 테스트
- [ ] React 앱 접속 (http://localhost:3000)
- [ ] 로그인 버튼 클릭
- [ ] Keycloak 로그인 페이지에서 사용자 정보 입력
- [ ] 로그인 후 사용자 정보 확인
- [ ] 로그아웃 테스트

**테스트 항목**:
- [ ] 전체 인증 플로우가 정상적으로 동작하는지 확인
- [ ] 로그인 후 세션이 올바르게 유지되는지 확인
- [ ] 로그아웃 후 세션이 정상적으로 종료되는지 확인
- [ ] 인증 상태가 프론트엔드와 백엔드 간에 일치하는지 확인

#### 4.3 API 엔드포인트 통합 테스트
- [ ] `GET /api/validate` (토큰 검증)
- [ ] `GET /api/protected/user` (사용자 정보)

**테스트 항목**:
- [ ] 인증된 요청이 정상적으로 처리되는지 확인
- [ ] 인증되지 않은 요청이 적절히 차단되는지 확인
- [ ] API 응답 형식이 올바른지 확인
- [ ] 에러 응답이 적절한 형식으로 반환되는지 확인

#### 4.4 에러 처리 및 예외 상황 테스트
- [ ] 잘못된 토큰으로 API 호출
- [ ] 만료된 토큰 처리
- [ ] 네트워크 에러 처리
- [ ] Keycloak 서버 다운 상황 처리

**테스트 항목**:
- [ ] 에러 상황에서 적절한 에러 메시지가 표시되는지 확인
- [ ] 사용자에게 명확한 피드백이 제공되는지 확인
- [ ] 시스템이 에러 상황에서 안정적으로 동작하는지 확인

#### 4.5 성능 및 부하 테스트 (샘플 범위 밖, 수행하지 않음)

---

### 🚀 Phase 5: 배포 및 운영 준비 (샘플 범위 밖, 수행하지 않음)
본 문서는 샘플 구현이므로 배포/운영 준비 항목은 범위 밖으로 제외합니다.

---

## 📊 Phase별 진행 상황

### Phase 1: Keycloak 인프라 설정 (0/8) - 0%
- [ ] Keycloak 서버 실행 (0/3)
- [ ] Realm 및 클라이언트 설정 (0/5)

### Phase 2: 백엔드 구현 (0/8) - 0%
- [ ] 프로젝트 구조 생성 (0/3)
- [ ] Domain Layer 구현 (0/2)
- [ ] Repository Layer 구현 (0/1)
- [ ] UseCase Layer 구현 (0/1)
- [ ] Handler Layer 구현 (0/1)

### Phase 3: 프론트엔드 구현 (0/6) - 0%
- [ ] 프로젝트 생성 및 기본 설정 (0/2)
- [ ] Keycloak 설정 및 연동 (0/1)
- [ ] API 서비스 구현 (0/1)
- [ ] 컴포넌트 구현 (0/1)
- [ ] 메인 앱 컴포넌트 및 라우팅 (0/1)

### Phase 4: 통합 테스트 및 검증 (0/4) - 0%
- [ ] 환경 시작 및 기본 연결 테스트 (0/3)
- [ ] 인증 플로우 통합 테스트 (0/4)
- [ ] API 엔드포인트 통합 테스트 (0/2)
- [ ] 에러 처리 및 예외 상황 테스트 (0/4)


**전체 진행률: 0%**

---

## 🧪 테스트 체크리스트 요약

### 단위 테스트
- [ ] Go 백엔드 코드 컴파일 테스트
- [ ] React 컴포넌트 렌더링 테스트
- [ ] API 엔드포인트 단위 테스트

### 통합 테스트
- [ ] Keycloak ↔ 백엔드 연동 테스트
- [ ] 백엔드 ↔ 프론트엔드 연동 테스트
- [ ] 전체 인증 플로우 테스트

### 성능 테스트 (샘플 범위 밖, 수행하지 않음)

### 보안 테스트
- [ ] 인증/인가 테스트
- [ ] 토큰 검증 테스트
- [ ] CORS 정책 테스트

### 사용성 테스트
- [ ] 사용자 인터페이스 테스트
- [ ] 에러 처리 테스트
- [ ] 반응형 디자인 테스트
