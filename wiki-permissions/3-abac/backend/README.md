# 3-abac backend

Go + Echo + GORM + SQLite 기반 ABAC 풀스택 샘플의 백엔드. 시리즈 마지막 편.

## 실행

```bash
go run main.go
# 기본: :8082, DB 파일 wiki-abac.db (자동 생성, 시드 적용)
# 1편(:8080), 2편(:8081)과 동시에 띄워 비교 시연 가능
```

환경변수: `DB_DSN`, `JWT_SECRET`, `ADDR` (모두 옵션).

## 시드 계정 (모두 비밀번호 password)

| Email | Department | Employment |
|---|---|---|
| alice@example.com | Engineering | fulltime |
| bob@example.com | Engineering | fulltime |
| carol@example.com | Marketing | fulltime |
| dave@example.com | Marketing | contract |

## 시드 페이지

| Title | Confidentiality | Department | Owner |
|---|---|---|---|
| Engineering Roadmap | internal | Engineering | alice |
| Q4 Marketing Plan | confidential | Marketing | carol |
| Public Onboarding Guide | public | (none) | alice |

## 정책 매트릭스

평가 우선순위 순:

1. **Owner** — 페이지 owner는 모든 액션 허용 (RBAC가 잃었던 owner 개념을 ABAC가 회복)
2. **Public** — public 페이지는 누구나 read (edit은 owner만)
3. **Department-match** — internal/confidential 페이지는 같은 부서가 아니면 거부
4. **Internal** — 같은 부서 사용자에게 read+edit
5. **Confidential** — 같은 부서 + **정규직** read+edit (contract 거부)
6. **Default Deny**

### 12 케이스 매트릭스 (4 사용자 × 3 페이지)

| | EngRoadmap (internal/Eng) | Q4MktPlan (confidential/Mkt) | Onboarding (public) |
|---|---|---|---|
| **alice** (Eng/fulltime) | owner → all | × 다른 부서 | read (public) |
| **bob** (Eng/fulltime) | read+edit (같은 부서) | × 다른 부서 | read (public) |
| **carol** (Mkt/fulltime) | × 다른 부서 | owner → all | read (public) |
| **dave** (Mkt/contract) | × 다른 부서 | × contract → confidential 거부 | read (public) |

## 주요 엔드포인트

| Method | Path | 설명 |
|---|---|---|
| POST | `/auth/login` | 로그인 → JWT + 사용자(department/employment_type) |
| GET | `/api/pages` | 정책 read 통과 페이지만 반환 |
| GET | `/api/pages/:id` | 페이지 상세 + Decision (`{can_read, can_edit, reason, policy}`) |
| PUT | `/api/pages/:id` | 정책 edit 통과 시 갱신 |
| GET | `/api/departments` | 부서 목록 |

## 빠른 시연 (cURL)

```bash
# 1. dave (Mkt/contract) 로그인 → Onboarding만 보임 (1개)
DAVE=$(curl -s -X POST localhost:8082/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"dave@example.com","password":"password"}' | jq -r .token)
curl -s -H "Authorization: Bearer $DAVE" localhost:8082/api/pages | jq '. | length'  # 1

# 2. dave가 Q4 (confidential/Mkt) 직접 접근 → 403 (contract → confidential 거부)
curl -s -i -H "Authorization: Bearer $DAVE" localhost:8082/api/pages/2 | head -1
# HTTP/1.1 403 Forbidden

# 3. bob (Eng/fulltime)이 Engineering Roadmap 상세 → can_edit 통과 + reason 표시
BOB=$(curl -s -X POST localhost:8082/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"bob@example.com","password":"password"}' | jq -r .token)
curl -s -H "Authorization: Bearer $BOB" localhost:8082/api/pages/1 | jq '.can_read, .can_edit'
# {"allowed":true,"reason":"same department, internal page","policy":"internal"}
```

## 1·2편과의 핵심 차이

| 영역 | 1편 ACL | 2편 RBAC | 3편 ABAC |
|---|---|---|---|
| 권한 데이터 | `ACLEntry` 1테이블 | `Role + Permission + UserRole + RolePermission` | 사용자/페이지 속성 컬럼 + Department |
| 평가 함수 | `EvaluateACL` 30줄 | `HasPermission` 1줄 | `EvaluateABAC` 60줄 (정책 우선순위) |
| 평가 출력 | bool | bool | `Decision{Allowed, Reason, Policy}` |
| 핵심 SQL | LEFT JOIN acl_entries | JOIN role_permissions + JOIN user_roles | 메모리 정책 평가 (DB JOIN 최소) |
| owner 처리 | 자동 short-circuit | 무시 (한계) | 정책으로 명시 (한계 회복) |
| 한계 시나리오 | cross product 폭발 | role explosion + 컨텍스트 부재 | 정책 설계의 복잡도 (정책 카탈로그 운영이 핵심) |

## 운영 환경 안내

본 샘플은 학습 목적상 **정책을 Go 함수로 hardcode** 했다. 운영에서는:

- **OPA (Open Policy Agent)** + Rego 정책 언어
- **Cedar** (Amazon)
- **Casbin**

같은 정책 엔진을 도입해 정책을 코드와 분리하고 런타임에 변경 가능하게 만드는 것이 일반적이다. ABAC의 데이터 모델(사용자/리소스 속성 + 환경 컨텍스트)은 동일하게 가져갈 수 있다.

## 테스트

```bash
go test ./...
```

## 디렉토리 구조

```
backend/
├── main.go              # 앱 진입점 + DI 와이어
├── config/              # SQLite 연결, 시드
├── domain/              # 엔티티 + repository 인터페이스 + EvaluateABAC (시리즈 핵심)
├── repository/          # User/Page/Department GORM repo
├── usecase/             # Auth/Page 비즈니스 로직 (정책 평가 통합)
├── http/
│   ├── handler/         # Auth/Page/Department HTTP 핸들러
│   ├── middleware/      # JWT 인증 미들웨어
│   └── router.go
└── pkg/
    ├── jwt/             # 1·2편과 동일
    └── passwordhash/    # 1·2편과 동일
```
