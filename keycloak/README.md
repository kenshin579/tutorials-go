# tutorials-go/keycloak 실행 가이드

이 저장소는 Keycloak 기반 인증 데모입니다. 로컬에서 Keycloak(도커), Go 백엔드, React 프론트엔드를 실행하여 전체 인증 흐름을 테스트할 수 있습니다.

## 사전 요구사항
- Docker / Docker Compose (Keycloak 실행용)
- Go 1.20+ (백엔드)
- Node.js 18+ 및 npm (프론트엔드)

## 1) Keycloak 실행
1. Keycloak 컨테이너 실행
   ```bash
   cd infra
   ./docker_run.sh
   ```
2. 브라우저에서 http://localhost:8080 접속 → Admin Console 로그인 (기본 계정은 스크립트/환경에 따름)
3. Realm/클라이언트 설정 (docs/todo.md 기준)
   - Realm: `myrealm`
   - 프론트엔드 클라이언트: `myclient`
     - Client Protocol: `openid-connect`
     - Access Type: `public`
     - Standard Flow Enabled: ON
     - Valid Redirect URIs: `http://localhost:3000/*`, `http://localhost:3000`
     - Web Origins: `http://localhost:3000`
   - 백엔드 클라이언트: `mybackend`
     - Client Protocol: `openid-connect`
     - Access Type: `confidential`
     - Service Accounts Enabled: ON
     - Web Origins: `http://localhost:3000`
     - 발급된 Client Secret을 복사해 둡니다.
   - 테스트 사용자: 예) `frank` 사용자 생성 및 비밀번호 설정

## 2) 백엔드 실행 (Go)
1. 백엔드 설정 파일에서 Keycloak 정보를 확인/수정
   - 파일: `backend/pkg/config/config.go`
   - 기본값:
     - BaseURL: `http://localhost:8080`
     - Realm: `myrealm`
     - ClientID: `mybackend`
     - ClientSecret: `your-client-secret` → 실제 백엔드 클라이언트의 시크릿 값으로 교체하세요.
2. 서버 실행
   ```bash
   cd backend
   go run cmd/server/main.go
   ```
3. 서버가 8081 포트에서 시작됩니다. (http://localhost:8081)

## 3) 프론트엔드 실행 (React + Vite)
1. 의존성 설치
   ```bash
   cd frontend
   npm install
   ```
2. 개발 서버 시작 (다음 중 하나)
  ```bash
  # 저장소 루트에서 실행 (frontend로 프록시)
  npm start
  # 또는
  npm run dev
  # 또는 프론트엔드 디렉터리에서 실행
  cd frontend && npm start
  ```
3. 브라우저에서 http://localhost:3000 접속

### (선택) 환경 변수로 설정 변경
Vite는 `import.meta.env` 를 통해 환경 변수를 읽습니다. 필요 시 `frontend/.env` 파일을 만들어 아래 값을 조정합니다.
```
VITE_KEYCLOAK_URL=http://localhost:8080
VITE_KEYCLOAK_REALM=myrealm
VITE_KEYCLOAK_CLIENT_ID=myclient
VITE_API_BASE_URL=http://localhost:8081
```

## 4) 동작 확인 (인증 플로우)
1. http://localhost:3000 에 접속
2. 로그인 버튼 클릭 → Keycloak 로그인 페이지로 이동 → `frank` 계정 등으로 로그인
3. 로그인 후 "보호된 페이지" 이동 → 사용자 정보와 토큰 유효성 확인
4. 로그아웃 버튼 클릭 시 세션 종료

## 5) API 엔드포인트 수동 테스트
- 토큰이 있는 상태에서:
  - `GET http://localhost:8081/api/validate`
  - `GET http://localhost:8081/api/protected/user`
- 프론트엔드는 자동으로 토큰을 헤더에 추가하고, 401 시 토큰 갱신을 시도합니다.

## 문제 해결 가이드
- 401/403/CORS 오류
  - Keycloak의 `Valid Redirect URIs`, `Web Origins` 에 `http://localhost:3000` 이 들어있는지 확인
  - 백엔드 CORS 설정은 http://localhost:3000 을 허용하도록 되어 있습니다.
- 로그인 후 무한 리다이렉트/빈 화면
  - 프론트엔드 `.env` 의 realm/clientId 가 Keycloak 설정과 일치하는지 확인
- 백엔드 인증 실패
  - `backend/pkg/config/config.go`의 ClientSecret이 실제 백엔드 클라이언트 시크릿과 일치해야 합니다.

## 포트 요약
- Keycloak: 8080
- Backend (Echo): 8081
- Frontend (Vite): 3000

## 구조 요약
- `infra/docker_run.sh`: Keycloak 도커 실행 스크립트
- `backend/`: Go 백엔드 (Echo)
  - `/api/validate`, `/api/protected/user`
- `frontend/`: React (Vite), `keycloak-js` 연동

필요 시 docs/todo.md의 체크리스트를 참고하여 설정/검증을 진행하세요.
