# wiki-permissions — 웹 권한 모델 비교 시리즈 (ACL/RBAC/ABAC)

블로그 시리즈 "웹 애플리케이션 권한 모델 비교 — ACL/RBAC/ABAC"의 풀스택 샘플 코드.

## 메타 컨텍스트

사내 위키 / 협업 문서 도구(Notion·Confluence 풍)를 가상 시나리오로 두고, 각 편이 그 안의 다른 측면을 다룬다. 같은 도메인 코어(User, Page) 위에 모델별 권한 데이터만 추가하는 방식이라, 모델 차이가 자연스럽게 코드 차원에서 드러난다.

## 구성

| 편 | 디렉토리 | 모델 | 시나리오 |
|---|---|---|---|
| 1편 | [`1-acl/`](./1-acl/) | Access Control List | 페이지마다 사용자에게 read/edit 직접 부여 |
| 2편 | [`2-rbac/`](./2-rbac/) | Role-Based Access Control | admin/editor/viewer 역할 기반 |
| 3편 | `3-abac/` (예정) | Attribute-Based Access Control | 분류 + 부서 + 고용형태 등 속성 기반 정책 |

## 공통 기술 스택

- Backend: Go + Echo + GORM + SQLite + JWT
- Frontend: React 19 + TypeScript + Vite + Tailwind v4

각 편의 코드는 self-contained이며, 하나의 디렉토리만 클론해도 독립 실행 가능하다.

## 동시 실행

각 편이 다른 포트를 사용하도록 설정되어 있어 1편과 2편을 동시에 띄워 비교 시연 가능하다.

| 편 | Backend | Frontend |
|---|---|---|
| 1편 | :8080 | :3000 |
| 2편 | :8081 | :3001 |
