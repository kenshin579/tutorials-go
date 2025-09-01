# Keycloak Tutorial Project

Keycloak을 이용한 로그인/로그아웃 샘플 페이지 프로젝트입니다.

## 프로젝트 구조

```
keycloak2/
├── backend/          # Go + Echo 백엔드 서버
├── frontend/         # React + TypeScript 프론트엔드
├── docs/            # 프로젝트 문서
│   ├── prd.md       # 프로젝트 요구사항 문서
│   ├── implementation.md  # 구현 가이드
│   └── todo.md      # 작업 목록
└── README.md        # 실행 가이드 (이 파일)
```

## 기술 스택

### Backend
- **언어**: Go 1.25
- **프레임워크**: Echo v4.13.4
- **인증**: JWT 토큰 검증
- **포트**: 8081

### Frontend
- **언어**: TypeScript
- **프레임워크**: React 18
- **라우팅**: React Router DOM
- **인증**: Keycloak JavaScript 어댑터
- **HTTP 클라이언트**: Axios
- **포트**: 3000

### 인증 서버
- **Keycloak**: Docker 컨테이너
- **포트**: 8080
- **Realm**: `myrealm`
- **Client**: `myclient`
- **Test User**: `myuser`

## 사전 요구사항

- Docker
- Go 1.25+
- Node.js 18+
- npm

## 실행 방법

### 1. Keycloak 실행

```bash
# Keycloak Docker 컨테이너 실행
docker run -d -p 8080:8080 \
  -e KEYCLOAK_ADMIN=admin \
  -e KEYCLOAK_ADMIN_PASSWORD=admin \
  --name keycloak-tutorial \
  quay.io/keycloak/keycloak:latest start-dev
```

### 2. Keycloak 설정

1. **관리자 콘솔 접속**
   - URL: http://localhost:8080
   - 계정: admin / admin

2. **Realm 생성**
   - Name: `myrealm`

3. **Client 생성**
   - Client ID: `myclient`
   - Client Type: `OpenID Connect`
   - Valid redirect URIs: `http://localhost:3000/*`
   - Web origins: `http://localhost:3000`

4. **Test User 생성**
   - Username: `myuser`
   - Email: `myuser@example.com`
   - First Name: `My`
   - Last Name: `User`
   - Password: `password`

### 3. Backend 실행

```bash
# backend 디렉토리로 이동
cd backend

# 의존성 설치
go mod tidy

# 서버 실행
go run main.go
```

서버가 정상적으로 실행되면 다음과 같은 메시지가 출력됩니다:
```
   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.13.4
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:8081
```

### 4. Frontend 실행

```bash
# frontend 디렉토리로 이동
cd frontend

# 의존성 설치
npm install

# 개발 서버 실행
npm start
```

브라우저가 자동으로 열리며 http://localhost:3000 에서 앱이 실행됩니다.

## 테스트 시나리오

### 기본 플로우
1. http://localhost:3000 접속
2. "Login with Keycloak" 버튼 클릭
3. Keycloak 로그인 페이지에서 `myuser` / `password`로 로그인
4. 사용자 정보 페이지에서 이름, 이메일 확인
5. "Logout" 버튼으로 로그아웃

### API 테스트
```bash
# 1. Keycloak에서 토큰 획득 (실제로는 프론트엔드에서 자동 처리)
# 2. 토큰을 사용하여 API 호출
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" \
     http://localhost:8081/api/user
```

## 주요 엔드포인트

### Backend API
- `GET /health` - 서버 상태 확인
- `GET /api/user` - 인증된 사용자 정보 조회 (JWT 토큰 필요)

### Frontend Routes
- `/` - 루트 (로그인 상태에 따라 리다이렉트)
- `/login` - 로그인 페이지
- `/profile` - 사용자 프로필 페이지 (보호된 라우트)

### Keycloak
- `http://localhost:8080` - Keycloak 관리자 콘솔
- `http://localhost:8080/realms/myrealm/protocol/openid-connect/certs` - JWKS 엔드포인트

## 트러블슈팅

### 포트 충돌
```bash
# 포트 사용 중인 프로세스 확인 및 종료
lsof -ti:8080 | xargs kill -9  # Keycloak
lsof -ti:8081 | xargs kill -9  # Backend
lsof -ti:3000 | xargs kill -9  # Frontend
```

### Docker 컨테이너 관리
```bash
# 컨테이너 상태 확인
docker ps -a | grep keycloak

# 컨테이너 중지
docker stop keycloak-tutorial

# 컨테이너 재시작
docker start keycloak-tutorial

# 컨테이너 제거
docker rm keycloak-tutorial
```

### 로그 확인
```bash
# Backend 로그 (실행 중인 터미널에서 확인)
# Frontend 로그 (브라우저 개발자 도구 Console 탭)
# Keycloak 로그
docker logs keycloak-tutorial
```

## 개발 환경 변수

### Backend (.env)
```bash
KEYCLOAK_URL=http://localhost:8080
KEYCLOAK_REALM=myrealm
KEYCLOAK_CLIENT_ID=myclient
SERVER_PORT=8081
```

### Frontend (.env)
```bash
REACT_APP_KEYCLOAK_URL=http://localhost:8080
REACT_APP_KEYCLOAK_REALM=myrealm
REACT_APP_KEYCLOAK_CLIENT_ID=myclient
REACT_APP_API_URL=http://localhost:8081/api
```

## 보안 고려사항

- 이 프로젝트는 **로컬 개발 환경**에서만 사용하도록 설계되었습니다
- HTTP 통신을 사용하므로 프로덕션 환경에서는 HTTPS 설정이 필요합니다
- JWT 토큰은 브라우저 메모리에 저장되며, 새로고침 시 재로그인이 필요할 수 있습니다

## 참고 문서

- [PRD 문서](docs/prd.md) - 프로젝트 요구사항
- [구현 가이드](docs/implementation.md) - 상세 구현 내용
- [작업 목록](docs/todo.md) - 개발 진행 상황

## 라이선스

이 프로젝트는 학습 목적으로 작성되었습니다.
