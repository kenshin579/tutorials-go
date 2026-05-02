# 2편 — RBAC (Role-Based Access Control)

사용자에게 역할(Role)을 주고, 역할에 권한(Permission)을 매핑하는 권한 모델.

## 구성

- [`backend/`](./backend/) — Go + Echo + GORM + SQLite (port :8081)
- [`frontend/`](./frontend/) — React 19 + TS + Vite + Tailwind v4 (port :3001)

## 빠른 시작

```bash
# Backend
cd backend && go run main.go &

# Frontend (다른 터미널)
cd frontend && npm install && npm run dev
```

브라우저: http://localhost:3001

> 1편(`:8080` / `:3000`)과 동시에 띄워 비교 시연 가능.

## 시드 시나리오

| 사용자 | 비밀번호 | Role | 가능한 액션 |
|---|---|---|---|
| alice@example.com | password | admin | 모든 페이지 모든 액션 + 사용자 role 관리 |
| bob@example.com | password | editor | 페이지 read / create / edit |
| carol@example.com | password | viewer | 페이지 read만 |
| dave@example.com | password | viewer | 페이지 read만 |

## 핵심 코드

- 권한 데이터 JOIN: [`backend/repository/permission_repository.go`](./backend/repository/permission_repository.go) — user → role → permission 3-hop JOIN
- RBAC 평가: [`backend/usecase/page_usecase.go`](./backend/usecase/page_usecase.go) — `HasPermission(perms, "resource:action")` 한 줄 lookup
- Frontend 사전 게이팅: [`frontend/src/components/PermissionGate.tsx`](./frontend/src/components/PermissionGate.tsx)

## 1편(ACL)과의 차이

| 영역 | 1편 ACL | 2편 RBAC |
|---|---|---|
| 권한 데이터 | `ACLEntry(page_id, user_id, action)` 1테이블 | `Role + Permission + UserRole + RolePermission` |
| 평가 함수 | 30줄 순수 함수 (owner short-circuit + edit→read 함의) | 한 줄 lookup |
| 핵심 SQL | LEFT JOIN acl_entries | JOIN role_permissions + JOIN user_roles |
| Frontend 게이팅 | 없음 (서버 403으로 처리) | `PermissionGate` 컴포넌트로 사전 노출 제어 |
| owner_id | 모든 액션 short-circuit | 무시 (메타데이터로만 유지) |
| 한계 시나리오 | 사용자/페이지 cross product 폭발 | "내가 만든 페이지만 수정" 표현 불가 (3편 ABAC 동기) |

## 관련 블로그 글

- 시리즈: 웹 권한 모델 비교
- 2편: RBAC — 워크스페이스 역할 (작성 예정)
