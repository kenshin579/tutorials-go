# RBAC (Role-Based Access Control) + Owner-Based Access Control

Go + React로 구현한 역할 기반 권한 관리 시스템 예제 프로젝트.

## 기술 스택

| 구분 | 기술 |
|------|------|
| Backend | Go 1.25, Echo v4, GORM, golang-jwt/jwt v5 |
| Frontend | React 19, TypeScript, React Router v7, Axios, Tailwind CSS v4 |
| Database | MySQL 8.0 (Docker) |

## 프로젝트 구조

```
rbac/
├── backend/
│   ├── config/          # DB 설정, 시드 데이터
│   ├── domain/          # 엔티티, Repository 인터페이스, 상태 전이 규칙
│   ├── repository/      # GORM Repository 구현체
│   ├── usecase/         # 비즈니스 로직
│   ├── http/
│   │   ├── handler/     # HTTP 핸들러
│   │   ├── middleware/  # JWT 인증, RBAC, Owner 미들웨어
│   │   └── router.go   # 라우트 설정
│   ├── pkg/jwt/         # JWT 토큰 생성/검증
│   └── main.go
├── frontend/
│   └── src/
│       ├── api/         # Axios 클라이언트 (토큰 자동 갱신)
│       ├── auth/        # AuthContext, ProtectedRoute, usePermission
│       ├── components/  # Layout, Sidebar, PermissionGate
│       └── pages/       # Dashboard, Products, Orders, Users, Roles
└── docker-compose.yml   # MySQL 8.0
```

## 실행 방법

### 1. MySQL 실행

```bash
cd rbac
docker compose up -d
```

### 2. Backend 실행

```bash
cd rbac/backend
go mod tidy
go run main.go
# 서버: http://localhost:8081
```

서버 시작 시 자동으로:
- 테이블 마이그레이션
- 시드 데이터 삽입 (Permission 17개, Role 3개, 테스트 사용자 3명, 샘플 상품/주문)

### 3. Frontend 실행

```bash
cd rbac/frontend
npm install
npm run dev
# 서버: http://localhost:5173
```

## 테스트 계정

| 이메일 | 비밀번호 | Role | 설명 |
|--------|----------|------|------|
| admin@example.com | admin123 | admin | 모든 권한 |
| manager@example.com | manager123 | manager | 상품/주문 관리 (제한적) |
| user@example.com | user123 | user | 조회 + 본인 주문만 |

## 권한 모델

### RBAC (Role-Based Access Control)

Role → Permission 매핑으로 **"무엇을 할 수 있는가"** 를 제어합니다.

**3개 Role, 17개 Permission:**

| Resource | Permission | admin | manager | user |
|----------|-----------|:-----:|:-------:|:----:|
| users | read, create, update, delete | O | - | - |
| roles | read, create, update, delete | O | - | - |
| products | read | O | O | O |
| products | create, update, status:update | O | O | - |
| products | delete | O | - | - |
| orders | read, create | O | O | O |
| orders | status:update | O | O | - |
| orders | cancel | O | O | O |

### Owner-Based Access Control

리소스 소유자 확인으로 **"자기 것만 수정할 수 있는가"** 를 제어합니다.

**적용 대상:**
- 상품 수정 (`products.created_by`) — admin은 bypass
- 주문 상세/취소 (`orders.ordered_by`) — admin, manager는 bypass

### 미들웨어 체인

```
요청 → JWT 인증 → RBAC Permission 체크 → Owner 소유권 체크 → Handler
```

## 주문 상태 전이

```
pending → confirmed → shipped → completed
  ↓           ↓
cancelled  cancelled
```

**Role별 전이 권한:**

| 전이 | admin | manager | user |
|------|:-----:|:-------:|:----:|
| pending → confirmed | O | - | - |
| pending → cancelled | O | O | O |
| confirmed → shipped | O | O | - |
| confirmed → cancelled | O | O | - |
| shipped → completed | O | - | - |

## API 엔드포인트

### Auth
| Method | Path | 설명 |
|--------|------|------|
| POST | `/api/auth/register` | 회원가입 |
| POST | `/api/auth/login` | 로그인 (JWT 토큰 발급) |
| POST | `/api/auth/refresh` | 토큰 갱신 |
| POST | `/api/auth/logout` | 로그아웃 |

### Products
| Method | Path | Permission |
|--------|------|-----------|
| GET | `/api/products` | products:read |
| GET | `/api/products/:id` | products:read |
| POST | `/api/products` | products:create |
| PUT | `/api/products/:id` | products:update + Owner |
| DELETE | `/api/products/:id` | products:delete |
| PATCH | `/api/products/:id/status` | products:status:update |

### Orders
| Method | Path | Permission |
|--------|------|-----------|
| GET | `/api/orders` | orders:read |
| GET | `/api/orders/:id` | orders:read + Owner |
| POST | `/api/orders` | orders:create |
| PATCH | `/api/orders/:id/status` | orders:status:update |
| PATCH | `/api/orders/:id/cancel` | orders:cancel + Owner |

### Users / Roles / Permissions
| Method | Path | Permission |
|--------|------|-----------|
| GET/PUT/DELETE | `/api/users/:id` | users:read/update/delete |
| POST/DELETE | `/api/users/:id/roles` | users:update |
| GET/POST/PUT/DELETE | `/api/roles` | roles:* |
| GET | `/api/permissions` | roles:read |

## 테스트 실행

```bash
cd rbac/backend
go test ./... -v
```

## 관련 블로그 포스트

> 블로그 시리즈 작성 예정 — `3_admin_blog_prd.md` 참조
