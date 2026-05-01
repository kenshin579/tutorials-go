# 1-acl backend

Go + Echo + GORM + SQLite 기반 ACL 풀스택 샘플의 백엔드.

## 실행

```bash
go run main.go
# 기본: :8080, DB 파일 wiki-acl.db (자동 생성, 시드 적용)
```

환경변수: `DB_DSN`, `JWT_SECRET`, `ADDR` (모두 옵션).

## 시드 계정

모든 사용자의 비밀번호: `password`

| Email | 역할(시나리오) |
|---|---|
| alice@example.com | Engineering Roadmap, Public Onboarding Guide owner / Q4 Marketing Plan read |
| bob@example.com | Engineering Roadmap edit / Q4·Onboarding read |
| carol@example.com | Q4 Marketing Plan owner / Engineering Roadmap·Onboarding read |
| dave@example.com | Public Onboarding Guide read만 |

## 주요 엔드포인트

| Method | Path | 권한 | 설명 |
|---|---|---|---|
| POST | `/auth/login` | public | 로그인 → JWT 발급 |
| GET | `/api/pages` | 인증 | 본인이 access 가능한 페이지 목록 |
| GET | `/api/pages/:id` | ACL `read` | 페이지 상세 |
| PUT | `/api/pages/:id` | ACL `edit` | 페이지 수정 |
| GET | `/api/pages/:id/acl` | 페이지 owner | 공유 목록 |
| POST | `/api/pages/:id/acl` | 페이지 owner | 권한 부여 |
| DELETE | `/api/pages/:id/acl/:userId` | 페이지 owner | 권한 회수 |

ACL 검증 규칙: 페이지 owner는 모든 action 가능. 그 외 사용자는 ACLEntry에 명시된 action만 가능. `read` 권한은 `edit` 권한이 있으면 자동 충족.

## 빠른 시연 (cURL)

```bash
# 1. alice 로그인
TOKEN=$(curl -s -X POST localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"password"}' | jq -r .token)

# 2. 본인이 접근 가능한 페이지 목록 (alice는 owner+ACL 합쳐서 3개)
curl -s -H "Authorization: Bearer $TOKEN" localhost:8080/api/pages | jq '.[].title'

# 3. dave 로그인 후 Engineering Roadmap 직접 접근 시도 → 403
DAVE_TOKEN=$(curl -s -X POST localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"dave@example.com","password":"password"}' | jq -r .token)
curl -s -i -H "Authorization: Bearer $DAVE_TOKEN" localhost:8080/api/pages/1
```

## 테스트

```bash
go test ./...
```

## 디렉토리 구조

```
backend/
├── main.go              # 앱 진입점 + DI 와이어
├── config/              # SQLite 연결, 시드 데이터
├── domain/              # 엔티티 + repository 인터페이스 + EvaluateACL 순수 함수
├── repository/          # GORM 기반 User/Page/ACL repository
├── usecase/             # Auth/Page/ACL 비즈니스 로직 (ACL 평가 통합)
├── http/
│   ├── handler/         # Auth/Page/ACL HTTP 핸들러
│   ├── middleware/      # JWT 인증 미들웨어
│   └── router.go        # 라우트 + 글로벌 미들웨어
└── pkg/
    ├── jwt/             # JWT 발급/검증 헬퍼
    └── passwordhash/    # bcrypt 래퍼
```
