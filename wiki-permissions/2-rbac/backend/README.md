# 2-rbac backend

Go + Echo + GORM + SQLite 기반 RBAC 풀스택 샘플의 백엔드.

## 실행

```bash
go run main.go
# 기본: :8081, DB 파일 wiki-rbac.db (자동 생성, 시드 적용)
# 1편(:8080)과 동시에 띄워 비교 시연 가능
```

환경변수: `DB_DSN`, `JWT_SECRET`, `ADDR` (모두 옵션).

## 시드 계정

모든 사용자의 비밀번호: `password`

| Email | Role | 가능한 액션 요약 |
|---|---|---|
| alice@example.com | admin | 모든 페이지 모든 액션 + 사용자 role 관리 |
| bob@example.com | editor | 페이지 read / create / edit (delete 불가) |
| carol@example.com | viewer | 페이지 read만 |
| dave@example.com | viewer | 페이지 read만 |

## 권한 매트릭스

| role | pages:read | pages:create | pages:edit | pages:delete | users:read | users:manage |
|---|---|---|---|---|---|---|
| admin | ✓ | ✓ | ✓ | ✓ | ✓ | ✓ |
| editor | ✓ | ✓ | ✓ | - | - | - |
| viewer | ✓ | - | - | - | - | - |

## 주요 엔드포인트

| Method | Path | 권한 | 설명 |
|---|---|---|---|
| POST | `/auth/login` | public | 로그인 → JWT + permissions/roles |
| GET | `/api/pages` | `pages:read` | 페이지 목록 (모든 페이지) |
| GET | `/api/pages/:id` | `pages:read` | 페이지 상세 |
| POST | `/api/pages` | `pages:create` | 페이지 생성 (owner = 요청자) |
| PUT | `/api/pages/:id` | `pages:edit` | 페이지 수정 |
| DELETE | `/api/pages/:id` | `pages:delete` | 페이지 삭제 (admin 전용) |
| GET | `/api/users` | `users:manage` | 모든 사용자 + 각자 roles |
| POST | `/api/users/:id/roles` | `users:manage` | 사용자에게 role 부여 |
| DELETE | `/api/users/:id/roles/:roleId` | `users:manage` | 사용자 role 회수 |
| GET | `/api/roles` | `users:manage` | 모든 role + 각 role의 permissions |

RBAC 평가 위치: HTTP 미들웨어가 아니라 **usecase 단**에서 `HasPermission(perms, "resource:action")` lookup. 1편(ACL)과 동일한 디자인 결정.

## 빠른 시연 (cURL)

```bash
# 1. alice (admin) 로그인 → 6개 권한
TOKEN=$(curl -s -X POST localhost:8081/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"alice@example.com","password":"password"}' | jq -r .token)
curl -s -H "Authorization: Bearer $TOKEN" localhost:8081/api/users | jq '. | length'  # 4

# 2. bob (editor) 로그인 → pages:delete 시도 시 403
BOB_TOKEN=$(curl -s -X POST localhost:8081/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"bob@example.com","password":"password"}' | jq -r .token)
curl -s -i -X DELETE -H "Authorization: Bearer $BOB_TOKEN" localhost:8081/api/pages/1 | head -1
# HTTP/1.1 403 Forbidden

# 3. carol (viewer) 로그인 → /api/users 접근 시 403 (users:manage 없음)
CAROL_TOKEN=$(curl -s -X POST localhost:8081/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"carol@example.com","password":"password"}' | jq -r .token)
curl -s -i -H "Authorization: Bearer $CAROL_TOKEN" localhost:8081/api/users | head -1
# HTTP/1.1 403 Forbidden
```

## 1편(ACL)과의 차이

| 영역 | 1편 ACL | 2편 RBAC |
|---|---|---|
| 데이터 모델 | ACLEntry 1테이블 | Role + Permission + 자동 join 2테이블 |
| 평가 함수 | `EvaluateACL` 30줄 (owner short-circuit + edit→read 함의) | `HasPermission` 한 줄 lookup |
| 핵심 SQL | LEFT JOIN acl_entries | JOIN role_permissions + JOIN user_roles |
| owner_id | 모든 액션 short-circuit | 무시 (메타데이터) |
| 한계 시나리오 | 사용자/페이지 cross product 폭발 | "내 페이지만 수정" 표현 불가 (3편 ABAC 동기) |

## 테스트

```bash
go test ./...
```

## 디렉토리 구조

```
backend/
├── main.go              # 앱 진입점 + DI 와이어
├── config/              # SQLite 연결, 시드 (사용자/페이지/role/permission/매트릭스)
├── domain/              # 엔티티 + repository 인터페이스
├── repository/          # User/Page/Role/Permission GORM repo (PermissionRepository.FindByUserID JOIN 핵심)
├── usecase/             # Auth/Page/Role 비즈니스 로직 + HasPermission
├── http/
│   ├── handler/         # Auth/Page/Role HTTP 핸들러
│   ├── middleware/      # JWT 인증 미들웨어
│   └── router.go        # 라우트 + 글로벌 미들웨어
└── pkg/
    ├── jwt/             # JWT 발급/검증 (1편과 동일)
    └── passwordhash/    # bcrypt 래퍼 (1편과 동일)
```
