# wiki-permissions — 웹 권한 모델 비교 시리즈 (ACL/RBAC/ABAC)

블로그 시리즈 "웹 애플리케이션 권한 모델 비교 — ACL/RBAC/ABAC"의 풀스택 샘플 코드. **시리즈 완결.**

## 메타 컨텍스트

사내 위키 / 협업 문서 도구(Notion·Confluence 풍)를 가상 시나리오로 두고, 각 편이 같은 도메인 코어 위에서 다른 권한 모델을 적용한다. 모델 차이가 자연스럽게 코드 차원에서 드러나도록 의도적으로 메타 컨텍스트를 공유한다.

## 구성

| 편 | 디렉토리 | 모델 | 시나리오 |
|---|---|---|---|
| 1편 | [`1-acl/`](./1-acl/) | Access Control List | 페이지마다 사용자에게 read/edit 직접 부여 |
| 2편 | [`2-rbac/`](./2-rbac/) | Role-Based Access Control | admin/editor/viewer 역할 기반 |
| 3편 | [`3-abac/`](./3-abac/) | Attribute-Based Access Control | 분류 + 부서 + 고용형태 등 속성 기반 정책 |

## 시리즈 종합 비교

| 영역 | 1편 ACL | 2편 RBAC | 3편 ABAC |
|---|---|---|---|
| 권한 데이터 | `ACLEntry` 1테이블 | Role + Permission + UserRole + RolePermission | 사용자/페이지 속성 + Department |
| 평가 함수 | 30줄 (owner short-circuit + edit→read) | 1줄 lookup | 60줄 (정책 우선순위) |
| 평가 출력 | bool | bool | `Decision{Allowed, Reason, Policy}` |
| 핵심 SQL | LEFT JOIN acl_entries | JOIN role_permissions + user_roles | (메모리 정책 평가) |
| owner 처리 | 자동 short-circuit | 무시 (한계) | 정책으로 명시 (한계 회복) |
| Frontend 게이팅 | 없음 | PermissionGate 사전 노출 제어 | Decision reason을 UX에 표시 |
| 한계 | 사용자/페이지 cross product 폭발 | role explosion + 컨텍스트 부재 | 정책 설계 복잡도 |

## 어느 모델을 언제 쓸까

- **ACL** — 사용자/리소스 수가 작고 개별 grant/share 의도가 분명할 때 (예: 파일 공유 SaaS)
- **RBAC** — 조직 구조에 맞춰 role을 도출 가능할 때 (대부분의 어드민 시스템)
- **ABAC** — 시간/위치/분류/소유 등 컨텍스트가 정책에 들어가야 할 때 (의료, 금융, 컴플라이언스)

> 실전에서는 RBAC + ABAC 또는 RBAC + ACL을 섞어 쓰는 경우가 흔하다. 본 시리즈는 비교를 위해 한 번에 하나씩 깊이 다뤘다.

## 공통 기술 스택

- Backend: Go + Echo + GORM + SQLite + JWT
- Frontend: React 19 + TypeScript + Vite + Tailwind v4

각 편의 코드는 self-contained이며, 하나의 디렉토리만 클론해도 독립 실행 가능하다.

## 동시 실행

각 편이 다른 포트를 사용하도록 설정되어 있어 세 편을 동시에 띄워 비교 시연 가능하다.

| 편 | Backend | Frontend |
|---|---|---|
| 1편 | :8080 | :3000 |
| 2편 | :8081 | :3001 |
| 3편 | :8082 | :3002 |
