# 3편 — ABAC (Attribute-Based Access Control)

사용자/리소스/환경의 **속성(attribute)** 으로 정책을 평가하는 권한 모델. 시리즈 마지막 편.

## 구성

- [`backend/`](./backend/) — Go + Echo + GORM + SQLite (port :8082)
- [`frontend/`](./frontend/) — React 19 + TS + Vite + Tailwind v4 (port :3002)

## 빠른 시작

```bash
# Backend
cd backend && go run main.go &

# Frontend (다른 터미널)
cd frontend && npm install && npm run dev
```

브라우저: http://localhost:3002

> 1편(`:8080` / `:3000`), 2편(`:8081` / `:3001`)과 동시에 띄워 비교 시연 가능.

## 시드 시나리오

| 사용자 | Department | Employment |
|---|---|---|
| alice@example.com | Engineering | fulltime |
| bob@example.com | Engineering | fulltime |
| carol@example.com | Marketing | fulltime |
| dave@example.com | Marketing | contract |

| Page | Confidentiality | Department |
|---|---|---|
| Engineering Roadmap | internal | Engineering |
| Q4 Marketing Plan | confidential | Marketing |
| Public Onboarding Guide | public | (none) |

비밀번호 모두 `password`.

## 정책 매트릭스

평가 우선순위:

1. **Owner** — 페이지 owner는 모든 액션 허용
2. **Public** — public 페이지는 누구나 read
3. **Department-match** — internal/confidential은 같은 부서가 아니면 거부
4. **Internal** — 같은 부서 read+edit
5. **Confidential** — 같은 부서 + 정규직만 read+edit
6. **Default Deny**

### 12 케이스 매트릭스

|  | EngRoadmap (internal/Eng) | Q4 (confidential/Mkt) | Onboarding (public) |
|---|---|---|---|
| alice (Eng/fulltime) | owner → all | × 다른 부서 | read |
| bob (Eng/fulltime) | read+edit (같은 부서) | × 다른 부서 | read |
| carol (Mkt/fulltime) | × 다른 부서 | owner → all | read |
| dave (Mkt/contract) | × 다른 부서 | × contract | read |

## 핵심 코드

- 정책 평가기: [`backend/domain/policy.go`](./backend/domain/policy.go) — `EvaluateABAC` (4 정책 + Decision struct)
- 정책 통합: [`backend/usecase/page_usecase.go`](./backend/usecase/page_usecase.go)
- Decision 표시: [`frontend/src/pages/PageDetailPage.tsx`](./frontend/src/pages/PageDetailPage.tsx)

## 1·2·3편 비교

| 영역 | 1편 ACL | 2편 RBAC | 3편 ABAC |
|---|---|---|---|
| 권한 데이터 | `ACLEntry` 1테이블 | Role + Permission + UserRole + RolePermission | 사용자/페이지 속성 + Department |
| 평가 함수 | `EvaluateACL` 30줄 | `HasPermission` 1줄 | `EvaluateABAC` 60줄 (정책 우선순위) |
| 평가 출력 | bool | bool | Decision struct (allowed + reason + policy) |
| 핵심 SQL | LEFT JOIN acl_entries | JOIN role_permissions + user_roles | 메모리 정책 평가 |
| owner 처리 | 자동 short-circuit | 무시 (한계) | 정책으로 명시 (한계 회복) |
| Frontend 표현 | 서버 403만 | PermissionGate 사전 게이팅 | Decision reason을 사용자에게 표시 |
| 한계 | cross product 폭발 | role explosion + 컨텍스트 부재 | 정책 설계의 복잡도 |

## 운영 환경 안내

본 샘플은 학습 목적상 정책을 Go 함수로 hardcode 했다. 운영에서는 [OPA + Rego](https://www.openpolicyagent.org/), [Cedar](https://www.cedarpolicy.com/), [Casbin](https://casbin.org/) 같은 정책 엔진 도입을 권장한다.

## 관련 블로그 글

- 시리즈: 웹 권한 모델 비교
- 3편: ABAC — 분류와 속성 기반 정책 (작성 예정)
